package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// 定义命令行参数
	codebaseName := flag.String("codebase", "", "代码库名称 (例如: troopers)")
	featureName := flag.String("feature", "", "功能名称 (例如: 新建集群)")
	codeListPath := flag.String("code-list", "", "代码文件列表路径 (可选，默认: ./output/{codebase}/full_code_list.txt)")
	outputPath := flag.String("output", "", "输出文件路径 (可选，默认: ./output/{codebase}/full_error_list.md)")
	baseDir := flag.String("base-dir", ".codebase/", "基准目录路径")

	flag.Parse()

	// 如果使用位置参数（兼容原 Python 脚本的调用方式）
	if *codebaseName == "" && len(os.Args) >= 3 {
		*codebaseName = os.Args[1]
		if len(os.Args) >= 3 {
			*featureName = os.Args[2]
		}
	}

	// 验证必需参数
	if *codebaseName == "" {
		fmt.Println("使用方法: go run main.go -codebase <codebase_name> [-feature <feature_name>]")
		fmt.Println("或者: go run main.go <codebase_name> <feature_name>")
		fmt.Println("示例: go run main.go troopers 新建集群")
		os.Exit(1)
	}

	// 设置默认路径
	if *codeListPath == "" {
		*codeListPath = filepath.Join("output", *codebaseName, "full_code_list.txt")
	}
	if *outputPath == "" {
		*outputPath = filepath.Join("output", *codebaseName, "full_error_list.md")
	}

	// 确保输出目录存在
	outputDir := filepath.Dir(*outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("创建输出目录失败: %v\n", err)
		os.Exit(1)
	}

	// 执行错误搜索
	processor := NewFileProcessor(*baseDir)
	searcher := NewErrorSearcher()

	errorCount, err := processor.ProcessFileList(*codeListPath, *outputPath, searcher)
	if err != nil {
		fmt.Printf("处理文件列表时出错: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("总共找到 %d 条错误信息\n", errorCount)
	fmt.Printf("结果已保存到 %s\n", *outputPath)
}

