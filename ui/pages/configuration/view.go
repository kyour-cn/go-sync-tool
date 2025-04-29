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

type ConfForm struct {
    projectName *widgets.LabeledInput

    erpDbType *widgets.DropDown
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
                Hint:           "单位简称如(巨蟹科技)",
            },
            erpDbType: widgets.NewDropDown(
                theme,
                widgets.NewDropDownOption("Mysql").WithValue("mysql"),
                widgets.NewDropDownOption("SqlServer").WithValue("sqlserver"),
                widgets.NewDropDownOption("Oracle").WithValue("oracle"),
            ),
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
            layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                return subFormInset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                    return v.confForm.projectName.Layout(gtx, theme)
                })
            }),
        )
    })

}
