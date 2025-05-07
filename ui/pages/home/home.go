package home

import (
    "app/assets"
    "app/internal/domain"
    "app/internal/global"
    "app/internal/task"
    "app/ui/chapartheme"
    "app/ui/widgets"
    "context"
    "gioui.org/layout"
    "gioui.org/op/paint"
    "gioui.org/unit"
    "gioui.org/widget"
    "gioui.org/widget/material"
    "github.com/go-gourd/gourd/event"
    "strconv"
    "strings"
)

type View struct {
    imageOp     paint.ImageOp
    startButton widget.Clickable

    //testModal *widgets.TestModal
}

func New() *View {

    // Load logo image
    data, err := assets.LoadImage("logo.png")
    if err != nil {
        panic(err)
    }

    c := &View{
        imageOp: paint.NewImageOp(data),
        //testModal: widgets.NewTestModal(),
    }

    return c
}

func (c *View) Layout(gtx layout.Context, theme *chapartheme.Theme) layout.Dimensions {

    //c.testModal.Layout(gtx, theme)

    if c.startButton.Clicked(gtx) {
        if global.State.Status == 1 {
            event.Trigger("task.start", context.Background())
            //c.testModal.SetVisible(true)
        } else if global.State.Status == 3 {
            event.Trigger("task.stop", context.Background())
            //c.testModal.SetVisible(false)
        }
    }

    return layout.Inset{
        Top: unit.Dp(30), Left: unit.Dp(100), Right: unit.Dp(100),
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
                    return material.H3(theme.Material(), domain.AppName).Layout(gtx)
                })
            }),
            layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                    return material.H5(theme.Material(), domain.Version).Layout(gtx)
                })
            }),

            // 显示运行状态
            layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                return layout.Inset{
                    Top: unit.Dp(30),
                }.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                    return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
                        layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                            taskMsg := "无"
                            var taskList []string
                            for _, v := range task.List {
                                if v.Config.Status {
                                    taskList = append(taskList, v.Label)
                                }
                            }
                            if len(taskList) > 0 {
                                taskMsg = strings.Join(taskList, "、")
                            }
                            return material.Body1(theme.Material(), "启用中任务: "+taskMsg).Layout(gtx)
                        }),
                        layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                            taskMsg := "无"
                            var taskList []string
                            for _, v := range task.List {
                                if v.Status {
                                    taskList = append(taskList, v.Label+"("+strconv.Itoa(v.DoneCount)+"/"+strconv.Itoa(v.DataCount)+")")
                                }
                            }
                            if len(taskList) > 0 {
                                taskMsg = strings.Join(taskList, "、")
                            }
                            return material.Body1(theme.Material(), "运行中任务: "+taskMsg).Layout(gtx)
                        }),
                    )
                })
            }),

            layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                return layout.Inset{
                    Top: unit.Dp(30),
                }.Layout(gtx, func(gtx layout.Context) layout.Dimensions {

                    btnText := "启动同步"
                    if global.State.Status != 1 {
                        btnText = "停止同步"
                    }
                    if global.State.Status == 2 || global.State.Status == 4 {
                        gtx = gtx.Disabled()
                        btnText = "请稍后..."
                    }

                    newBtn := widgets.Button(theme.Material(), &c.startButton, nil, widgets.IconPositionStart, btnText)
                    newBtn.Color = theme.ButtonTextColor
                    newBtn.Background = theme.SendButtonBgColor

                    return newBtn.Layout(gtx, theme)
                })
            }),
        )
    })
}
