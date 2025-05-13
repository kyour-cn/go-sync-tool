package main

import (
	"fmt"
	"net"
)

func main() {
	// 连接服务端控制通道
	controlConn, err := net.Dial("tcp", "public-server-ip:7000")
	if err != nil {
		fmt.Println("Control connect error:", err)
		return
	}
	defer controlConn.Close()

	// 注册客户端标识
	_, err = controlConn.Write([]byte("web"))
	if err != nil {
		fmt.Println("Register error:", err)
		return
	}
	fmt.Println("Connected to server")

	// 监听服务端转发的请求
	go handleControl(controlConn)

	// 保持运行
	select {}
}

// 处理服务端转发的请求
func handleControl(controlConn net.Conn) {
	for {
		// 接收服务端数据
		buffer := make([]byte, 1024)
		n, err := controlConn.Read(buffer)
		if err != nil {
			fmt.Println("Control read error:", err)
			return
		}

		// 连接内网服务
		localConn, err := net.Dial("tcp", "127.0.0.1:8080") // 内网服务地址
		if err != nil {
			fmt.Println("Local connect error:", err)
			return
		}

		// 发送初始数据到内网服务
		_, err = localConn.Write(buffer[:n])
		if err != nil {
			localConn.Close()
			continue
		}

		// 双向转发
		go forward(localConn, controlConn)
		forward(controlConn, localConn)
	}
}

// 双向数据转发
func forward(src, dst net.Conn) {
	defer src.Close()
	defer dst.Close()
	buffer := make([]byte, 1024)
	for {
		n, err := src.Read(buffer)
		if err != nil {
			return
		}
		_, err = dst.Write(buffer[:n])
		if err != nil {
			return
		}
	}
}
