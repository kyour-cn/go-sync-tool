package codeeditor

import (
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"

	"app/ui/apptheme"
	"app/ui/widgets"

	"github.com/oligo/gvcode"
	wgvcode "github.com/oligo/gvcode/widget"
)

const (
	CodeLanguageJSON       = "JSON"
	CodeLanguageYAML       = "YAML"
	CodeLanguageXML        = "XML"
	CodeLanguagePython     = "Python"
	CodeLanguageGolang     = "Golang"
	CodeLanguageJava       = "Java"
	CodeLanguageJavaScript = "JavaScript"
	CodeLanguageRuby       = "Ruby"
	CodeLanguageShell      = "Shell"
	CodeLanguageDotNet     = "Shell"
	CodeLanguageProperties = "properties"
)

type CodeEditor struct {
	editor *gvcode.Editor
	// popup  completion.CompletionPopup

	code string

	theme *apptheme.Theme

	styledCode string
	styles     []*gvcode.TextStyle

	lexer     chroma.Lexer
	codeStyle *chroma.Style

	lang string

	onChange func(text string)

	font font.FontFace

	border widget.Border

	beatufier   widget.Clickable
	loadExample widget.Clickable

	withBeautify bool

	onLoadExample func()

	vScrollbar      widget.Scrollbar
	vScrollbarStyle material.ScrollbarStyle
}

func NewCodeEditor(code string, lang string, theme *apptheme.Theme) *CodeEditor {

	c := &CodeEditor{
		theme:  theme,
		editor: &gvcode.Editor{},
		code:   code,
		lang:   lang,
	}

	// cm := &completion.DefaultCompletion{Editor: c.editor}
	// // set completion triggers
	// cm.SetTriggers(gvcode.AutoTrigger{})
	//
	// // set popup widget to let user navigate the candidates.
	// c.popup = *completion.NewCompletionPopup(c.editor, cm)
	// cm.SetPopup(func(gtx layout.Context, items []gvcode.CompletionCandidate) layout.Dimensions {
	//	c.popup.TextSize = unit.Sp(12)
	//	c.popup.Size = image.Point{
	//		X: gtx.Dp(unit.Dp(400)),
	//		Y: gtx.Dp(unit.Dp(200)),
	//	}
	//
	//	return c.popup.Layout(gtx, theme.Material(), items)
	// })

	c.editor.WithOptions(
		gvcode.WithShaperParams(c.font.Font, unit.Sp(12), text.Start, unit.Sp(16), 1),
		gvcode.WithTabWidth(4),
		gvcode.WithSoftTab(true),
		gvcode.WrapLine(true),
		// 	gvcode.WithAutoCompletion(cm),
	)
	c.vScrollbarStyle = material.Scrollbar(theme.Material(), &c.vScrollbar)

	c.border = widget.Border{
		Color:        theme.BorderColor,
		Width:        unit.Dp(1),
		CornerRadius: unit.Dp(4),
	}

	c.lexer = getLexer(lang)

	style := styles.Get("dracula")
	if style == nil {
		style = styles.Fallback
	}

	c.codeStyle = style

	c.editor.SetText(code)

	return c
}

func getLexer(lang string) chroma.Lexer {
	lexer := lexers.Get(lang)
	if lexer == nil {
		lexer = lexers.Fallback
	}

	return chroma.Coalesce(lexer)
}

func (c *CodeEditor) WithBeautifier(enabled bool) {
	c.withBeautify = enabled
}

func (c *CodeEditor) SetOnChanged(f func(text string)) {
	c.onChange = f
}

func (c *CodeEditor) SetReadOnly(readOnly bool) {
	c.editor.WithOptions(gvcode.ReadOnlyMode(readOnly))
}

func (c *CodeEditor) SetOnLoadExample(f func()) {
	c.onLoadExample = f
}

func (c *CodeEditor) SetCode(code string) {
	c.editor.SetText(code)
	c.code = code
	c.editor.UpdateTextStyles(c.stylingText(c.editor.Text()))
}

func (c *CodeEditor) SetLanguage(lang string) {
	c.lang = lang
	c.lexer = getLexer(lang)
	c.editor.UpdateTextStyles(c.stylingText(c.editor.Text()))
}

func (c *CodeEditor) Code() string {
	return c.editor.Text()
}

func (c *CodeEditor) Layout(gtx layout.Context, theme *apptheme.Theme, hint string) layout.Dimensions {
	if c.styledCode == "" {
		// First time styling
		c.editor.UpdateTextStyles(c.stylingText(c.editor.Text()))
	}

	if !c.editor.ReadOnly() {
		if ev, ok := c.editor.Update(gtx); ok {
			if _, ok := ev.(gvcode.ChangeEvent); ok {
				st := c.stylingText(c.editor.Text())
				c.styles = st
				c.editor.UpdateTextStyles(st)
				if c.onChange != nil {
					c.onChange(c.editor.Text())
					c.code = c.editor.Text()
				}
			}
		}
	}

	if c.loadExample.Clicked(gtx) {
		c.onLoadExample()
	}

	flexH := layout.Flex{Axis: layout.Horizontal}

	if c.withBeautify {
		macro := op.Record(gtx.Ops)
		c.beautyButton(gtx, theme)
		defer op.Defer(gtx.Ops, macro.Stop())
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis:    layout.Horizontal,
				Spacing: layout.SpaceStart,
			}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if c.onLoadExample == nil {
						return layout.Dimensions{}
					}

					btn := widgets.Button(theme.Material(), &c.loadExample, widgets.RefreshIcon, widgets.IconPositionStart, "Load Example")
					btn.Color = theme.ButtonTextColor
					btn.Inset = layout.Inset{
						Top: unit.Dp(4), Bottom: unit.Dp(4),
						Left: unit.Dp(4), Right: unit.Dp(4),
					}

					return btn.Layout(gtx, theme)
				}),
			)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return c.border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return flexH.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{
							Top:    unit.Dp(4),
							Bottom: unit.Dp(4),
							Left:   unit.Dp(8),
							Right:  unit.Dp(4),
						}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return c.editorStyle(gtx, hint)
						})
					}),
				)
			})
		}),
	)
}

func (c *CodeEditor) beautyButton(gtx layout.Context, theme *apptheme.Theme) layout.Dimensions {
	if c.beatufier.Clicked(gtx) {
		c.SetCode(BeautifyCode(c.lang, c.code))
		if c.onChange != nil {
			c.onChange(c.editor.Text())
		}
	}

	return layout.Inset{Bottom: unit.Dp(4), Right: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.SE.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			btn := widgets.Button(theme.Material(), &c.beatufier, widgets.CleanIcon, widgets.IconPositionStart, "Beautify")
			btn.Color = theme.ButtonTextColor
			btn.Inset = layout.Inset{
				Top: 4, Bottom: 4,
				Left: 4, Right: 4,
			}
			return btn.Layout(gtx, theme)
		})
	})
}

func (c *CodeEditor) editorStyle(gtx layout.Context, _ string) layout.Dimensions {
	es := wgvcode.NewEditor(c.theme.Material(), c.editor)
	es.Font.Typeface = "Source Code Pro"
	es.TextSize = unit.Sp(12)
	es.LineHeightScale = 1.3

	es.SelectionColor = c.theme.TextSelectionColor
	editorDims := es.Layout(gtx)

	layout.E.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		viewportStart, viewportEnd := c.editor.ViewPortRatio()
		return c.vScrollbarStyle.Layout(gtx, layout.Vertical, viewportStart, viewportEnd)
	})

	if delta := c.vScrollbar.ScrollDistance(); delta != 0 {
		c.editor.ScrollByRatio(gtx, delta)
	}

	return editorDims
}

func (c *CodeEditor) stylingText(text string) []*gvcode.TextStyle {
	if c.styledCode == text {
		return c.styles
	}

	// nolint:prealloc
	var textStyles []*gvcode.TextStyle

	offset := 0

	iterator, err := c.lexer.Tokenise(nil, text)
	if err != nil {
		return textStyles
	}

	for _, token := range iterator.Tokens() {
		entry := c.codeStyle.Get(token.Type)

		textStyle := &gvcode.TextStyle{
			TextRange: gvcode.TextRange{
				Start: offset,
				End:   offset + len([]rune(token.Value)),
			},
			Color: rgbToOp(c.theme.Fg),
			// Background: rgbToOp(c.theme.Bg),
		}

		if entry.Colour.IsSet() {
			textStyle.Color = chromaColorToOp(entry.Colour)
		}

		textStyles = append(textStyles, textStyle)
		offset = textStyle.End
	}

	c.styledCode = text
	c.styles = textStyles

	return textStyles
}

func chromaColorToOp(textColor chroma.Colour) op.CallOp {
	ops := new(op.Ops)

	m := op.Record(ops)
	paint.ColorOp{Color: color.NRGBA{
		R: textColor.Red(),
		G: textColor.Green(),
		B: textColor.Blue(),
		A: 0xff,
	}}.Add(ops)
	return m.Stop()
}

func rgbToOp(color color.NRGBA) op.CallOp {
	ops := new(op.Ops)

	m := op.Record(ops)
	paint.ColorOp{Color: color}.Add(ops)
	return m.Stop()
}
