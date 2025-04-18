package sync

import (
	"app/internal/domain"
	"app/ui/chapartheme"
	"app/ui/widgets"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type View struct {
	logs []domain.Log

	list *widget.List

	clearButton *widget.Clickable

	startButton widget.Clickable
}

func New() *View {
	c := &View{
		list: &widget.List{
			List: layout.List{
				Axis: layout.Vertical,
			},
		},
		clearButton: &widget.Clickable{},
	}

	return c
}

func (c *View) Layout(gtx layout.Context, theme *chapartheme.Theme) layout.Dimensions {

	return layout.Inset{
		Top: unit.Dp(30), Left: unit.Dp(250), Right: unit.Dp(250),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,

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

			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{
					Top: unit.Dp(30), Left: unit.Dp(250), Right: unit.Dp(250),
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
