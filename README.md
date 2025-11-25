# Go 错误信息搜索工具

这是 `3.grep_error_list.py` 的 Go 语言实现版本，用于在 Go 代码库中搜索和提取错误信息。

## 功能

- 从文件列表中读取 Go 文件路径
- 在每个文件中搜索各种错误处理模式（log.Error, fmt.Errorf, errors.New 等）
- 提取错误消息字符串
- 将结果以 Markdown 表格格式输出

## 文件结构

- `main.go` - 主程序入口，处理命令行参数
- `error_search.go` - 错误模式搜索和提取逻辑
- `file_processor.go` - 文件读取和路径处理

## 使用方法

### 基本用法

```bash
# 使用位置参数（兼容原 Python 脚本）
go run main.go troopers 新建集群

# 使用命令行标志
go run main.go -codebase troopers -feature "新建集群"

# 指定自定义路径
go run main.go -codebase troopers \
  -code-list ./custom/code_list.txt \
  -output ./custom/error_list.md \
  -base-dir .codebase/
```

### 编译后使用

```bash
# 编译
go build -o grep_error_list

# 运行
./grep_error_list -codebase troopers
```

## 命令行参数

- `-codebase` - 代码库名称（必需）
- `-feature` - 功能名称（可选）
- `-code-list` - 代码文件列表路径（默认: `./output/{codebase}/full_code_list.txt`）
- `-output` - 输出文件路径（默认: `./output/{codebase}/full_error_list.md`）
- `-base-dir` - 基准目录路径（默认: `.codebase/`）

## 支持的错误模式

程序会搜索以下类型的错误处理模式：

- 标准日志库：`log.Error`, `log.Warn`, `log.Fatal`
- Kubernetes klog：`klog.Error`, `klog.Warning`, `klog.Fatal`
- 标准错误处理：`fmt.Errorf`, `errors.New`, `errors.Wrap`
- 自定义日志器：`logger.Error`, `glog.Error`
- Panic 模式：`panic`, `log.Panic`
- HTTP 错误响应：`http.Error`
- Kubernetes 特定：`field.Error`, `apierrors.New`

## 输出格式

输出为 Markdown 表格格式：

```markdown
# 相关错误信息汇总

| 报错日志 | 文件路径 | 行号 |
| -------- | -------- | ---- |
| error message here | path/to/file.go | 123 |
```

## 示例

```bash
# 搜索 gin 代码库的错误信息
go run main.go gin

# 输出会保存到: ./output/gin/full_error_list.md
```

