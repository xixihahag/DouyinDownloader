package main

import (
	"douyin/internal/theme"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {

	myApp := app.New()
	myApp.Settings().SetTheme(&theme.MyTheme{}) // 设置主题支持中文
	myWindow := myApp.NewWindow("Entry Widget")

	input := widget.NewEntry()
	input.SetPlaceHolder("Enter url...")

	// log 日志
	logContent := container.NewVBox()
	scroll := container.NewScroll(logContent)

	dy := NewDouyin(logContent, scroll)

	content := container.NewVBox(input, widget.NewButton("Download", func() {
		if err := dy.Get(input.Text); err != nil {
			log.Println(err)
		}
	}))

	// logContent.Add(
	// 	widget.NewLabel("this is add label"),
	// )
	// scroll.ScrollToBottom()

	myWindow.SetContent(container.NewVSplit(content, scroll))
	myWindow.Resize(fyne.NewSize(500, 500))
	myWindow.ShowAndRun()

	// 单个视频
	// err := dy.Get("3.33 WZm:/  https://v.douyin.com/6oGRLEE/ 复制此链接，打开Dou音搜索，直接观看视频！")
	// 主页视频
	// err := dy.Get("长按复制此条消息，打开抖音搜索，查看TA的更多作品。 https://v.douyin.com/6otwUTD/")
}
