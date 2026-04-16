package services

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	appconfig "video-toolkit/config"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type ConvertService struct {
	App *application.App
}

// ConvertFolder converts all video files in a folder to MP3
func (b *ConvertService) ConvertFolder(folder string) ([]string, error) {
	var list []string
	files, err := os.ReadDir(folder)
	if err != nil {
		return nil, err
	}

	// 获取 FFmpeg 路径
	ffmpegPath := appconfig.AppCfg.FFmpegPath
	if ffmpegPath == "" {
		ffmpegPath = "ffmpeg" // 回退到默认值
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		ext := filepath.Ext(file.Name())
		if ext == ".mp4" || ext == ".mov" || ext == ".mkv" || ext == ".avi" {
			full := filepath.Join(folder, file.Name())
			out := strings.TrimSuffix(full, ext) + ".mp3"
			cmd := exec.Command(ffmpegPath, "-i", full, "-q:a", "0", "-y", out)
			cmd.CombinedOutput()
			list = append(list, "✅ "+file.Name())
		}
	}
	return list, nil
}

// ToMP3 converts a video file to MP3
func (f *ConvertService) ToMP3(input string) (string, error) {
	out := strings.TrimSuffix(input, filepath.Ext(input)) + ".mp3"

	// 获取 FFmpeg 路径
	ffmpegPath := appconfig.AppCfg.FFmpegPath
	if ffmpegPath == "" {
		ffmpegPath = "ffmpeg" // 回退到默认值
	}

	cmd := exec.Command(ffmpegPath, "-i", input, "-q:a", "0", "-y", out)
	_, err := cmd.CombinedOutput()
	println(out)
	println(err)
	return out, err
}

// 正则匹配 ffmpeg 进度输出
// 匹配格式：time=00:00:10.12  duration=00:01:00.45
var progressRegex = regexp.MustCompile(`time=(\d{2}):(\d{2}):(\d{2}\.\d+).*duration=(\d{2}):(\d{2}):(\d{2}\.\d+)`)

// ToMP3 转换视频为 MP3，支持实时进度回调
func (f *ConvertService) ToMP3Progress(input string, convertId string) (string, error) {
	// 生成输出文件名
	out := strings.TrimSuffix(input, filepath.Ext(input)) + ".mp3"

	// 获取 FFmpeg 路径
	ffmpegPath := appconfig.AppCfg.FFmpegPath
	if ffmpegPath == "" {
		ffmpegPath = "ffmpeg" // 回退到默认值
	}

	// 构建 ffmpeg 命令：禁止输出冗余日志，只输出进度
	cmd := exec.Command(ffmpegPath,
		"-i", input,
		"-q:a", "0", // 最高音质
		"-y",                  // 覆盖输出文件
		"-vn",                 // 禁用视频流（纯音频更快）
		"-progress", "pipe:2", // 把进度输出到 stderr
		"-nostats",
		out,
	)

	// 获取 stderr 管道（实时读取进度）
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", fmt.Errorf("创建管道失败: %w", err)
	}

	// 启动命令
	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("启动命令失败: %w", err)
	}

	// 协程：实时读取并解析进度
	go func() {
		scanner := bufio.NewScanner(stderr)
		// 加大缓冲区，避免大文件扫描失败
		buf := make([]byte, 1024*1024)
		scanner.Buffer(buf, 1024*1024)

		var totalSeconds float64

		for scanner.Scan() {
			line := scanner.Text()
			// 跳过空行
			if line == "" {
				continue
			}

			// 解析总时长和当前时长
			match := progressRegex.FindStringSubmatch(line)
			if len(match) == 7 {
				// 解析当前时间
				currentH, _ := strconv.ParseFloat(match[1], 64)
				currentM, _ := strconv.ParseFloat(match[2], 64)
				currentS, _ := strconv.ParseFloat(match[3], 64)
				currentSeconds := currentH*3600 + currentM*60 + currentS

				// 第一次解析总时长
				if totalSeconds == 0 {
					totalH, _ := strconv.ParseFloat(match[4], 64)
					totalM, _ := strconv.ParseFloat(match[5], 64)
					totalS, _ := strconv.ParseFloat(match[6], 64)
					totalSeconds = totalH*3600 + totalM*60 + totalS
				}

				// 计算进度（0~100）
				if totalSeconds > 0 {
					progress := (currentSeconds / totalSeconds) * 100
					// 防止进度超过 100
					if progress > 100 {
						progress = 100
					}

					// 发送进度事件
					if f.App != nil {
						data := map[string]interface{}{
							"convertId": convertId,
							"progress":  progress,
						}
						jsonData, err := json.Marshal(data)
						if err == nil {
							f.App.Event.Emit("convert_progress", string(jsonData))
						}
					}
				}
			}
		}
	}()

	// 等待命令执行完成
	err = cmd.Wait()
	if err != nil {
		return "", fmt.Errorf("转换失败: %w", err)
	}

	// 发送完成事件
	if f.App != nil {
		data := map[string]interface{}{
			"convertId": convertId,
			"progress":  100.0,
		}
		jsonData, err := json.Marshal(data)
		if err == nil {
			f.App.Event.Emit("convert_progress", string(jsonData))
		}
	}

	return out, nil
}

// 在 convertservice.go 中添加此方法
func (b *ConvertService) ConvertTo16kWave(inputPath string) (string, error) {
	// 获取 FFmpeg 路径
	ffmpegPath := appconfig.AppCfg.FFmpegPath
	if ffmpegPath == "" {
		ffmpegPath = "ffmpeg"
	}

	// 创建临时输出路径 (例如: input_16k.wav)
	outputPath := strings.TrimSuffix(inputPath, filepath.Ext(inputPath)) + "_16k.wav"

	// ffmpeg 参数:
	// -ar 16000: 设置采样率
	// -ac 1: 设置为单声道 (ASR 通常需要单声道)
	// -y: 覆盖已存在文件
	cmd := exec.Command(ffmpegPath, "-i", inputPath, "-ar", "16000", "-ac", "1", "-y", outputPath)

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("FFmpeg 转换失败: %v", err)
	}

	return outputPath, nil
}
