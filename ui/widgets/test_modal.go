package widgets

import (
    "app/ui/chapartheme"
    "gioui.org/layout"
    "gioui.org/op"
    "gioui.org/unit"
    "gioui.org/widget"
    "gioui.org/widget/material"
    "gioui.org/x/component"
)

type TestModal struct {
    visible bool
}

func NewTestModal() *TestModal {
    return &TestModal{}
}
func (c *TestModal) layout(gtx layout.Context, theme *chapartheme.Theme) layout.Dimensions {
    border := widget.Border{
        Color:        theme.TableBorderColor,
        CornerRadius: unit.Dp(4),
        Width:        unit.Dp(2),
    }

    return layout.N.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
        return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
            gtx.Constraints.Max.X = gtx.Dp(600)
            gtx.Constraints.Max.Y = gtx.Dp(300)
            return component.NewModalSheet(component.NewModal()).Layout(gtx, theme.Material(), &component.VisibilityAnimation{}, func(gtx layout.Context) layout.Dimensions {
                return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                    return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
                        layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                            return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                                return material.H3(theme.Material(), "巨蟹科技B2B电商同步系统").Layout(gtx)
                            })
                        }),
                    )
                })
            })
        })
    })
}

func (c *TestModal) SetVisible(visible bool) {
    c.visible = visible
}

func (c *TestModal) Layout(gtx layout.Context, theme *chapartheme.Theme) layout.Dimensions {
    if !c.visible {
        return layout.Dimensions{}
    }

    ops := op.Record(gtx.Ops)
    dims := c.layout(gtx, theme)
    defer op.Defer(gtx.Ops, ops.Stop())

    return dims
}
