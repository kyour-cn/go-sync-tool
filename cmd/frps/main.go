package main

import (
	"fmt"
	"net"
	"sync"
)

var (
	clients = make(map[string]net.Conn) // 存储客户端连接
	mutex   sync.Mutex                  // 保护 clients 并发访问
)

func main() {
	// 启动控制通道，监听客户端连接
	go startControlServer(":7000")
	// 启动代理端口，监听外部请求
	startProxyServer(":9000")
}

// 控制通道：处理客户端连接
func startControlServer(addr string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("Control server error:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Control server started on", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}
		go handleClient(conn)
	}
}

// 处理客户端注册
func handleClient(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Client read error:", err)
		return
	}
	clientID := string(buffer[:n]) // 客户端发送的标识（如 "web"）
	fmt.Println("Client registered:", clientID)

	mutex.Lock()
	clients[clientID] = conn
	mutex.Unlock()

	// 保持连接
	for {
		_, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Client disconnected:", clientID)
			mutex.Lock()
			delete(clients, clientID)
			mutex.Unlock()
			return
		}
	}
}

// 代理端口：处理外部请求
func startProxyServer(addr string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("Proxy server error:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Proxy server started on", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Proxy accept error:", err)
			continue
		}
		go handleProxy(conn)
	}
}

// 转发外部请求到客户端
func handleProxy(conn net.Conn) {
	defer conn.Close()

	mutex.Lock()
	clientConn, exists := clients["web"] // 假设映射到 "web" 客户端
	mutex.Unlock()

	if !exists {
		fmt.Println("No client available")
		return
	}

	// 双向转发
	go forward(conn, clientConn)
	forward(clientConn, conn)
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
