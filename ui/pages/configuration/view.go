package configuration

import (
    "app/ui/chapartheme"
    "app/ui/widgets"
    "gioui.org/layout"
    "gioui.org/unit"
    "gioui.org/widget"
    "gioui.org/widget/material"
)

type View struct {
    targetEnvEditor *widget.Editor
}

func New(theme *chapartheme.Theme) *View {

    c := &View{}
    c.targetEnvEditor = &widget.Editor{SingleLine: true}

    return c
}

// TODO 待实现
func (v *View) Layout(gtx layout.Context, theme *chapartheme.Theme) layout.Dimensions {
    return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,

        layout.Rigid(func(gtx layout.Context) layout.Dimensions {
            return layout.Inset{Left: unit.Dp(1), Right: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                gtx.Constraints.Min.X = gtx.Dp(unit.Dp(100))
                gtx.Constraints.Max.X = gtx.Dp(unit.Dp(100))
                editor := material.Editor(theme.Material(), v.targetEnvEditor, "Target Key")
                editor.SelectionColor = theme.TextSelectionColor
                return editor.Layout(gtx)
            })
        }),
        layout.Rigid(func(gtx layout.Context) layout.Dimensions {
            return widgets.DrawLine(gtx, theme.TableBorderColor, unit.Dp(35), unit.Dp(1))
        }),
        layout.Rigid(func(gtx layout.Context) layout.Dimensions {
            return layout.Inset{Left: unit.Dp(8), Right: unit.Dp(1)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                gtx.Constraints.Min.X = gtx.Dp(unit.Dp(75))
                return material.Label(theme.Material(), theme.TextSize, "From").Layout(gtx)
            })
        }),
        layout.Rigid(func(gtx layout.Context) layout.Dimensions {
            return widgets.DrawLine(gtx, theme.TableBorderColor, unit.Dp(35), unit.Dp(1))
        }),
        layout.Rigid(func(gtx layout.Context) layout.Dimensions {
            return layout.Inset{Left: unit.Dp(8), Right: unit.Dp(1)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                gtx.Constraints.Min.X = gtx.Dp(unit.Dp(40))
                return material.Label(theme.Material(), theme.TextSize, "Status").Layout(gtx)
            })
        }),
        layout.Rigid(func(gtx layout.Context) layout.Dimensions {
            return widgets.DrawLine(gtx, theme.TableBorderColor, unit.Dp(35), unit.Dp(1))
        }),
        layout.Rigid(func(gtx layout.Context) layout.Dimensions {
            return layout.Inset{Left: unit.Dp(8), Right: unit.Dp(1)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                return material.Label(theme.Material(), theme.TextSize, "Source Key/JSON Path").Layout(gtx)
            })
        }),
    )

}
