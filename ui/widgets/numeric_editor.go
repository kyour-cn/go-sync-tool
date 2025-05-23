package widgets

import (
	"strconv"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"app/ui/apptheme"
)

type NumericEditor struct {
	widget.Editor
}

func (n *NumericEditor) Value() int {
	v, _ := strconv.Atoi(n.Text())
	return v
}

func (n *NumericEditor) Layout(gtx layout.Context, theme *apptheme.Theme) layout.Dimensions {
	for {
		event, ok := n.Update(gtx)
		if !ok {
			break
		}

		if _, ok := event.(widget.ChangeEvent); ok {
			if n.Text() != "" {
				if _, err := strconv.Atoi(n.Text()); err != nil {
					n.SetText(n.Text()[:len(n.Text())-1])
				}
			}
		}
	}

	editor := material.Editor(theme.Material(), &n.Editor, "0")
	editor.SelectionColor = theme.TextSelectionColor
	return editor.Layout(gtx)
}
