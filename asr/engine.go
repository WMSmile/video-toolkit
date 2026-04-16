package asr

// AsrConfig 统一配置
type AsrConfig struct {
	ModelPath            string
	TokensPath           string
	VadModelPath         string // 新增：VAD 模型路径 (silero_vad.onnx)
	PunctuationModelPath string // 新增：标点模型路径
	HotwordsPath         string // 新增：热词文件路径
	NumThreads           int
}

// Recognizer 定义了 asrservice 需要调用的行为
type Recognizer interface {
	Process(samples []float32) string
	Flush() string // 新增：强制处理缓冲区中剩余的所有音频
	Close()
}