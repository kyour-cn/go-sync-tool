package main

import (
    "app/internal/global"
    "flag"
    "log"
    "os"

    mainApp "app/ui/app"
    "gioui.org/app"
    "gioui.org/unit"
)

func main() {
    flag.Parse()

    go func() {
        var w app.Window
        w.Option(app.Title(global.AppName+" ("+global.Version+")"), app.Size(unit.Dp(1200), unit.Dp(800)))

        mainUI, err := mainApp.New(&w, global.Version)
        if err != nil {
            log.Fatal(err)
        }

        if err := mainUI.Run(); err != nil {
            log.Fatal(err)
        }
        os.Exit(0)
    }()

    app.Main()
}
