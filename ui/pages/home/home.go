package home

import (
    "app/assets"
    "app/ui/chapartheme"
    "app/ui/widgets"
    "fmt"
    "gioui.org/layout"
    "gioui.org/op/paint"
    "gioui.org/unit"
    "gioui.org/widget"
    "gioui.org/widget/material"
)

type View struct {
    imageOp     paint.ImageOp
    startButton widget.Clickable
}

func New() *View {

    // Load logo image
    data, err := assets.LoadImage("logo.png")
    if err != nil {
        panic(err)
    }

    c := &View{
        imageOp: paint.NewImageOp(data),
    }

    return c
}

func (c *View) Layout(gtx layout.Context, theme *chapartheme.Theme) layout.Dimensions {

    if c.startButton.Clicked(gtx) {
        fmt.Println("Start button clicked")
    }

    return layout.Inset{
        Top: unit.Dp(30), Left: unit.Dp(250), Right: unit.Dp(250),
    }.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
        return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,

            // Logo
            layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                return widget.Image{
                    Src:      c.imageOp,
                    Fit:      widget.Unscaled,
                    Position: layout.Center,
                    Scale:    0.5,
                }.Layout(gtx)
            }),

            layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                    return material.H3(theme.Material(), "巨蟹科技B2B电商同步系统").Layout(gtx)
                })
            }),
            layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                    return material.H5(theme.Material(), "V1.0.0").Layout(gtx)
                })
            }),

            // 显示运行状态
            layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                return layout.Inset{
                    Top: unit.Dp(30),
                }.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                    return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
                        layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                            return material.H6(theme.Material(), "运行状态: 空闲中").Layout(gtx)
                        }),
                        layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                            return material.H6(theme.Material(), "同步任务: 商品资料、商品库存、客户资料、客户地址、订单、发货单、业务员、新客户").Layout(gtx)
                        }),
                    )
                })
            }),

            layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                return layout.Inset{
                    Top: unit.Dp(30),
                }.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                    newBtn := widgets.Button(theme.Material(), &c.startButton, nil, widgets.IconPositionStart, "启动同步")
                    newBtn.Color = theme.ButtonTextColor
                    newBtn.Background = theme.SendButtonBgColor
                    return newBtn.Layout(gtx, theme)
                })
            }),
        )
    })
}
