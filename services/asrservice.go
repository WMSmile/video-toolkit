package services

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"os/exec"
	"video-toolkit/asr"
	appconfig "video-toolkit/config"
)

type ASRService struct {
	recognizer asr.Recognizer
}

// Init 初始化 ASR 服务
func (s *ASRService) Init() error {
	// 创建 ASR 配置
	asrConfig := asr.AsrConfig{
		ModelPath:            appconfig.AppCfg.ModelPath,
		TokensPath:           appconfig.AppCfg.TokensPath,
		VadModelPath:         appconfig.AppCfg.VadModelPath,
		PunctuationModelPath: appconfig.AppCfg.PunctuationModelPath,
		HotwordsPath:         appconfig.AppCfg.HotwordsPath,
		NumThreads:           4,
	}

	// 创建识别器
	recognizer, err := asr.NewRecognizer(asrConfig)
	if err != nil {
		return fmt.Errorf("创建识别器失败: %v", err)
	}

	s.recognizer = recognizer
	return nil
}

// Recognize 识别单个音频文件
func (s *ASRService) Recognize(filePath string) (string, error) {
	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", fmt.Errorf("文件不存在: %v", filePath)
	}

	// 检查文件扩展名
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext != ".mp3" && ext != ".wav" && ext != ".flac" && ext != ".m4a" {
		return "", fmt.Errorf("不支持的文件格式: %v", ext)
	}

	// 初始化识别器（如果尚未初始化）
	if s.recognizer == nil {
		if err := s.Init(); err != nil {
			return "", err
		}
	}

	// 模拟音频样本（实际项目中需要从文件中读取）
	samples := make([]float32, 16000) // 1秒的音频样本

	// 处理音频
	result := s.recognizer.Process(samples)
	return result, nil
}

// BatchRecognize 批量识别音频文件
func (s *ASRService) BatchRecognize(filePaths []string) (map[string]string, error) {
	// 初始化识别器（如果尚未初始化）
	if s.recognizer == nil {
		if err := s.Init(); err != nil {
			return nil, err
		}
	}

	results := make(map[string]string)
	for _, filePath := range filePaths {
		// 检查文件是否存在
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			results[filePath] = "错误: 文件不存在"
			fmt.Printf("文件不存在: %s\n", filePath)
			continue
		}

		// 检查文件扩展名
		ext := strings.ToLower(filepath.Ext(filePath))
		// if ext != ".mp3" && ext != ".wav" && ext != ".flac" && ext != ".m4a" {
		// 	results[filePath] = "错误: 不支持的文件格式"
		// 	fmt.Printf("不支持的文件格式: %s\n", filePath)
		// 	continue
		// }

		// 使用 ffprobe 获取采样率 (假设 ffprobe 在环境路径中)
		// 实际项目中可以把 ffprobe 路径也放在 config 里
		checkCmd := exec.Command(appconfig.AppCfg.FFprobePath, "-v", "error", "-select_streams", "a:0",
			"-show_entries", "stream=sample_rate", "-of", "default=noprint_wrappers=1:nokey=1", filePath)

		out, err := checkCmd.Output()
		currentSR := strings.TrimSpace(string(out))
		fmt.Printf("当前文件 %s 的采样率: %s\n", filePath, currentSR)

		targetFilePath := filePath
		isTempFile := false

		if err == nil && currentSR != "16000" && ext != ".wav" {
			fmt.Printf("当前采样率 %s 不符合要求，正在转换...\n", currentSR)
			// 调用 ConvertService 转换
			convertService := &ConvertService{}
			newPath, err := convertService.ConvertTo16kWave(filePath)
			if err != nil {
				results[filePath] = fmt.Sprintf("采样率转换失败: %v\n", err)
				continue
			}
			targetFilePath = newPath
			isTempFile = true // 标记是临时文件，后续需删除
		}

		fmt.Printf("转换后的文件路径: %s\n", targetFilePath)

		var totalResult strings.Builder

		// 采用 2 秒一个分块进行流式读取
		err = s.processLargeWav(targetFilePath, func(chunk []float32) error {
			// 实时处理分块
			part := s.recognizer.Process(chunk)
			if part != "" {
				totalResult.WriteString(part)
			}
			return nil
		})

		// 关键：文件读完后，强制刷出 VAD 缓冲区里最后一段话
		lastPart := s.recognizer.Flush()
		totalResult.WriteString(lastPart)
		result := totalResult.String()

		// 清理临时文件
		if isTempFile {
			os.Remove(targetFilePath)
		}
		results[filePath] = result
	}

	return results, nil
}

// 流式读取并处理
func (s *ASRService) processLargeWav(path string, handler func([]float32) error) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Seek(44, 0); err != nil {
		return fmt.Errorf("seek error: %v", err)
	}

	// 建议：每次读取 0.5 秒 (8000 个采样点)，这样 VAD 可以更灵敏地切割
	const samplesPerChunk = 8000
	byteBuffer := make([]byte, samplesPerChunk*2)

	for {
		n, err := io.ReadAtLeast(file, byteBuffer, 2) // 只要有数据就读，不强制读满 58 秒
		if n <= 0 {
			break
		}

		numSamples := n / 2
		samples := make([]float32, numSamples)
		for i := 0; i < numSamples; i++ {
			low := uint16(byteBuffer[i*2])
			high := uint16(byteBuffer[i*2+1])
			val := int16(high<<8 | low)
			samples[i] = float32(val) / 32768.0
		}

		if err := handler(samples); err != nil {
			return err
		}

		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		}
	}
	return nil
}
