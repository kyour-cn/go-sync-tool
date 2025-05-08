package app

import (
	"app/internal/domain"
	"app/internal/global"
	"gioui.org/text"
	"gioui.org/widget/material"
	"image"

	"app/ui/apptheme"
	"app/ui/widgets"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
)

type Sidebar struct {
	Theme *apptheme.Theme

	flatButtons []*widgets.FlatButton
	Buttons     []*SideBarButton
	list        *widget.List

	cache *op.Ops

	clickables []*widget.Clickable

	selectedIndex int

	OnSelectedChanged func(index int)
}

type SideBarButton struct {
	Icon *widget.Icon
	Text string
}

func NewSidebar(theme *apptheme.Theme) *Sidebar {
	s := &Sidebar{
		Theme: theme,
		cache: new(op.Ops),

		Buttons: []*SideBarButton{
			{Icon: widgets.HomeIcon, Text: "首页"},
			//{Icon: widgets.MenuIcon, Text: "Envs"},
			//{Icon: widgets.WorkspacesIcon, Text: "Workspaces"},
			//{Icon: widgets.FileFolderIcon, Text: "Proto files"},
			// {Icon: widgets.TunnelIcon, Text: "Tunnels"},
			//{Icon: widgets.ConsoleIcon, Text: "Console"},
			// {Icon: widgets.LogsIcon, Text: "Logs"},
			// {Icon: widgets.SettingsIcon, Text: "Settings"},

			{Icon: widgets.SwapHoriz, Text: "同步"},
			{Icon: widgets.MenuIcon, Text: "配置"},
			{Icon: widgets.ConsoleIcon, Text: "日志"},
		},
		list: &widget.List{
			List: layout.List{
				Axis: layout.Vertical,
			},
		},
	}

	s.clickables = make([]*widget.Clickable, 0)
	for range s.Buttons {
		s.clickables = append(s.clickables, &widget.Clickable{})
	}

	s.makeButtons(theme)

	return s
}

func (s *Sidebar) makeButtons(theme *apptheme.Theme) {
	s.flatButtons = make([]*widgets.FlatButton, 0)
	for i, b := range s.Buttons {
		s.flatButtons = append(s.flatButtons, &widgets.FlatButton{
			Icon:              b.Icon,
			Text:              b.Text,
			IconPosition:      widgets.FlatButtonIconTop,
			Clickable:         s.clickables[i],
			SpaceBetween:      unit.Dp(5),
			BackgroundPadding: unit.Dp(1),
			CornerRadius:      5,
			MinWidth:          unit.Dp(70),
			BackgroundColor:   theme.SideBarBgColor,
			TextColor:         theme.SideBarTextColor,
			ContentPadding:    unit.Dp(5),
		})
	}
}

func (s *Sidebar) SelectedIndex() int {
	return s.selectedIndex
}

func (s *Sidebar) Layout(gtx layout.Context, theme *apptheme.Theme) layout.Dimensions {
	for i, c := range s.clickables {
		for c.Clicked(gtx) {
			s.selectedIndex = i
			if s.OnSelectedChanged != nil {
				s.OnSelectedChanged(i)
			}
		}
	}

	return layout.Background{}.Layout(gtx,
		func(gtx layout.Context) layout.Dimensions {
			if !theme.IsDark() {
				defer clip.UniformRRect(image.Rectangle{Max: gtx.Constraints.Min}, 0).Push(gtx.Ops).Pop()
				paint.Fill(gtx.Ops, theme.SideBarBgColor)
			}
			return layout.Dimensions{Size: gtx.Constraints.Min}
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Left:  unit.Dp(2),
				Right: unit.Dp(2),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{
							Axis:    layout.Vertical,
							Spacing: layout.SpaceBetween,
						}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return s.list.Layout(gtx, len(s.Buttons), func(gtx layout.Context, i int) layout.Dimensions {
									btn := s.flatButtons[i]
									if s.selectedIndex == i {
										btn.TextColor = theme.SideBarTextColor
									} else {
										btn.TextColor = widgets.Disabled(theme.SideBarTextColor)
									}

									return btn.Layout(gtx, theme)
								})
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{
									Axis:    layout.Vertical,
									Spacing: layout.SpaceBetween,
								}.Layout(gtx,
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {

										statusName := "初始化中"
										if global.State.Status == 1 {
											statusName = "待启动"
										} else if global.State.Status == 2 {
											statusName = "启动中"
										} else if global.State.Status == 3 {
											statusName = "运行中"
										} else if global.State.Status == 4 {
											statusName = "停止中"
										}

										return layout.Inset{
											Bottom: unit.Dp(5),
										}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
											gtx.Constraints.Min.X = gtx.Dp(70)
											st := material.Subtitle1(theme.Theme, statusName)
											st.Alignment = text.Middle
											return st.Layout(gtx)
										})
									}),
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										return layout.Inset{
											Bottom: unit.Dp(5),
										}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
											gtx.Constraints.Min.X = gtx.Dp(70)
											st := material.Subtitle1(theme.Theme, domain.Version)
											st.Alignment = text.Middle
											return st.Layout(gtx)
										})
									}),
								)
							}),
						)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if !theme.IsDark() {
							return layout.Dimensions{}
						}
						return widgets.DrawLine(gtx, theme.SeparatorColor, unit.Dp(gtx.Constraints.Max.Y), unit.Dp(1))
					}),
				)
			})
		},
	)
}
