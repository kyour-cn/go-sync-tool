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

    confForm *ConfForm
}

var (
    // 虚拟列表，用于创建滚动布局
    virtualList = &widget.List{List: layout.List{Axis: layout.Vertical}}
)

type ConfForm struct {
    projectName *widgets.LabeledInput

    // 商城数据库
    shopDbType *widgets.DropDown
    shopDbHost *widgets.LabeledInput
    shopDbPort *widgets.LabeledInput
    shopDbName *widgets.LabeledInput
    shopDbUser *widgets.LabeledInput
    shopDbPass *widgets.LabeledInput

    // ERP数据库
    erpDbType *widgets.DropDown
    erpDbHost *widgets.LabeledInput
    erpDbPort *widgets.LabeledInput
    erpDbName *widgets.LabeledInput
    erpDbUser *widgets.LabeledInput
    erpDbPass *widgets.LabeledInput
}

func New(theme *chapartheme.Theme) *View {

    c := &View{
        confForm: &ConfForm{
            projectName: &widgets.LabeledInput{
                Label:          "项目名称",
                SpaceBetween:   5,
                MinEditorWidth: unit.Dp(150),
                MinLabelWidth:  unit.Dp(80),
                Editor:         widgets.NewPatternEditor(),
                Hint:           "公司简称",
            },
            shopDbType: widgets.NewDropDown(
                theme,
                widgets.NewDropDownOption("Mysql").WithValue("mysql"),
            ),
            shopDbHost: &widgets.LabeledInput{
                Label:          "地址",
                SpaceBetween:   5,
                MinEditorWidth: unit.Dp(150),
                MinLabelWidth:  unit.Dp(80),
                Editor:         widgets.NewPatternEditor(),
                Hint:           "IP或域名",
            },
            shopDbPort: &widgets.LabeledInput{
                Label:          "端口",
                SpaceBetween:   5,
                MinEditorWidth: unit.Dp(150),
                MinLabelWidth:  unit.Dp(80),
                Editor:         widgets.NewPatternEditor(),
            },
            shopDbName: &widgets.LabeledInput{
                Label:          "数据库名",
                SpaceBetween:   5,
                MinEditorWidth: unit.Dp(150),
                MinLabelWidth:  unit.Dp(80),
                Editor:         widgets.NewPatternEditor(),
            },
            shopDbUser: &widgets.LabeledInput{
                Label:          "用户名",
                SpaceBetween:   5,
                MinEditorWidth: unit.Dp(150),
                MinLabelWidth:  unit.Dp(80),
                Editor:         widgets.NewPatternEditor(),
            },
            shopDbPass: &widgets.LabeledInput{
                Label:          "密码",
                SpaceBetween:   5,
                MinEditorWidth: unit.Dp(150),
                MinLabelWidth:  unit.Dp(80),
                Editor:         widgets.NewPatternEditor(),
            },
            erpDbType: widgets.NewDropDown(
                theme,
                widgets.NewDropDownOption("Mysql").WithValue("mysql"),
                widgets.NewDropDownOption("SqlServer").WithValue("sqlserver"),
                widgets.NewDropDownOption("Oracle").WithValue("oracle"),
            ),
            erpDbHost: &widgets.LabeledInput{
                Label:          "地址",
                SpaceBetween:   5,
                MinEditorWidth: unit.Dp(150),
                MinLabelWidth:  unit.Dp(80),
                Editor:         widgets.NewPatternEditor(),
                Hint:           "IP或域名",
            },
            erpDbPort: &widgets.LabeledInput{
                Label:          "端口",
                SpaceBetween:   5,
                MinEditorWidth: unit.Dp(150),
                MinLabelWidth:  unit.Dp(80),
                Editor:         widgets.NewPatternEditor(),
            },
            erpDbName: &widgets.LabeledInput{
                Label:          "数据库名",
                SpaceBetween:   5,
                MinEditorWidth: unit.Dp(150),
                MinLabelWidth:  unit.Dp(80),
                Editor:         widgets.NewPatternEditor(),
            },
            erpDbUser: &widgets.LabeledInput{
                Label:          "用户名",
                SpaceBetween:   5,
                MinEditorWidth: unit.Dp(150),
                MinLabelWidth:  unit.Dp(80),
                Editor:         widgets.NewPatternEditor(),
            },
            erpDbPass: &widgets.LabeledInput{
                Label:          "密码",
                SpaceBetween:   5,
                MinEditorWidth: unit.Dp(150),
                MinLabelWidth:  unit.Dp(80),
                Editor:         widgets.NewPatternEditor(),
            },
        },
    }

    return c
}

// TODO 待实现
func (v *View) Layout(gtx layout.Context, theme *chapartheme.Theme) layout.Dimensions {
    topButtonInset := layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(4)}

    subFormInset := layout.Inset{
        Top:    unit.Dp(8),
        Bottom: unit.Dp(4),
        Left:   unit.Dp(10),
    }

    return material.List(theme.Material(), virtualList).Layout(gtx, 1, func(gtx layout.Context, _ int) layout.Dimensions {

        return layout.Inset{
            Left:  unit.Dp(10),
            Right: unit.Dp(10),
        }.Layout(gtx, func(gtx layout.Context) layout.Dimensions {

            return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
                layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                    return topButtonInset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                        return v.confForm.projectName.Layout(gtx, theme)
                    })
                }),

                layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                    return topButtonInset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                        return material.Label(theme.Material(), theme.TextSize, "商城配置:").Layout(gtx)
                    })
                }),

                // 商城数据库
                layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                    return subFormInset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                        return layout.Flex{
                            Axis:      layout.Horizontal,
                            Alignment: layout.Middle,
                        }.Layout(gtx,
                            layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                                gtx.Constraints.Min.X = gtx.Dp(85)
                                return material.Label(theme.Material(), theme.TextSize, "数据库类型").Layout(gtx)
                            }),
                            layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                                v.confForm.shopDbType.MaxWidth = unit.Dp(162)
                                return v.confForm.shopDbType.Layout(gtx, theme)

                            }),
                        )
                    })
                }),
                layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                    return subFormInset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                        return v.confForm.shopDbHost.Layout(gtx, theme)
                    })
                }),
                layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                    return subFormInset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                        return v.confForm.shopDbPort.Layout(gtx, theme)
                    })
                }),
                layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                    return subFormInset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                        return v.confForm.shopDbName.Layout(gtx, theme)
                    })
                }),
                layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                    return subFormInset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                        return v.confForm.shopDbUser.Layout(gtx, theme)
                    })
                }),
                layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                    return subFormInset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                        return v.confForm.shopDbPass.Layout(gtx, theme)
                    })
                }),

                // ERP 数据库
                layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                    return topButtonInset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                        return material.Label(theme.Material(), theme.TextSize, "ERP配置:").Layout(gtx)
                    })
                }),
                layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                    return subFormInset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                        return layout.Flex{
                            Axis:      layout.Horizontal,
                            Alignment: layout.Middle,
                        }.Layout(gtx,
                            layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                                gtx.Constraints.Min.X = gtx.Dp(85)
                                return material.Label(theme.Material(), theme.TextSize, "数据库类型").Layout(gtx)
                            }),
                            layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                                v.confForm.erpDbType.MaxWidth = unit.Dp(162)
                                return v.confForm.erpDbType.Layout(gtx, theme)

                            }),
                        )
                    })
                }),
                layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                    return subFormInset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                        return v.confForm.erpDbHost.Layout(gtx, theme)
                    })
                }),
                layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                    return subFormInset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                        return v.confForm.erpDbPort.Layout(gtx, theme)
                    })
                }),
                layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                    return subFormInset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                        return v.confForm.erpDbName.Layout(gtx, theme)
                    })
                }),
                layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                    return subFormInset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                        return v.confForm.erpDbUser.Layout(gtx, theme)
                    })
                }),
                layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                    return subFormInset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                        return v.confForm.erpDbPass.Layout(gtx, theme)
                    })
                }),
            )

        })
    })

}
