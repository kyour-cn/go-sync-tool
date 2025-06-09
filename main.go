package main

import (
	"app/internal/domain"
	"app/internal/initialize"
	mainApp "app/ui/app"
	"errors"
	"fmt"
	"gioui.org/app"
	"gioui.org/unit"
	"golang.org/x/sys/windows"
	"log"
	"os"
	"syscall"
)

//go:generate rsrc -ico assets/images/favicon.ico -manifest assets/app.manifest -o main.syso
func main() {

	// 创建一个具有唯一名称的互斥体
	mutexName := "Global\\JxSyncUtilV3"
	mutex, err := windows.CreateMutex(nil, false, windows.StringToUTF16Ptr(mutexName))
	if err != nil {
		fmt.Println("Error creating mutex:", err)
		return
	}
	if errCode := syscall.GetLastError(); errors.Is(errCode, windows.ERROR_ALREADY_EXISTS) {
		fmt.Println("Another instance of the program is already running.")
		return
	} else if errCode != nil {
		fmt.Println("Error checking for existing instance:", errCode)
		return
	}
	// 在程序退出时释放互斥体资源
	defer func(handle windows.Handle) {
		err := windows.CloseHandle(handle)
		if err != nil {
			os.Exit(1)
		}
	}(mutex)

	// 初始化
	initialize.InitApp()

	go func() {
		var w app.Window
		w.Option(
			app.Title(domain.AppName+" ("+domain.Version+")"),
			app.Size(unit.Dp(900), unit.Dp(600)),
		)

		mainUI, err := mainApp.New(&w, domain.Version)
		if err != nil {
			log.Fatal(err)
		}

		if err := mainUI.Run(); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	// 启动托盘图标
	go mainApp.RunNotifyIcon()

	app.Main()
}
