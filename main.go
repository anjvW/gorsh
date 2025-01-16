package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"golang.org/x/term"
)

func main() {
	// 获取监听端口，默认为1234
	port := "1234"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	// 监听指定端口
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Listening on port " + port + "...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// 启动一个新的goroutine来处理连接
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	var err error // 声明err变量

	// 设置终端为原始模式
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Error setting terminal to raw mode:", err)
		return
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState) // 恢复终端状态

	// 启动一个新的bash会话
	_, err = conn.Write([]byte("python -c 'import pty; pty.spawn(\"/bin/bash\")'\n"))
	if err != nil {
		fmt.Println("错误:", err)
		os.Exit(1)
	}

	// 等待一段时间，确保远程服务器准备好
	time.Sleep(1 * time.Second)

	// 获取终端的行数和列数
	var width, height int
	if runtime.GOOS == "windows" {
		// 在 Windows 上使用默认值
		width, height = 80, 24
	} else {
		width, height, err = term.GetSize(int(os.Stdin.Fd()))
		if err != nil {
			fmt.Println("Error getting terminal size:", err)
			return
		}
	}

	// 发送初始化命令
	commands := []string{
		"reset\n",
		"export SHELL=bash\n",
		"export TERM=xterm-256color\n",
		fmt.Sprintf("stty rows %d columns %d\n", height, width), // 使用动态计算的行列数
	}

	for _, cmd := range commands {
		_, err = conn.Write([]byte(cmd))
		if err != nil {
			fmt.Println("错误:", err)
			os.Exit(1)
		}
	}

	// 启动一个goroutine来读取远程服务器的输出
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				fmt.Println("从远程服务器读取数据失败:", err)
				term.Restore(int(os.Stdin.Fd()), oldState) // 恢复终端状态
				os.Exit(1)                                 // 退出程序
			}
			fmt.Print(string(buf[:n])) // 打印远程服务器的输出
		}
	}()

	// 捕获Ctrl+C信号
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)

	// 从标准输入读取数据并发送到远程服务器
	buf := make([]byte, 1) // 每次读取一个字节
	for {
		select {
		case <-signalChan:
			// 捕获到Ctrl+C信号，发送给远程服务器
			_, err := conn.Write([]byte{3}) // 发送Ctrl+C的ASCII码（3）
			if err != nil {
				fmt.Println("Error sending Ctrl+C to remote:", err)
				os.Exit(1)
			}
		default:
			_, err := os.Stdin.Read(buf)
			if err != nil {
				fmt.Println("Error reading from stdin:", err)
				term.Restore(int(os.Stdin.Fd()), oldState) // 恢复终端状态
				os.Exit(1)                                 // 退出程序
			}

			// 检查是否为EOF（Ctrl + D）
			if buf[0] == 4 { // 4 是 Ctrl + D 的 ASCII 值
				fmt.Println("退出程序，关闭连接！")
				conn.Close()
				term.Restore(int(os.Stdin.Fd()), oldState)
				os.Exit(0)
			}

			// 发送输入到远程连接
			_, err = conn.Write(buf)
			if err != nil {
				fmt.Println("Error sending data to remote:", err)
				term.Restore(int(os.Stdin.Fd()), oldState) // 恢复终端状态
				os.Exit(1)                                 // 退出程序
			}
		}
	}
}
