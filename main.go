package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"log"
	"time"

	"github.com/fyne-io/calculator/socks5server"
)

type AppUI struct {
	window     fyne.Window
	server     *socks5server.Server
	statusText *widget.Label
	startBtn   *widget.Button
	stopBtn    *widget.Button
}

func NewAppUI() *AppUI {
	cfg := socks5server.Config{
		Port:        1080,
		Username:    "admin",
		Password:    "123456",
		IdleTimeout: 30 * time.Minute,
	}

	server, _ := socks5server.New(cfg)

	return &AppUI{
		server: server,
	}
}

func (ui *AppUI) Run() {
	a := app.New()
	ui.window = a.NewWindow("SOCKS5 代理控制器")

	// 界面组件
	ui.statusText = widget.NewLabel("服务状态: 已停止")
	ui.startBtn = widget.NewButtonWithIcon("启动服务", theme.MediaPlayIcon(), nil)
	ui.stopBtn = widget.NewButtonWithIcon("停止服务", theme.MediaStopIcon(), nil)

	// 按钮样式
	ui.startBtn.Importance = widget.HighImportance
	ui.stopBtn.Importance = widget.DangerImportance
	ui.stopBtn.Disable()

	// 事件绑定
	ui.startBtn.OnTapped = func() {
		go func() {
			if err := ui.server.Start(); err != nil {
				log.Printf("启动失败: %v", err)
			}
			ui.updateUI()
		}()
	}

	ui.stopBtn.OnTapped = func() {
		go func() {
			if err := ui.server.Stop(); err != nil {
				log.Printf("停止失败: %v", err)
			}
			ui.updateUI()
		}()
	}

	// 布局
	content := container.NewVBox(
		container.NewCenter(ui.statusText),
		layout.NewSpacer(),
		container.NewGridWithColumns(2, ui.startBtn, ui.stopBtn),
	)

	ui.window.SetContent(content)
	ui.window.Resize(fyne.NewSize(300, 150))
	ui.window.ShowAndRun()
}

func (ui *AppUI) updateUI() {
	switch ui.server.Status() {
	case "running":
		ui.statusText.SetText(fmt.Sprintf("服务运行中 (端口 %d)", ui.server.GetPort()))
		ui.startBtn.Disable()
		ui.stopBtn.Enable()
	case "stopped":
		ui.statusText.SetText("服务已停止")
		ui.startBtn.Enable()
		ui.stopBtn.Disable()
	case "error":
		ui.statusText.SetText("服务异常")
		ui.startBtn.Enable()
		ui.stopBtn.Disable()
	}
}

func main() {
	NewAppUI().Run()
}
