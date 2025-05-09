package console

import (
    "app/internal/global"
    "fmt"
    "time"

    "gioui.org/layout"
    "gioui.org/unit"
    "gioui.org/widget"
    "gioui.org/widget/material"

    "app/ui/apptheme"
)

type Console struct {
    //logs []domain.Log

    list *widget.List

    clearButton *widget.Clickable
}

func New() *Console {
    c := &Console{
        list: &widget.List{
            List: layout.List{
                Axis: layout.Vertical,
            },
        },
        clearButton: &widget.Clickable{},
    }

    //c.logs = make([]domain.Log, 0)
    //
    //c.logs = append(c.logs, domain.Log{
    //    Message: "Welcome to the console!",
    //    Time:    time.Now(),
    //    Level:   "info",
    //})

    // bus.Subscribe(state.LogSubmitted, c.handleIncomingLog)
    return c
}

func (c *Console) logLayout(gtx layout.Context, theme *apptheme.Theme, log *global.Log) layout.Dimensions {
    textColor := theme.Palette.Fg
    switch log.Level {
    case "INFO":
        textColor = apptheme.LightGreen
    case "ERROR":
        textColor = apptheme.LightRed
    case "WARN":
        textColor = apptheme.LightYellow
    }

    return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
        layout.Rigid(func(gtx layout.Context) layout.Dimensions {
            l := material.Label(theme.Material(), theme.TextSize, fmt.Sprintf("[%s] ", log.Time.Format(time.DateTime)))
            l.Color = textColor
            return l.Layout(gtx)
        }),
        layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
            return material.Label(theme.Material(), theme.TextSize, log.Message).Layout(gtx)
        }),
    )
}

func (c *Console) Layout(gtx layout.Context, theme *apptheme.Theme) layout.Dimensions {
    return layout.Inset{
        Top:    unit.Dp(15),
        Left:   unit.Dp(5),
        Bottom: unit.Dp(5),
        Right:  unit.Dp(5),
    }.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
        return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
            layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Spacing: layout.SpaceStart}.Layout(gtx,
                    layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                        if c.clearButton.Clicked(gtx) {
                            //c.logs = make([]domain.Log, 0)
                            global.Logs = make([]global.Log, 0)
                        }
                        return layout.Inset{Bottom: unit.Dp(5)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                            return material.Button(theme.Material(), c.clearButton, "Clear").Layout(gtx)
                        })
                    }),
                )
            }),
            layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
                return widget.Border{
                    Color:        theme.BorderColor,
                    Width:        unit.Dp(1),
                    CornerRadius: unit.Dp(4),
                }.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                    return layout.UniformInset(unit.Dp(5)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                        return material.List(theme.Material(), c.list).Layout(gtx, len(global.Logs), func(gtx layout.Context, i int) layout.Dimensions {
                            return c.logLayout(gtx, theme, &global.Logs[i])
                        })
                    })
                })
            }),
        )
    })
}
