package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/fyne-io/calculator/oneapi"
)

func main() {
	a := app.New()
	w := a.NewWindow("One-API控制器")

	// 配置参数
	cfg := oneapi.Config{
		Port:     8080,
		APIToken: "123456",
	}

	// UI组件
	statusLabel := widget.NewLabel("服务状态: 未运行")
	startBtn := widget.NewButton("启动服务", nil)
	stopBtn := widget.NewButton("停止服务", nil)
	stopBtn.Disable()

	// 按钮逻辑
	startBtn.OnTapped = func() {
		if err := oneapi.Start(cfg); err != nil {
			statusLabel.SetText("启动失败: " + err.Error())
		} else {
			statusLabel.SetText("服务状态: 运行中 (端口 8080)")
			startBtn.Disable()
			stopBtn.Enable()
		}
	}

	stopBtn.OnTapped = func() {
		if err := oneapi.Stop(); err != nil {
			statusLabel.SetText("停止失败: " + err.Error())
		} else {
			statusLabel.SetText("服务状态: 已停止")
			startBtn.Enable()
			stopBtn.Disable()
		}
	}

	// 布局
	controls := container.NewGridWithColumns(2, startBtn, stopBtn)
	content := container.NewVBox(
		statusLabel,
		layout.NewSpacer(),
		controls,
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(300, 300))
	w.ShowAndRun()
}
