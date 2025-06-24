package configuration

import (
	"app/internal/config"
	"app/internal/global"
	"app/ui/apptheme"
	"app/ui/widgets"
	"context"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/go-gourd/gourd/event"
	"log/slog"
	"strconv"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

type View struct {
	targetEnvEditor *widget.Editor

	confForm *ConfForm

	// 保存按钮
	saveButton widget.Clickable
}

var (
	// 虚拟列表，用于创建滚动布局
	virtualList = &widget.List{List: layout.List{Axis: layout.Vertical}}
)

type ConfForm struct {
	// 项目名称
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

	// ERP编码
	erpDbCode *widgets.DropDown
}

func New(theme *apptheme.Theme) *View {

	c := &View{
		confForm: &ConfForm{
			projectName: &widgets.LabeledInput{
				Label:          "项目名称",
				SpaceBetween:   5,
				MinEditorWidth: unit.Dp(150),
				MinLabelWidth:  unit.Dp(80),
				Editor:         widgets.NewPatternEditor(),
				Hint:           "请输入平台名称或公司简称",
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
				//widgets.NewDropDownOption("SqlServer(ODBC)").WithValue("sqlserver-odbc"),
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

			erpDbCode: widgets.NewDropDown(
				theme,
				widgets.NewDropDownOption("GBK").WithValue("0"),
				widgets.NewDropDownOption("UTF8").WithValue("1"),
				widgets.NewDropDownOption("智能识别").WithValue("2"),
			),
		},
	}

	// 初始化获取配置
	shopDbConf, err := config.GetDBConfig("shop")
	if err != nil {
		shopDbConf = &config.DbConfig{}
	}

	c.confForm.erpDbType.SetSelectedByValue(shopDbConf.Type)
	c.confForm.shopDbHost.SetText(shopDbConf.Host)
	c.confForm.shopDbPort.SetText(strconv.Itoa(shopDbConf.Port))
	c.confForm.shopDbName.SetText(shopDbConf.Database)
	c.confForm.shopDbUser.SetText(shopDbConf.User)
	c.confForm.shopDbPass.SetText(shopDbConf.Pass)

	erpDbConf, err := config.GetDBConfig("erp")
	if err != nil {
		erpDbConf = &config.DbConfig{}
	}
	c.confForm.erpDbType.SetSelectedByValue(erpDbConf.Type)
	c.confForm.erpDbHost.SetText(erpDbConf.Host)
	c.confForm.erpDbPort.SetText(strconv.Itoa(erpDbConf.Port))
	c.confForm.erpDbName.SetText(erpDbConf.Database)
	c.confForm.erpDbUser.SetText(erpDbConf.User)
	c.confForm.erpDbPass.SetText(erpDbConf.Pass)

	appconf, err := config.GetAppConfig()
	if err != nil {
		appconf = &config.AppConfig{}
		slog.Error("Get app config error: " + err.Error())
	}

	c.confForm.projectName.SetText(appconf.ProjectName)
	c.confForm.erpDbCode.SetSelectedByValue(strconv.Itoa(appconf.ErpEncoding))

	return c
}

func (v *View) Save() {

	// 保存配置
	shopDbPort, _ := strconv.Atoi(v.confForm.shopDbPort.Text())
	shopDbConf := &config.DbConfig{
		Type:     v.confForm.shopDbType.GetSelected().Value,
		Host:     v.confForm.shopDbHost.Text(),
		Port:     shopDbPort,
		Database: v.confForm.shopDbName.Text(),
		User:     v.confForm.shopDbUser.Text(),
		Pass:     v.confForm.shopDbPass.Text(),
	}
	err := config.SetDBConfig("shop", shopDbConf)
	if err != nil {
		slog.Error("Set shop db config error: " + err.Error())
	}

	erpDbPort, _ := strconv.Atoi(v.confForm.erpDbPort.Text())
	erpDbConf := &config.DbConfig{
		Type:     v.confForm.erpDbType.GetSelected().Value,
		Host:     v.confForm.erpDbHost.Text(),
		Port:     erpDbPort,
		Database: v.confForm.erpDbName.Text(),
		User:     v.confForm.erpDbUser.Text(),
		Pass:     v.confForm.erpDbPass.Text(),
	}
	err = config.SetDBConfig("erp", erpDbConf)
	if err != nil {
		slog.Error("Set erp db config error: " + err.Error())
	}

	appconf, err := config.GetAppConfig()
	if err != nil {
		appconf = &config.AppConfig{}
		slog.Error("Get app config error: " + err.Error())
	}

	appconf.ProjectName = v.confForm.projectName.Text()
	appconf.ErpEncoding, _ = strconv.Atoi(v.confForm.erpDbCode.GetSelected().Value)
	global.State.ErpEncoding = appconf.ErpEncoding
	err = config.SetAppConfig(appconf)
	if err != nil {
		slog.Error("Set app config error: " + err.Error())
	}
}

func (v *View) Layout(gtx layout.Context, theme *apptheme.Theme) layout.Dimensions {
	topButtonInset := layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(4)}

	subFormInset := layout.Inset{
		Top:    unit.Dp(8),
		Bottom: unit.Dp(4),
		Left:   unit.Dp(10),
	}

	if v.saveButton.Clicked(gtx) {

		// 提示保存成功
		params := context.WithValue(context.Background(), "modalMsg", "保存成功")
		event.Trigger("modal.message", params)

		v.Save()

	}

	return material.List(theme.Material(), virtualList).Layout(gtx, 1, func(gtx layout.Context, _ int) layout.Dimensions {

		return layout.Inset{
			Left:  unit.Dp(10),
			Right: unit.Dp(10),
		}.Layout(gtx, func(gtx C) D {

			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return topButtonInset.Layout(gtx, func(gtx C) D {
						return v.confForm.projectName.Layout(gtx, theme)
					})
				}),

				layout.Rigid(func(gtx C) D {
					return topButtonInset.Layout(gtx, func(gtx C) D {
						return material.Label(theme.Material(), theme.TextSize, "商城配置:").Layout(gtx)
					})
				}),

				layout.Rigid(func(gtx C) D {
					return layout.Flex{}.Layout(gtx,
						layout.Flexed(.5, func(gtx C) D {
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									return subFormInset.Layout(gtx, func(gtx C) D {
										return layout.Flex{
											Axis:      layout.Horizontal,
											Alignment: layout.Middle,
										}.Layout(gtx,
											layout.Rigid(func(gtx C) D {
												gtx.Constraints.Min.X = gtx.Dp(85)
												return material.Label(theme.Material(), theme.TextSize, "数据库类型").Layout(gtx)
											}),
											layout.Rigid(func(gtx C) D {
												return v.confForm.shopDbType.Layout(gtx, theme)
											}),
										)
									})
								}),
								layout.Rigid(func(gtx C) D {
									return subFormInset.Layout(gtx, func(gtx C) D {
										return v.confForm.shopDbPort.Layout(gtx, theme)
									})
								}),
								layout.Rigid(func(gtx C) D {
									return subFormInset.Layout(gtx, func(gtx C) D {
										return v.confForm.shopDbUser.Layout(gtx, theme)
									})
								}),
							)
						}),
						layout.Flexed(.5, func(gtx C) D {
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									return subFormInset.Layout(gtx, func(gtx C) D {
										return v.confForm.shopDbHost.Layout(gtx, theme)
									})
								}),
								layout.Rigid(func(gtx C) D {
									return subFormInset.Layout(gtx, func(gtx C) D {
										return v.confForm.shopDbName.Layout(gtx, theme)
									})
								}),
								layout.Rigid(func(gtx C) D {
									return subFormInset.Layout(gtx, func(gtx C) D {
										return v.confForm.shopDbPass.Layout(gtx, theme)
									})
								}),
							)
						}),
					)
				}),

				//layout.Rigid(func(gtx C) D {
				//    return subFormInset.Layout(gtx, func(gtx C) D {
				//        return layout.Flex{
				//            Axis:      layout.Horizontal,
				//            Alignment: layout.Middle,
				//        }.Layout(gtx,
				//            layout.Rigid(func(gtx C) D {
				//                gtx.Constraints.Min.X = gtx.Dp(85)
				//                return material.Label(theme.Material(), theme.TextSize, "数据库类型").Layout(gtx)
				//            }),
				//            layout.Rigid(func(gtx C) D {
				//                v.confForm.shopDbType.MaxWidth = unit.Dp(162)
				//                return v.confForm.shopDbType.Layout(gtx, theme)
				//
				//            }),
				//        )
				//    })
				//}),
				//layout.Rigid(func(gtx C) D {
				//    return subFormInset.Layout(gtx, func(gtx C) D {
				//        return v.confForm.shopDbHost.Layout(gtx, theme)
				//    })
				//}),
				//layout.Rigid(func(gtx C) D {
				//    return subFormInset.Layout(gtx, func(gtx C) D {
				//        return v.confForm.shopDbPort.Layout(gtx, theme)
				//    })
				//}),
				//layout.Rigid(func(gtx C) D {
				//    return subFormInset.Layout(gtx, func(gtx C) D {
				//        return v.confForm.shopDbName.Layout(gtx, theme)
				//    })
				//}),
				//layout.Rigid(func(gtx C) D {
				//    return subFormInset.Layout(gtx, func(gtx C) D {
				//        return v.confForm.shopDbUser.Layout(gtx, theme)
				//    })
				//}),
				//layout.Rigid(func(gtx C) D {
				//    return subFormInset.Layout(gtx, func(gtx C) D {
				//        return v.confForm.shopDbPass.Layout(gtx, theme)
				//    })
				//}),

				// ERP 数据库
				layout.Rigid(func(gtx C) D {
					return topButtonInset.Layout(gtx, func(gtx C) D {
						return material.Label(theme.Material(), theme.TextSize, "ERP配置:").Layout(gtx)
					})
				}),

				layout.Rigid(func(gtx C) D {
					return layout.Flex{}.Layout(gtx,
						layout.Flexed(.5, func(gtx C) D {
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									return subFormInset.Layout(gtx, func(gtx C) D {
										return layout.Flex{
											Axis:      layout.Horizontal,
											Alignment: layout.Middle,
										}.Layout(gtx,
											layout.Rigid(func(gtx C) D {
												gtx.Constraints.Min.X = gtx.Dp(85)
												return material.Label(theme.Material(), theme.TextSize, "数据库类型").Layout(gtx)
											}),
											layout.Rigid(func(gtx C) D {
												return v.confForm.erpDbType.Layout(gtx, theme)
											}),
										)
									})
								}),
								layout.Rigid(func(gtx C) D {
									return subFormInset.Layout(gtx, func(gtx C) D {
										return v.confForm.erpDbPort.Layout(gtx, theme)
									})
								}),
								layout.Rigid(func(gtx C) D {
									return subFormInset.Layout(gtx, func(gtx C) D {
										return v.confForm.erpDbUser.Layout(gtx, theme)
									})
								}),
								layout.Rigid(func(gtx C) D {
									return subFormInset.Layout(gtx, func(gtx C) D {
										return layout.Flex{
											Axis:      layout.Horizontal,
											Alignment: layout.Middle,
										}.Layout(gtx,
											layout.Rigid(func(gtx C) D {
												gtx.Constraints.Min.X = gtx.Dp(85)
												return material.Label(theme.Material(), theme.TextSize, "字符编码").Layout(gtx)
											}),
											layout.Rigid(func(gtx C) D {
												return v.confForm.erpDbCode.Layout(gtx, theme)
											}),
										)
									})
								}),
							)
						}),
						layout.Flexed(.5, func(gtx C) D {
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									return subFormInset.Layout(gtx, func(gtx C) D {
										return v.confForm.erpDbHost.Layout(gtx, theme)
									})
								}),
								layout.Rigid(func(gtx C) D {
									return subFormInset.Layout(gtx, func(gtx C) D {
										return v.confForm.erpDbName.Layout(gtx, theme)
									})
								}),
								layout.Rigid(func(gtx C) D {
									return subFormInset.Layout(gtx, func(gtx C) D {
										return v.confForm.erpDbPass.Layout(gtx, theme)
									})
								}),
							)
						}),
					)
				}),
				//layout.Rigid(func(gtx C) D {
				//    return subFormInset.Layout(gtx, func(gtx C) D {
				//        return layout.Flex{
				//            Axis:      layout.Horizontal,
				//            Alignment: layout.Middle,
				//        }.Layout(gtx,
				//            layout.Rigid(func(gtx C) D {
				//                gtx.Constraints.Min.X = gtx.Dp(85)
				//                return material.Label(theme.Material(), theme.TextSize, "数据库类型").Layout(gtx)
				//            }),
				//            layout.Rigid(func(gtx C) D {
				//                v.confForm.erpDbType.MaxWidth = unit.Dp(162)
				//                return v.confForm.erpDbType.Layout(gtx, theme)
				//
				//            }),
				//        )
				//    })
				//}),
				//layout.Rigid(func(gtx C) D {
				//    return subFormInset.Layout(gtx, func(gtx C) D {
				//        return v.confForm.erpDbHost.Layout(gtx, theme)
				//    })
				//}),
				//layout.Rigid(func(gtx C) D {
				//    return subFormInset.Layout(gtx, func(gtx C) D {
				//        return v.confForm.erpDbPort.Layout(gtx, theme)
				//    })
				//}),
				//layout.Rigid(func(gtx C) D {
				//    return subFormInset.Layout(gtx, func(gtx C) D {
				//        return v.confForm.erpDbName.Layout(gtx, theme)
				//    })
				//}),
				//layout.Rigid(func(gtx C) D {
				//    return subFormInset.Layout(gtx, func(gtx C) D {
				//        return v.confForm.erpDbUser.Layout(gtx, theme)
				//    })
				//}),
				//layout.Rigid(func(gtx C) D {
				//    return subFormInset.Layout(gtx, func(gtx C) D {
				//        return v.confForm.erpDbPass.Layout(gtx, theme)
				//    })
				//}),

				// 保存按钮
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(10)}.Layout(gtx, func(gtx C) D {
						return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								gtx.Constraints.Max.X = gtx.Dp(162)
								newBtn := widgets.Button(theme.Material(), &v.saveButton, nil, widgets.IconPositionStart, "保存")
								newBtn.Color = theme.ButtonTextColor
								newBtn.Background = theme.SendButtonBgColor
								return newBtn.Layout(gtx, theme)
							}),
						)
					})
				}),
			)

		})
	})

}
