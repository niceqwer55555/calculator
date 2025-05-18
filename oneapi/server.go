package oneapi

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var (
	server *http.Server
	mu     sync.Mutex
)

// 服务配置（可扩展）
type Config struct {
	Port     int
	APIToken string
}

// Start 启动服务（非阻塞）
func Start(cfg Config) error {
	mu.Lock()
	defer mu.Unlock()

	if server != nil {
		return fmt.Errorf("service already running")
	}

	mux := http.NewServeMux()
	setupRoutes(mux, cfg.APIToken) // 自定义路由

	server = &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", cfg.Port),
		Handler: mux,
	}

	go func() {
		fmt.Printf("One-API服务启动，端口 %d\n", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("服务错误: %v\n", err)
		}
	}()

	return nil
}

// Stop 停止服务
func Stop() error {
	mu.Lock()
	defer mu.Unlock()

	if server == nil {
		return fmt.Errorf("service not running")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("停止失败: %v", err)
	}

	server = nil
	fmt.Println("服务已停止")
	return nil
}

// 示例路由配置
func setupRoutes(mux *http.ServeMux, token string) {
	mux.HandleFunc("/api/v1/data", authMiddleware(token, handleData))
}

func authMiddleware(token string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-API-Token") != token {
			http.Error(w, "未授权", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

func handleData(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `{"status": "ok"}`)
}
