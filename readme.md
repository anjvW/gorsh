# Go TCP Shell Server

这是一个用 Go 语言编写反弹shell时获取交互式Shell，允许用户通过 TCP 连接与远程 bash 会话进行交互。能够处理 `Ctrl + C` 和 `Ctrl + D` 信号。

## 特性

- 监听指定的 TCP 端口（默认为 1234）。
- 支持通过命令行参数指定监听端口。
- 启动一个交互式 bash 会话。
- 动态设置终端的行数和列数。
- 处理 `Ctrl + C` 和 `Ctrl + D` 信号。

## 安装

确保您已经安装了 Go 语言环境。然后，您可以克隆该项目，运行，或下载release版本。

## 使用

在终端中运行以下命令以启动服务器：

```bash
go run main.go [port]
```

其中 `[port]` 是可选的，您可以指定要监听的端口。如果未指定，默认将使用 1234 端口。

例如，要在 1234 端口上运行服务器：

```bash
go run main.go 1234
```

## 连接到服务器

您可以使用 `netcat` 或其他 TCP 客户端连接到服务器：

```bash
nc 127.0.0.1 1234
```


## 终止程序

- 按 `Ctrl + D` 发送 EOF 信号并关闭连接。

## 依赖

该项目依赖于 `golang.org/x/term` 包。您可以通过以下命令安装它：

```bash
go get golang.org/x/term
```

## 许可证

该项目使用 MIT 许可证。有关详细信息，请参阅 [LICENSE](LICENSE) 文件。
