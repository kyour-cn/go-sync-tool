package main

import (
    "app/internal/domain"
    "app/internal/initialize"
    mainApp "app/ui/app"
    "gioui.org/app"
    "gioui.org/unit"
    "log"
    "os"
)

//go:generate rsrc -ico assets/images/favicon.ico -manifest assets/app.manifest -o main.syso
func main() {

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
