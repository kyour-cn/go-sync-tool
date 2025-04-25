package widgets

import (
    "context"
    "gioui.org/layout"
    "gioui.org/widget/material"
    "gioui.org/x/component"
    "github.com/go-gourd/gourd/event"
    "time"
)

type Tip struct {
    th        *material.Theme
    Tooltip   component.Tooltip
    IsShow    bool
    CloseTime time.Time
}

func (t *Tip) Layout(gtx layout.Context) layout.Dimensions {
    if !t.IsShow {
        return layout.Dimensions{}
    }
    return t.Tooltip.Layout(gtx)
}

func (t *Tip) Show(content string, expire int) {
    t.Tooltip.Text.Text = content
    t.IsShow = true
    t.CloseTime = time.Now().Add(time.Second * time.Duration(expire))

    ctx := context.Background()

    event.Trigger("window.invalidate", ctx)

    go func() {
        time.Sleep(time.Second * 2)
        if time.Now().After(t.CloseTime) {
            t.IsShow = false
            event.Trigger("window.invalidate", ctx)
        }
    }()
}

func NewTip(theme *material.Theme) *Tip {

    if theme == nil {
        theme = material.NewTheme()
    }

    tip := Tip{
        th:      theme,
        Tooltip: component.DesktopTooltip(theme, ""),
    }
    tip.th = theme

    return &tip
}
