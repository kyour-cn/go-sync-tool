package sync

import (
    "app/internal/task"
    "app/ui/chapartheme"
    "app/ui/widgets"
    "app/ui/widgets/codeeditor"
    "fmt"
    "gioui.org/layout"
    "gioui.org/unit"
    "gioui.org/widget"
    "gioui.org/widget/material"
    giox "gioui.org/x/component"
    "image/color"
)

type View struct {
    startButton widget.Clickable
    split       widgets.SplitView
    treeView    *widgets.TreeView
    codeEdit    *codeeditor.CodeEditor
    // 当前选中的名称
    selectedNode string
}

func New(theme *chapartheme.Theme) *View {
    // 初始化编辑器
    codeEditor := codeeditor.NewCodeEditor("test", codeeditor.CodeLanguageShell, theme)

    // 左侧列表
    var leftTreeNode []*widgets.TreeNode

    for _, node := range task.List {
        tn := &widgets.TreeNode{
            Text:       node.Label,
            Identifier: node.Name,
        }
        if node.Type == 1 {
            tn.Prefix = "W"
            // 红色
            tn.PrefixColor = color.NRGBA{
                R: 0xff, G: 0x73, B: 0x73, A: 0xff,
            }
        } else {
            tn.Prefix = "Q"
            // 绿色
            tn.PrefixColor = color.NRGBA{
                R: 0x8b, G: 0xc3, B: 0x4a, A: 0xff,
            }
        }

        leftTreeNode = append(leftTreeNode, tn)
    }

    leftTree := widgets.NewTreeView(leftTreeNode)

    c := &View{
        split: widgets.SplitView{
            // Ratio:       -0.64,
            Resize: giox.Resize{
                Ratio: 0.19,
            },
            BarWidth: unit.Dp(2),
        },
        treeView:     leftTree,
        codeEdit:     codeEditor,
        selectedNode: "请选择",
    }

    // 设置点击事件
    leftTree.OnNodeClick(func(node *widgets.TreeNode) {
        for _, t := range task.List {

            if t.Name == node.Identifier {
                SetCodeEdit(c, theme, node.Text)
                c.selectedNode = t.Label
                break
            }
        }

        //c.codeEdit = codeeditor.NewCodeEditor(node.Text, codeeditor.CodeLanguageShell, theme)
        //codeEditor.SetCode(node.Text)
    })

    // 设置默认编辑器
    SetCodeEdit(c, theme, "请点击左侧选择编辑对象")
    c.codeEdit.SetReadOnly(true)

    return c
}

func SetCodeEdit(c *View, theme *chapartheme.Theme, code string) {
    c.codeEdit = codeeditor.NewCodeEditor(code, codeeditor.CodeLanguageShell, theme)
    c.codeEdit.SetReadOnly(false)
}

func (v *View) Layout(gtx layout.Context, theme *chapartheme.Theme) layout.Dimensions {

    if v.startButton.Clicked(gtx) {
        fmt.Println("Save button clicked")
    }

    return v.split.Layout(gtx, theme,
        func(gtx layout.Context) layout.Dimensions {
            return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
                layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
                    return layout.Inset{Top: unit.Dp(10), Right: 0}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                        return v.treeView.Layout(gtx, theme)
                    })
                }),
            )
        },
        func(gtx layout.Context) layout.Dimensions {
            return layout.Inset{
                Top: unit.Dp(30), Left: unit.Dp(10), Right: unit.Dp(10),
            }.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,

                    layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                        return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
                            layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                                return material.H5(theme.Material(), v.selectedNode).Layout(gtx)
                            }),
                            // 填充中间空间
                            layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
                                return layout.Spacer{}.Layout(gtx)
                                //return layout.Dimensions{} // 不渲染任何内容，仅占用空间
                            }),
                            // 保存按钮
                            layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                                newBtn := widgets.Button(theme.Material(), &v.startButton, nil, widgets.IconPositionStart, "保存")
                                newBtn.Color = theme.ButtonTextColor
                                newBtn.Background = theme.SendButtonBgColor
                                return newBtn.Layout(gtx, theme)
                            }),
                        )
                    }),

                    // 代码编辑器
                    layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                        return layout.Inset{
                            Top: unit.Dp(10),
                        }.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                            return v.codeEdit.Layout(gtx, theme, "Shell")
                        })
                    }),

                )
            })
        },
    )

}
