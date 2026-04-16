package services

import (
	"fmt"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type FileService struct{}

// SelectFile 适配 alpha.74 版本
func (s *FileService) SelectFile(title string, filterName string, patterns string) (string, error) {
	// 1. 初始化选项
	opts := &application.OpenFileDialogOptions{
		Title:                title,
		CanChooseFiles:       true,
		CanChooseDirectories: false,
	}
	if patterns != "" {
		opts.Filters = []application.FileFilter{
			{
				DisplayName: filterName,
				Pattern:     patterns,
			},
		}
	}
	fmt.Printf("收到的过滤器规则: %s\n", patterns)

	// 2. 关键修复：在 alpha.74 中，直接初始化结构体
	// 注意：这里使用 {} 初始化结构体，而不是把它当函数调用
	dialog := application.OpenFileDialogStruct{}

	// 3. 设置选项并运行
	// 使用单选模式 PromptForSingleSelection
	dialog.SetOptions(opts)
	return dialog.PromptForSingleSelection()
}

// SelectDirectory 适配 alpha.74 版本
func (s *FileService) SelectDirectory(title string) (string, error) {
	opts := &application.OpenFileDialogOptions{
		Title:                title,
		CanChooseDirectories: true,
		CanChooseFiles:       false,
	}

	dialog := application.OpenFileDialogStruct{}
	dialog.SetOptions(opts)

	return dialog.PromptForSingleSelection()
}

// SelectMultipleFiles 适配 alpha.74 版本
func (s *FileService) SelectMultipleFiles(title string, filterName string, patterns string) ([]string, error) {
	opts := &application.OpenFileDialogOptions{
		Title:                title,
		CanChooseDirectories: false,
		CanChooseFiles:       true,
	}
	if patterns != "" {
		opts.Filters = []application.FileFilter{
			{
				DisplayName: filterName,
				Pattern:     patterns,
			},
		}
	}

	fmt.Printf("收到的过滤器规则: %s\n", patterns)

	dialog := application.OpenFileDialogStruct{}
	dialog.SetOptions(opts)

	return dialog.PromptForMultipleSelection()
}
