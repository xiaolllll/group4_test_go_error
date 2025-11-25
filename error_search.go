package main

import (
	"regexp"
	"strings"
)

// ErrorSearcher 负责搜索和提取错误信息
type ErrorSearcher struct {
	errorPatterns        []*regexp.Regexp
	errorMessagePatterns []*regexp.Regexp
}

// NewErrorSearcher 创建新的错误搜索器
func NewErrorSearcher() *ErrorSearcher {
	searcher := &ErrorSearcher{}

	// 初始化错误模式（基于 Kubernetes 项目实际错误处理模式）
	errorPatternStrings := []string{
		// 标准日志库错误模式
		`log\.Errorf?\([^)]*"[^"]*"[^)]*\)`,     // log.Errorf/Error("...", ...)
		`log\.Warnf?\([^)]*"[^"]*"[^)]*\)`,     // log.Warnf/Warn("...", ...)
		`log\.Fatalf?\([^)]*"[^"]*"[^)]*\)`,    // log.Fatalf/Fatal("...", ...)

		// Kubernetes klog日志库
		`klog\.Errorf?\([^)]*"[^"]*"[^)]*\)`,   // klog.Errorf/Error("...", ...)
		`klog\.Warningf?\([^)]*"[^"]*"[^)]*\)`, // klog.Warningf/Warning("...", ...)
		`klog\.Fatalf?\([^)]*"[^"]*"[^)]*\)`,   // klog.Fatalf/Fatal("...", ...)
		`klog\.Infof?\([^)]*"[Ee]rror[^"]*"[^)]*\)`, // klog.Info包含error

		// 标准错误处理
		`fmt\.Errorf\([^)]*"[^"]*"[^)]*\)`,      // fmt.Errorf("...", ...)
		`errors\.New\([^)]*"[^"]*"[^)]*\)`,      // errors.New("...", ...)
		`errors\.Wrapf?\([^)]*"[^"]*"[^)]*\)`,   // errors.Wrap/Wrapf("...", ...)

		// 上下文错误包装
		`fmt\.Errorf\([^)]*"[^"]*"[^)]*%w[^)]*\)`, // fmt.Errorf("... %w", err)
		`errors\.Wrap\([^)]*"[^"]*"[^)]*\)`,      // errors.Wrap(err, "...")

		// 返回错误消息
		`return\s+fmt\.Errorf\([^)]*"[^"]*"[^)]*\)`, // return fmt.Errorf("...")
		`return\s+errors\.New\([^)]*"[^"]*"[^)]*\)`, // return errors.New("...")
		`return\s+errors\.Wrapf?\([^)]*"[^"]*"[^)]*\)`, // return errors.Wrap(...)

		// 自定义日志器
		`logger\.Errorf?\([^)]*"[^"]*"[^)]*\)`,   // logger.Errorf/Error("...", ...)
		`glog\.Errorf?\([^)]*"[^"]*"[^)]*\)`,    // glog.Errorf/Error("...", ...)
		`glog\.Warningf?\([^)]*"[^"]*"[^)]*\)`,  // glog.Warningf/Warning("...", ...)

		// Panic模式
		`panic\([^)]*"[^"]*"[^)]*\)`,            // panic("...")
		`log\.Panicf?\([^)]*"[^"]*"[^)]*\)`,     // log.Panicf/Panic("...", ...)

		// 事件记录
		`Eventf?\([^)]*"[^"]*"[^"]*"[^"]*"[^)]*\)`, // Event/Eventf(...)
		`Recorder\.Eventf?\([^)]*"[^"]*"[^"]*"[^"]*"[^)]*\)`, // recorder.Event(...)

		// HTTP错误响应
		`http\.Error\([^)]*"[^"]*"[^)]*\)`,       // http.Error(w, "...", ...)

		// Kubernetes特定模式
		`field\.Error\([^)]*"[^"]*"[^)]*\)`,      // field.Error(...)
		`apierrors\.New[^)]*\([^)]*"[^"]*"[^)]*\)`, // apierrors.New...("...")

		// 条件错误消息
		`if\s+.*\{\s*return\s+fmt\.Errorf\([^)]*"[^"]*"[^)]*\)`, // if ... { return fmt.Errorf(...) }
		`if\s+.*\{\s*return\s+errors\.New\([^)]*"[^"]*"[^)]*\)`, // if ... { return errors.New(...) }
	}

	// 编译正则表达式
	searcher.errorPatterns = make([]*regexp.Regexp, 0, len(errorPatternStrings))
	for _, pattern := range errorPatternStrings {
		re, err := regexp.Compile(pattern)
		if err != nil {
			// 如果编译失败，跳过该模式
			continue
		}
		searcher.errorPatterns = append(searcher.errorPatterns, re)
	}

	// 匹配错误消息字符串的模式
	searcher.errorMessagePatterns = []*regexp.Regexp{
		regexp.MustCompile(`"([^"]*)"`), // 提取双引号中的错误消息
	}

	return searcher
}

// ErrorInfo 表示一个错误信息
type ErrorInfo struct {
	Index        int
	ErrorMessage string
	FilePath     string
	LineNum      int
	FullLine     string
}

// SearchErrors 在文件内容中搜索错误信息
func (es *ErrorSearcher) SearchErrors(content string, filePath string) []ErrorInfo {
	var errors []ErrorInfo
	lines := strings.Split(content, "\n")

	for lineNum, line := range lines {
		lineContent := strings.TrimSpace(line)
		if lineContent == "" {
			continue
		}

		// 检查是否匹配错误模式
		isError := false
		for _, pattern := range es.errorPatterns {
			if pattern.MatchString(lineContent) {
				isError = true
				break
			}
		}

		// 如果是错误信息
		if isError {
			// 提取错误消息字符串
			errorMessage := es.extractErrorMessage(lineContent)

			// 只有当错误消息不为空时才记录
			if errorMessage != "" && strings.TrimSpace(errorMessage) != "" {
				errors = append(errors, ErrorInfo{
					ErrorMessage: errorMessage,
					FilePath:     filePath,
					LineNum:      lineNum + 1, // 行号从1开始
					FullLine:     lineContent,
				})
			}
		}
	}

	return errors
}

// extractErrorMessage 从行内容中提取错误消息
func (es *ErrorSearcher) extractErrorMessage(lineContent string) string {
	for _, pattern := range es.errorMessagePatterns {
		matches := pattern.FindStringSubmatch(lineContent)
		if len(matches) > 1 {
			return matches[1]
		}
	}
	return ""
}

