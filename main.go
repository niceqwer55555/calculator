package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
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
}
