//go:build linux

package asr

import (
	"fmt"
	"strings"

	sherpa "github.com/k2-fsa/sherpa-onnx-go-linux"
)

type LinuxRecognizer struct {
	recognizer *sherpa.OfflineRecognizer
	vad        *sherpa.VoiceActivityDetector
	punc       *sherpa.OfflinePunctuation // 新增：标点处理器
}

func NewRecognizer(c AsrConfig) (Recognizer, error) {
	// 1. ASR 识别器配置
	config := &sherpa.OfflineRecognizerConfig{}
	config.FeatConfig.SampleRate = 16000
	config.FeatConfig.FeatureDim = 80
	config.ModelConfig.Paraformer.Model = c.ModelPath
	config.ModelConfig.Tokens = c.TokensPath
	config.ModelConfig.NumThreads = c.NumThreads
	config.ModelConfig.Debug = 0

	// --- 新增：传入热词路径 ---
	if c.HotwordsPath != "" {
		config.HotwordsFile = c.HotwordsPath
		// 可选：设置热词分值（分数越高，越容易识别出热词，默认通常是 1.5 到 2.0）
		config.HotwordsScore = 2.0
	}
	// -----------------------

	recognizer := sherpa.NewOfflineRecognizer(config)
	if recognizer == nil {
		return nil, fmt.Errorf("创建识别器失败")
	}

	// 2. VAD 配置
	// 注意：根据你提供的函数签名，需要传入 bufferSizeInSeconds
	// 建议设置为 60 秒左右，这个 buffer 决定了 VAD 内部能缓存多长的未处理音频
	vadConfig := &sherpa.VadModelConfig{
		SileroVad: sherpa.SileroVadModelConfig{
			Model:              c.VadModelPath,
			MinSilenceDuration: 0.5,
			MinSpeechDuration:  0.25,
			Threshold:          0.5,
		},
		SampleRate: 16000,
	}

	// 第二个参数是缓冲区大小（秒）
	vad := sherpa.NewVoiceActivityDetector(vadConfig, 60.0)
	if vad == nil {
		return nil, fmt.Errorf("创建 VAD 失败")
	}

	// 3. 标点恢复配置 (重点在这里)
	puncConfig := &sherpa.OfflinePunctuationConfig{
		Model: sherpa.OfflinePunctuationModelConfig{
			CtTransformer: c.PunctuationModelPath,
			NumThreads:    c.NumThreads,
			Debug:         0,
			Provider:      "cpu",
		},
	}

	punc := sherpa.NewOfflinePunctuation(puncConfig)
	if punc == nil {
		// 如果标点模型加载失败，可以选择记录日志或返回错误
		fmt.Println("警告: 标点模型加载失败，将输出纯文本")
	}

	return &LinuxRecognizer{
		recognizer: recognizer,
		vad:        vad,
		punc:       punc,
	}, nil
}

func (l *LinuxRecognizer) Process(samples []float32) string {
	if l == nil || l.recognizer == nil || l.vad == nil || len(samples) == 0 {
		return ""
	}

	// 1. 将样本喂入 VAD 内部缓冲区
	l.vad.AcceptWaveform(samples)

	var resultTexts []string

	// 2. 只要 VAD 切分出了完整的语音段就进行处理
	for !l.vad.IsEmpty() {
		segment := l.vad.Front()

		// 防御性检查：确保 segment 和 Samples 有效
		if segment.Samples == nil || len(segment.Samples) == 0 {
			l.vad.Pop()
			continue
		}

		// 创建识别流 (OfflineStream 必须在拿到完整 Samples 后再处理)
		stream := sherpa.NewOfflineStream(l.recognizer)
		if stream == nil {
			l.vad.Pop()
			continue
		}

		// 按照官方建议一次性喂入该段所有采样
		stream.AcceptWaveform(16000, segment.Samples)
		l.recognizer.Decode(stream)

		res := stream.GetResult()
		if res.Text != "" {
			fmt.Printf("识别片段结果: %s\n", res.Text) // 调试日志
			resultTexts = append(resultTexts, res.Text)
		}

		// 处理完必须弹出，否则内存会持续增长
		l.vad.Pop()
	}

	rawText := strings.Join(resultTexts, "")

	// 如果有标点模型，则转换
	if l.punc != nil && rawText != "" {
		// 这里的 AddPunct 会把 "你好吗我很好" 变成 "你好吗？我很好。"
		return l.punc.AddPunct(rawText)
	}

	return rawText
}

// Flush 强制处理缓冲区中剩余的所有音频
func (l *LinuxRecognizer) Flush() string {
	if l.vad == nil {
		return ""
	}

	// 某些版本的 sherpa 可以在此处通过输入一组静音采样来强制触发检测
	// 或者如果库支持，直接调用 l.vad.Flush()

	var resultTexts []string

	// 强制读取剩余段落
	for !l.vad.IsEmpty() {
		segment := l.vad.Front()

		stream := sherpa.NewOfflineStream(l.recognizer)
		stream.AcceptWaveform(16000, segment.Samples)
		l.recognizer.Decode(stream)

		res := stream.GetResult()
		if res.Text != "" {
			resultTexts = append(resultTexts, res.Text)
		}
		l.vad.Pop()
	}

	rawText := strings.Join(resultTexts, "")

	// 如果有标点模型，则转换
	if l.punc != nil && rawText != "" {
		// 这里的 AddPunct 会把 "你好吗我很好" 变成 "你好吗？我很好。"
		return l.punc.AddPunct(rawText)
	}

	return rawText
}

func (l *LinuxRecognizer) Close() {
	if l.vad != nil {
		// 释放 VAD C内存
		sherpa.DeleteVoiceActivityDetector(l.vad)
		l.vad = nil
	}
	if l.recognizer != nil {
		sherpa.DeleteOfflineRecognizer(l.recognizer)
		l.recognizer = nil
	}
	if l.punc != nil {
		// 释放标点处理器内存
		sherpa.DeleteOfflinePunc(l.punc)
		l.punc = nil
	}
}