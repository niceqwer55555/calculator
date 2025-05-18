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

	// 配置服务器
	conf := &gosocks5.Config{
		AuthMethods: []gosocks5.Authenticator{auth},
		IdleTimeout: config.IdleTimeout,
	}

	server := gosocks5.NewServer(conf)
	return &Server{
		server: server,
		config: config,
		status: "stopped",
	}, nil
}

func (s *Server) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.status == "running" {
		return fmt.Errorf("service already running")
	}

	addr := fmt.Sprintf("0.0.0.0:%d", s.config.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	s.listener = listener
	s.status = "running"

	go func() {
		defer listener.Close()
		fmt.Printf("[SOCKS5] 服务启动于 %s\n", addr)
		if err := s.server.Serve(listener); err != nil {
			s.mu.Lock()
			s.status = "error"
			s.mu.Unlock()
			fmt.Printf("[SOCKS5] 服务错误: %v\n", err)
		}
	}()

	return nil
}

func (s *Server) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.status != "running" {
		return fmt.Errorf("service not running")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.listener.Close(); err != nil {
		return err
	}

	s.status = "stopped"
	fmt.Println("[SOCKS5] 服务已停止")
	return nil
}

func (s *Server) Status() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.status
}
