package main

import (
    "flag"
    "log"
    _ "net/http/pprof"
    "os"

    mainApp "app/ui/app"
    "gioui.org/app"
    "gioui.org/unit"
)

var version = "v1.0.0"

func main() {
    flag.Parse()

    go func() {
        var w app.Window
        w.Option(app.Title("巨蟹ERP同步工具(B2B) "+version), app.Size(unit.Dp(1200), unit.Dp(800)))

        mainUI, err := mainApp.New(&w, version)
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
