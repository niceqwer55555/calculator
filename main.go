package main

import (
	"calculator/socks5server"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
<<<<<<< HEAD
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/fyne-io/calculator/localserver" // 替换为你的包路径
	"os"
)

func main() {
	a := app.New()
	os.Setenv("FYNE_FONT", "小米兰亭字体.ttf")
	a.SetIcon(theme.FyneLogo())

	w := a.NewWindow("demo")

	// 启动服务按钮
	startBtn := widget.NewButton("Start Server", func() {
		localserver.Start()
	})

	// 停止服务按钮
	stopBtn := widget.NewButton("Stop Server", func() {
		localserver.Stop()
	})

	// 布局
	grid := container.NewGridWithColumns(2, startBtn, stopBtn)
	w.SetContent(grid)
	w.Resize(fyne.NewSize(300, 300))
	w.ShowAndRun()
=======
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"log"
	"time"
)

type ServiceController struct {
	server   *socks5server.Server
	startBtn *widget.Button
	stopBtn  *widget.Button
	status   *widget.Label
}

func NewServiceController() *ServiceController {
	cfg := socks5server.Config{
		Port:        1080,
		Username:    "admin",
		Password:    "123456",
		IdleTimeout: 30 * time.Minute,
	}
	server, _ := socks5server.New(cfg)

	return &ServiceController{
		server: server,
	}
}

func (sc *ServiceController) createUI() fyne.CanvasObject {
	sc.startBtn = widget.NewButtonWithIcon("启动代理", theme.MediaPlayIcon(), nil)
	sc.stopBtn = widget.NewButtonWithIcon("停止代理", theme.MediaStopIcon(), nil)
	sc.status = widget.NewLabel("服务状态: 已停止")

	sc.startBtn.OnTapped = func() {
		go func() {
			if err := sc.server.Start(); err != nil {
				log.Println("启动失败:", err)
			}
			sc.updateUI()
		}()
	}

	sc.stopBtn.OnTapped = func() {
		go func() {
			if err := sc.server.Stop(); err != nil {
				log.Println("停止失败:", err)
			}
			sc.updateUI()
		}()
	}

	sc.updateUI()
	return container.NewVBox(
		container.NewHBox(
			sc.startBtn,
			sc.stopBtn,
			layout.NewSpacer(),
			sc.status,
		),
		widget.NewSeparator(),
	)
}

func (sc *ServiceController) updateUI() {
	switch sc.server.Status() {
	case "running":
		sc.status.SetText("服务状态: 运行中 (端口 1080)")
		sc.startBtn.Disable()
		sc.stopBtn.Enable()
	case "stopped":
		sc.status.SetText("服务状态: 已停止")
		sc.startBtn.Enable()
		sc.stopBtn.Disable()
	case "error":
		sc.status.SetText("服务状态: 异常")
		sc.startBtn.Enable()
		sc.stopBtn.Disable()
	}
}

func main() {
	app := app.New()
	window := app.NewWindow("SOCKS5代理控制器")

	controller := NewServiceController()
	content := controller.createUI()

	window.SetContent(content)
	window.Resize(fyne.NewSize(400, 200))
	window.ShowAndRun()
>>>>>>> c36bb4c0fe53c1b358f0223e203cd7b1e7c885ad
}
