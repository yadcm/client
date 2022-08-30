package views

import (
	"yadcmc/internal/app/components"
	"yadcmc/internal/app/crypto"
	"yadcmc/internal/app/model"

	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/vmihailenco/msgpack/v5"
)

type UnlockProfileView struct {
	*View
	editing       bool
	profileRaw    []byte
	msgContainer  *widgets.Paragraph
	inputPassword *components.TextEdit
}

func (v *UnlockProfileView) OnUnmount() {
	v.inputPassword.Unfocus()
}

func (v *UnlockProfileView) OnInput(key string) {
	defer v.Render()
	if v.editing {
		v.inputPassword.OnInput(key)
	}
	switch key {
	case "P":
		if v.editing {
			return
		}
		v.editing = true
		v.inputPassword.Focus()
	case "<Escape>", "<Enter>":
		if v.editing {
			v.editing = false
			return
		}
		v.decryptProfile()
	}
}

func (v *UnlockProfileView) decryptProfile() {
	v.msgContainer.Text = "Trying to decrypt file"
	v.Render()
	key := crypto.KeyFromPassword(v.inputPassword.Value())
	decrypted, err := crypto.AesCBCDecrypt(v.profileRaw, key)
	if err != nil {
		v.msgContainer.Text = err.Error()
		v.Render()
		return
	}

	var profile model.Profile
	if err := msgpack.Unmarshal(decrypted, &profile); err != nil {
		//@todo invalid password?
		v.msgContainer.Text = err.Error()
		v.Render()
		return
	}

	v.goToMainView(profile)
}

func (v *UnlockProfileView) goToMainView(profile model.Profile) {
	//@todo
}

func NewUnlockProfileView(args ViewArgs) Mountable {
	if len(args) == 0 {
		return nil
	}
	profileRaw, casted := args[0].([]byte)
	if !casted {
		return nil
	}
	var grid *termui.Grid = termui.NewGrid()
	var view UnlockProfileView
	view.profileRaw = profileRaw
	// @todo need a better way to do it
	view.View = NewView(grid, view.OnInput, view.OnUnmount)

	view.inputPassword = components.NewTextEdit("<P>ASSWORD", true)
	view.msgContainer = widgets.NewParagraph()
	view.msgContainer.Border = false
	grid.Set(
		termui.NewRow(1.0/5,
			termui.NewCol(1.0/1, view.inputPassword),
		),
		termui.NewRow(1.0/5,
			termui.NewCol(1.0/1, view.msgContainer),
		),
	)
	return view
}
