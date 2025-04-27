package app

import (
    "app/internal/config"
    "app/ui/chapartheme"
    "app/ui/fonts"
    "app/ui/pages/console"
    "app/ui/pages/home"
    "app/ui/pages/sync"
    "app/ui/widgets"
    "context"
    "gioui.org/app"
    "gioui.org/layout"
    "gioui.org/op"
    "gioui.org/op/clip"
    "gioui.org/op/paint"
    "gioui.org/text"
    "gioui.org/unit"
    "gioui.org/widget/material"
    "github.com/go-gourd/gourd/event"
    "image"
)

type UI struct {
    Theme  *chapartheme.Theme
    window *app.Window

    sideBar *Sidebar
    header  *Header

    modal *widgets.MessageModal

    currentPage int

    consolePage *console.Console

    // 子页面
    homeView *home.View
    syncView *sync.View

    // 提示
    tip *widgets.Tip
}

// New creates a new UI using the Go Fonts.
func New(w *app.Window, appVersion string) (*UI, error) {
    u := &UI{
        window: w,
    }

    // 监听 window.invalidate 事件刷新UI
    event.Listen("window.invalidate", func(context.Context) {
        w.Invalidate()
    })

    fontCollection, err := fonts.Prepare()
    if err != nil {
        return nil, err
    }

    appConf, err := config.GetAppConfig()
    if err != nil {
        return nil, err
    }

    theme := material.NewTheme()
    theme.Shaper = text.NewShaper(text.WithCollection(fontCollection))
    u.Theme = chapartheme.New(theme, appConf.IsDark)
    // console need to be initialized before other pages as its listening for logs
    u.consolePage = console.New()

    // 头部
    u.header = NewHeader(w, u.Theme)

    // 侧边栏
    u.sideBar = NewSidebar(u.Theme)

    // 切换页面
    u.sideBar.OnSelectedChanged = func(index int) {
        u.currentPage = index
    }

    u.homeView = home.New()
    u.syncView = sync.New(u.Theme)

    u.header.OnThemeSwitched = u.onThemeChange

    u.tip = widgets.NewTip(nil)

    // 监听tips显示
    event.Listen("tips.show", func(ctx context.Context) {
        tips := ctx.Value("tipMsg")
        u.tip.Show(tips.(string), 2)
    })

    // 监听modal.message事件
    event.Listen("modal.message", func(ctx context.Context) {
        msg := ctx.Value("modalMsg")
        u.showModal(msg.(string))
    })

    return u, u.load()
}

func (u *UI) showModal(message string) {
    u.modal = widgets.NewMessageModal("提示", message, widgets.MessageModalTypeInfo, func(text string) {
        u.modal.Hide()
    }, widgets.ModalOption{Text: "确定"})
    u.modal.Show()
}

func (u *UI) onThemeChange(isDark bool) error {
    u.Theme.Switch(isDark)

    appConf, err := config.GetAppConfig()
    if err != nil {
        return err
    }

    appConf.IsDark = isDark
    if err := config.SetAppConfig(appConf); err != nil {
        return err
    }

    return nil
}

func (u *UI) load() error {

    appConf, err := config.GetAppConfig()
    if err != nil {
        return err
    }

    u.header.SetTheme(appConf.IsDark)
    u.Theme.Switch(appConf.IsDark)

    return nil
}

func (u *UI) Run() error {
    // ops are the operations from the UI
    var ops op.Ops

    for {
        switch e := u.window.Event().(type) {
        // this is sent when the application should re-render.
        case app.FrameEvent:
            gtx := app.NewContext(&ops, e)
            // render and handle UI.
            u.Layout(gtx)
            // render and handle the operations from the UI.
            e.Frame(gtx.Ops)
        // this is sent when the application is closed.
        case app.DestroyEvent:
            return e.Err
        }
    }
}

// Layout displays the main program layout.
func (u *UI) Layout(gtx layout.Context) layout.Dimensions {
    // set the background color
    macro := op.Record(gtx.Ops)
    rect := image.Rectangle{
        Max: image.Point{
            X: gtx.Constraints.Max.X,
            Y: gtx.Constraints.Max.Y,
        },
    }
    paint.FillShape(gtx.Ops, u.Theme.Palette.Bg, clip.Rect(rect).Op())
    background := macro.Stop()

    background.Add(gtx.Ops)

    u.modal.Layout(gtx, u.Theme)

    return layout.Stack{}.Layout(gtx,

        // 固定元素
        layout.Expanded(func(gtx layout.Context) layout.Dimensions {
            gtx.Constraints.Min.X = gtx.Constraints.Max.X
            return layout.Flex{Axis: layout.Vertical, Spacing: 0}.Layout(gtx,

                // 头部
                layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                    return u.header.Layout(gtx, u.Theme)
                }),

                // 主体
                layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
                    return layout.Flex{Axis: layout.Horizontal, Spacing: 0}.Layout(gtx,

                        // 侧边栏
                        layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                            return u.sideBar.Layout(gtx, u.Theme)
                        }),

                        // 内容页
                        layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
                            switch u.currentPage {
                            case 0:
                                return u.homeView.Layout(gtx, u.Theme)
                            case 1:
                                return u.syncView.Layout(gtx, u.Theme)
                            case 3:
                                return u.consolePage.Layout(gtx, u.Theme)
                            }
                            return layout.Dimensions{}
                        }),
                    )
                }),
            )
        }),
        // 悬浮元素
        layout.Stacked(func(gtx layout.Context) layout.Dimensions {

            if !u.tip.IsShow {
                return layout.Dimensions{}
            }

            // 获取转换比例 (物理像素/DP)
            pxPerDp := gtx.Metric.PxPerDp

            // 使用 Inset 定位悬浮元素
            inset := layout.Inset{
                Top:  unit.Dp((float32(gtx.Constraints.Max.Y) / pxPerDp) / 2),
                Left: unit.Dp((float32(gtx.Constraints.Max.X) / pxPerDp) / 2),
            }
            return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                return u.tip.Layout(gtx)
            })
        }),

    )
}
