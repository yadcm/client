package components

import (
	"regexp"
	"strings"

	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

const passwdMask string = "*"

type TextEdit struct {
	*widgets.Paragraph
	value    string
	password bool
	focus    bool
}

func (t *TextEdit) Focus() {
	t.BorderStyle.Modifier = termui.ModifierBold
	t.BorderStyle.Fg = termui.ColorMagenta
	t.focus = true
}
func (t *TextEdit) Unfocus() {
	t.BorderStyle.Modifier = termui.ModifierClear
	t.BorderStyle.Fg = termui.ColorWhite
	t.focus = false
}
func (t *TextEdit) PopValue() {
	valueLen := len(t.value)
	if valueLen == 0 {
		return
	}
	t.value = t.value[:valueLen-1]
	t.formatText()
}
func (t *TextEdit) WriteValue(value string) {
	t.value += value
	t.formatText()
}
func (t *TextEdit) Clear() {
	t.value = ""
	t.formatText()
}
func (t *TextEdit) Value() string {
	return t.value
}
func (t *TextEdit) formatText() {
	if !t.password {
		t.Text = t.value
		return
	}
	m := regexp.MustCompile(".")
	t.Text = m.ReplaceAllString(t.value, passwdMask)
}

func (t *TextEdit) OnInput(key string) {
	if !t.focus {
		return
	}
	switch key {
	case "<Backspace>":
		t.PopValue()
		return
	case "<Space>":
		t.WriteValue(" ")
		return
	case "<Escape>", "<Enter>":
		t.Unfocus()
	default:
		if strings.HasPrefix(key, "<") && len(key) > 1 {
			return
		}
		t.WriteValue(key)
	}
}

func NewTextEdit(title string, password bool) *TextEdit {
	var view TextEdit
	view.Paragraph = widgets.NewParagraph()
	view.Title = title
	view.password = password
	return &view
}
