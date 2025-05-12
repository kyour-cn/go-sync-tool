//go:build windows

package app

import (
	"app/assets"
	"app/internal/global"
	"context"
	"github.com/go-gourd/gourd/event"
	"log"
)

var (
	// 托盘图标
	ni *walk.NotifyIcon
)

func RunNotifyIcon() {
	w, err := walk.NewMainWindow()
	if err != nil {
		log.Fatal(err)
	}
	ni, err = walk.NewNotifyIcon(w)
	if err != nil {
		log.Fatal(err)
	}

	// 图标设置
	ico, err := assets.LoadImage("icon.png")
	resource, err := walk.NewIconFromImageForDPI(ico, 32)
	if err != nil {
		return
	}

	_ = ni.SetIcon(resource)
	if err := ni.SetToolTip("点击打开主界面"); err != nil {
		log.Fatal(err)
	}
	// 设置托盘图标按钮事件
	ni.MouseDown().Attach(func(x, y int, button walk.MouseButton) {
		global.State.HideWindow = false
	})

	_ = ni.SetVisible(true)

	// 监听事件 - 显示通知
	event.Listen("notify_show", func(ctx context.Context) {
		title := ctx.Value("notify_title")
		info := ctx.Value("notify_info")
		if err := ni.ShowInfo(title.(string), info.(string)); err != nil {
			log.Fatal(err)
		}
	})

	// Run the message loop.
	w.Run()
}
