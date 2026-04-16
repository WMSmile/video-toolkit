package config

import (
	"fmt"
	"os"
	"path/filepath"

	"runtime"

	"gopkg.in/yaml.v3"
)

// ===================== 【全局配置变量】全程序通用 =====================
type AppConfig struct {
	DebugMode            bool   `yaml:"debug_mode"`
	FFmpegPath           string `yaml:"ffmpeg_path"`
	FFprobePath          string `yaml:"ffprobe_path"`
	Theme                string `yaml:"theme"`
	PrimaryColor         string `yaml:"primary_color"`
	ModelPath            string `yaml:"model_path"`
	TokensPath           string `yaml:"tokens_path"`
	VadModelPath         string `yaml:"vad_model_path"`         // VAD 模型路径 (silero_vad.onnx)
	PunctuationModelPath string `yaml:"punctuation_model_path"` // 标点模型路径
	HotwordsPath         string `yaml:"hotwords_path"`          // 热词文件路径
}

// 全局配置实例
var AppCfg AppConfig

// InitConfig 初始化配置
func InitConfig() error {
	// 设置默认值
	AppCfg = AppConfig{
		DebugMode:            false,
		FFmpegPath:           "ffmpeg", // 默认使用系统 PATH 中的 ffmpeg
		FFprobePath:          "ffprobe",
		ModelPath:            "./resources/models/sherpa-onnx-paraformer-zh-2023-09-14/model.int8.onnx", // 实际项目中需要设置正确的模型路径
		TokensPath:           "./resources/models/sherpa-onnx-paraformer-zh-2023-09-14/tokens.txt",
		VadModelPath:         "./resources/models/silero_vad/silero_vad_v5.onnx",                                                       // 实际项目中需要设置正确的词表路径
		PunctuationModelPath: "./resources/models/csukuangfj/sherpa-onnx-punct-ct-transformer-zh-en-vocab272727-2024-04-12/model.onnx", // 实际项目中需要设置正确的词表路径
		HotwordsPath:         "",
	}

	// 加载配置文件
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	// 检查配置文件是否存在
	if _, err := os.Stat(configPath); err == nil {
		// 读取配置文件
		data, err := os.ReadFile(configPath)
		if err != nil {
			return err
		}

		// 解析 YAML
		if err := yaml.Unmarshal(data, &AppCfg); err != nil {
			return err
		}
	}

	return nil
}

// getConfigPath 获取配置文件路径
func getConfigPath() (string, error) {
	// 获取用户配置目录
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	fmt.Println("用户配置目录:", configDir)
	// 创建应用配置目录
	appConfigDir := filepath.Join(configDir, "video-toolkit")
	if err := os.MkdirAll(appConfigDir, 0755); err != nil {
		return "", err
	}

	// 返回配置文件路径
	return filepath.Join(appConfigDir, "config.yaml"), nil
}

func GetPlatform() string {
	return runtime.GOOS // macOS 会返回 "darwin", Windows 返回 "windows"
}

// SaveConfig 保存配置到文件
func SaveConfig() error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	// 序列化配置为 YAML
	data, err := yaml.Marshal(AppCfg)
	if err != nil {
		return err
	}

	// 写入文件
	return os.WriteFile(configPath, data, 0644)
}

/**

// 如何使用配置

// 修改配置
AppCfg.DebugMode = true
AppCfg.FFmpegPath = "/usr/local/bin/ffmpeg"

// 保存配置
if err := SaveConfig(); err != nil {
    log.Printf("保存配置失败: %v", err)
}

// 读取 FFmpeg 路径
ffmpegPath := AppCfg.FFmpegPath

// 检查调试模式
if AppCfg.DebugMode {
    log.Println("调试模式已启用")
}

3. 配置文件示例 (config.yaml)
debug_mode: true
ffmpeg_path: /usr/local/bin/ffmpeg


*/
