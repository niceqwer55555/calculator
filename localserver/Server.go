package localserver

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

// 启动服务（非阻塞）
func Start() {
	mu.Lock()
	defer mu.Unlock()

	if server != nil {
		return // 防止重复启动
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from Go server!")
	})

	server = &http.Server{
		Addr:    "0.0.0.0:8080", // Android 必须绑定到 0.0.0.0
		Handler: mux,
	}

	go func() {
		fmt.Println("Server starting on :8080")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Printf("Server error: %v\n", err)
		}
	}()
}

// 停止服务
func Stop() {
	mu.Lock()
	defer mu.Unlock()

	if server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			fmt.Printf("Server shutdown error: %v\n", err)
		}
		server = nil
		fmt.Println("Server stopped")
	}
}
