package app

import (
	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"app/ui/apptheme"
	"app/ui/widgets"
)

type Header struct {
	w *app.Window

	theme       *apptheme.Theme
	envDropDown *widgets.DropDown

	// modal is used to show error and messages to the user
	modal *widgets.MessageModal

	selectedWorkspace string
	//workspaceDropDown *widgets.DropDown

	themeSwitcherClickable widget.Clickable

	isDarkMode    bool
	iconDarkMode  material.LabelStyle
	iconLightMode material.LabelStyle

	OnThemeSwitched func(isDark bool) error
}

func NewHeader(w *app.Window, theme *apptheme.Theme) *Header {
	h := &Header{
		w:     w,
		theme: theme,
		//workspacesState: workspacesState,
	}
	h.iconDarkMode = widgets.MaterialIcons("dark_mode", theme)
	h.iconLightMode = widgets.MaterialIcons("light_mode", theme)

	h.envDropDown = widgets.NewDropDown(theme)
	h.envDropDown.MaxWidth = unit.Dp(150)
	return h
}

func (h *Header) showError(err error) {
	h.modal = widgets.NewMessageModal("Error", err.Error(), widgets.MessageModalTypeErr, func(_ string) {
		h.modal.Hide()
	}, widgets.ModalOption{Text: "Ok"})
	h.modal.Show()
}

func (h *Header) SetTheme(isDark bool) {
	h.isDarkMode = isDark
}

func (h *Header) themeSwitchIcon() material.LabelStyle {
	if h.isDarkMode {
		h.iconDarkMode = widgets.MaterialIcons("dark_mode", h.theme)
		return h.iconDarkMode
	}
	h.iconLightMode = widgets.MaterialIcons("light_mode", h.theme)
	return h.iconLightMode
}

func (h *Header) Layout(gtx layout.Context, theme *apptheme.Theme) layout.Dimensions {
	inset := layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(4), Left: unit.Dp(4)}

	if h.themeSwitcherClickable.Clicked(gtx) {
		h.isDarkMode = !h.isDarkMode
		if h.OnThemeSwitched != nil {
			if err := h.OnThemeSwitched(h.isDarkMode); err != nil {
				h.showError(err)
			}
			h.w.Invalidate()
		}
	}

	content := layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Spacing: layout.SpaceBetween}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Inset{Left: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return material.H6(h.theme.Material(), "巨蟹科技").Layout(gtx)
							})
						}),
						//layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						//    return layout.Inset{Left: unit.Dp(20), Right: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						//        return h.workspaceDropDown.Layout(gtx, theme)
						//    })
						//}),
					)
				}),
				//layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				//    gtx.Constraints.Max.X /= 3
				//    return h.headerSearch.Layout(gtx, theme)
				//}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return widgets.Clickable(gtx, &h.themeSwitcherClickable, unit.Dp(4), h.themeSwitchIcon().Layout)
						}),
						//layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						//    return layout.Inset{Left: unit.Dp(20), Right: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						//        return h.envDropDownenvDropDown.Layout(gtx, theme)
						//    })
						//}),
					)
				}),
			)
		})
	})

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		content,
		widgets.DrawLineFlex(theme.SeparatorColor, unit.Dp(1), unit.Dp(gtx.Constraints.Max.X)),
	)
}
