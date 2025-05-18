package socks5server

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/ginuerzh/gosocks5"
)

type Config struct {
	Port        int
	Username    string
	Password    string
	IdleTimeout time.Duration
}

type Server struct {
	server   *gosocks5.Server
	listener net.Listener
	config   Config
	mu       sync.Mutex
	status   string // stopped, running, error
}

func New(config Config) (*Server, error) {
	// 创建认证器
	creds := gosocks5.StaticCredentials{
		config.Username: config.Password,
	}
	auth := gosocks5.UserPassAuthenticator{Credentials: creds}

	// 配置服务器参数
	conf := &gosocks5.Config{
		AuthMethods: []gosocks5.Authenticator{auth},
		IdleTimeout: config.IdleTimeout,
	}

	return &Server{
		server: gosocks5.NewServer(conf),
		config: config,
		status: "stopped",
	}, nil
}

// 添加导出方法
func (s *Server) GetPort() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.config.Port
}

// Start 启动服务（非阻塞）
func (s *Server) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.status == "running" {
		return fmt.Errorf("服务已在运行")
	}

	addr := fmt.Sprintf("0.0.0.0:%d", s.config.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("端口监听失败: %v", err)
	}

	s.listener = listener
	s.status = "running"

	go func() {
		defer listener.Close()
		fmt.Printf("✅ SOCKS5 服务启动于 %s\n", addr)
		if err := s.server.Serve(listener); err != nil {
			s.mu.Lock()
			s.status = "error"
			s.mu.Unlock()
			fmt.Printf("❌ 服务异常: %v\n", err)
		}
	}()

	return nil
}

// Stop 停止服务
func (s *Server) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.status != "running" {
		return fmt.Errorf("服务未运行")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.listener.Close(); err != nil {
		return fmt.Errorf("停止服务失败: %v", err)
	}

	s.status = "stopped"
	fmt.Println("🛑 服务已停止")
	return nil
}

// Status 获取服务状态
func (s *Server) Status() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.status
}
