package sync

import (
	"app/internal/config"
	"app/internal/global"
	"app/internal/task"
	"app/ui/apptheme"
	"app/ui/widgets"
	"app/ui/widgets/codeeditor"
	"context"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	giox "gioui.org/x/component"
	"github.com/go-gourd/gourd/event"
	"image/color"
	"log/slog"
	"strconv"
	"time"
)

type View struct {
	startButton  widget.Clickable
	split        widgets.SplitView
	treeView     *widgets.TreeView
	codeEdit     *codeeditor.CodeEditor
	selectedNode *task.Task // 当前选中
	nodeStatus   *widget.Bool
	editing      bool
	intervalTime *widgets.LabeledInput
}

var leftTree *widgets.TreeView

func New(theme *apptheme.Theme) *View {

	// 左侧列表
	var leftTreeNode []*widgets.TreeNode

	leftTree = widgets.NewTreeView(leftTreeNode)

	updateNode := func() {

		var newTreeNode []*widgets.TreeNode

		for _, node := range task.List {
			tn := &widgets.TreeNode{
				Text:       node.Label,
				Identifier: node.Name,
				MenuOptions: []string{
					"校验SQL",
					"清除缓存",
				},
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

			if !node.Config.Status {
				// 灰色
				tn.PrefixColor = color.NRGBA{
					R: 0x8b, G: 0x8b, B: 0x8b, A: 0xff,
				}
			}

			newTreeNode = append(newTreeNode, tn)
		}

		leftTree.SetNodes(newTreeNode)

	}

	updateNode()
	event.Listen("task.config", func(ctx context.Context) {
		updateNode()
	})

	c := &View{
		split: widgets.SplitView{
			// Ratio:       -0.64,
			Resize: giox.Resize{
				Ratio: 0.19,
			},
			BarWidth: unit.Dp(2),
		},
		treeView:   leftTree,
		nodeStatus: new(widget.Bool),
		intervalTime: &widgets.LabeledInput{
			Label:          "任务间隔（秒）",
			SpaceBetween:   5,
			MinEditorWidth: unit.Dp(10),
			MinLabelWidth:  unit.Dp(10),
			Editor:         widgets.NewPatternEditor(),
		},
	}

	// 点击菜单子项
	leftTree.SetOnMenuItemClick(func(node *widgets.TreeNode, item string) {
		switch item {
		case "校验SQL":
			slog.Info("开始校验SQL")
			err := task.ValidateSql(node.Identifier)
			if err != nil {
				// 提示框
				params := context.WithValue(context.Background(), "modalMsg", err.Error())
				event.Trigger("modal.message", params)
				slog.Error("校验SQL失败：" + err.Error())
			} else {
				event.Trigger("tips.show", context.WithValue(context.Background(), "tipMsg", "校验SQL成功"))
			}
		case "清除缓存":
			slog.Info("清除缓存")
			err := task.ClearCache(node.Identifier)
			if err != nil {
				// 提示框
				params := context.WithValue(context.Background(), "modalMsg", err.Error())
				event.Trigger("modal.message", params)
				slog.Error("清除缓存失败：" + err.Error())
			} else {
				event.Trigger("tips.show", context.WithValue(context.Background(), "tipMsg", "清除缓存成功"))
			}

		default:
			slog.Warn("操作不支持：" + item)
		}
	})

	// 设置点击事件
	leftTree.OnNodeClick(func(node *widgets.TreeNode) {
		for _, t := range task.List {

			if t.Name == node.Identifier {
				// 获取配置
				_conf, err := config.GetTaskConfig(t.Name)
				if err != nil {
					_conf = &config.TaskConfig{}
				}
				SetCodeEdit(c, theme, _conf.Sql)

				// 设置状态
				c.nodeStatus.Value = _conf.Status

				// 设置间隔
				c.intervalTime.SetText(strconv.Itoa(_conf.IntervalTime))

				c.selectedNode = &t
				break
			}
		}
	})

	// 设置默认编辑器
	SetCodeEdit(c, theme, "请点击左侧选择编辑对象")
	c.codeEdit.SetReadOnly(true)

	return c
}

func SetCodeEdit(c *View, theme *apptheme.Theme, code string) {
	c.codeEdit = codeeditor.NewCodeEditor(code, "", theme)
	c.codeEdit.SetReadOnly(false)
	c.codeEdit.SetOnChanged(func(text string) {
		c.editing = true
	})

	c.editing = false

	// 刷新窗口
	event.Trigger("window.invalidate", context.Background())

	// 这里主要是解决代码编辑器的行号显示有延迟
	go func() {
		time.Sleep(time.Millisecond * 100)
		event.Trigger("window.invalidate", context.Background())
	}()
}

// 保存
func (v *View) onSave() {

	if global.State.Status != 1 {
		msg := "正在运行中，不允许更改配置，请先停止任务"
		slog.Info(msg)
		// 提示框
		params := context.WithValue(context.Background(), "modalMsg", msg)
		event.Trigger("modal.message", params)

		return
	}

	_conf, err := config.GetTaskConfig(v.selectedNode.Name)
	if err != nil {
		slog.Info("Get sql not found: " + err.Error())
		_conf = &config.TaskConfig{
			Name: v.selectedNode.Name,
		}
	}

	// 更新数据
	_conf.Sql = v.codeEdit.Code()
	_conf.Status = v.nodeStatus.Value
	_conf.IntervalTime, _ = strconv.Atoi(v.intervalTime.Text())

	// 保存到配置
	err = config.SetTaskConfig(v.selectedNode.Name, _conf)
	if err != nil {
		slog.Warn("Set sql error: " + err.Error())
	}

	// 提示保存成功
	params := context.WithValue(context.Background(), "tipMsg", "保存成功")
	event.Trigger("tips.show", params)

	slog.Info("保存 " + v.selectedNode.Name + " 配置成功成功")

	v.editing = false
}

func (v *View) Layout(gtx layout.Context, theme *apptheme.Theme) layout.Dimensions {

	// 保存按钮
	if v.startButton.Clicked(gtx) && v.selectedNode != nil {
		v.onSave()
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
				Top: unit.Dp(10), Left: unit.Dp(10), Right: unit.Dp(10), Bottom: unit.Dp(10),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,

					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
							// 标题
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								text := "请选择"
								if v.selectedNode != nil {
									text = v.selectedNode.Label
								}
								if v.editing {
									text += " （未保存）"
								}
								return material.H5(theme.Material(), text).Layout(gtx)
							}),
							// 填充中间空间
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								return layout.Spacer{}.Layout(gtx)
								//return layout.Dimensions{} // 不渲染任何内容，仅占用空间
							}),

							// 右侧
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										// 最大宽度
										gtx.Constraints.Max.X = gtx.Dp(unit.Dp(160))
										return v.intervalTime.Layout(gtx, theme)
									}),
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										return material.Body1(theme.Material(), " 是否启用：").Layout(gtx)
									}),

									// 是否启用
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										s := material.Switch(theme.Material(), v.nodeStatus, "开关")
										s.Color.Enabled = theme.SwitchBgColor
										s.Color.Disabled = theme.Palette.Fg
										return layout.Inset{
											Left:  unit.Dp(10),
											Right: unit.Dp(10),
											Top:   unit.Dp(5),
										}.Layout(gtx,
											s.Layout,
										)
									}),
									// 保存按钮
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										if v.selectedNode == nil {
											gtx = gtx.Disabled()
										}
										newBtn := widgets.Button(theme.Material(), &v.startButton, nil, widgets.IconPositionStart, "保存")
										newBtn.Color = theme.ButtonTextColor
										newBtn.Background = theme.SendButtonBgColor
										return newBtn.Layout(gtx, theme)
									}),
								)
							}),
						)
					}),

					// 描述
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						text := " "
						if v.selectedNode != nil {
							text = v.selectedNode.Description
						}
						return material.Body1(theme.Material(), text).Layout(gtx)
					}),

					// 任务自定义UI
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if v.selectedNode == nil || v.selectedNode.Handle == nil {
							return layout.Dimensions{}
						}
						return v.selectedNode.Handle.ConfigLayout(gtx, theme, v.selectedNode)
					}),

					// 代码编辑器
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						// 写入类型暂不支持自定义SQL
						if v.selectedNode == nil || v.selectedNode.Type == 1 {
							return layout.Dimensions{}
						}
						return layout.Inset{
							Top: unit.Dp(10),
						}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return v.codeEdit.Layout(gtx, theme, "")
						})
					}),
				)
			})
		},
	)

}
