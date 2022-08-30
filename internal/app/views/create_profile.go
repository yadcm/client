package views

import (
	"errors"
	"os"
	"yadcmc/internal/app/components"
	"yadcmc/internal/app/crypto"
	"yadcmc/internal/app/model"

	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/vmihailenco/msgpack/v5"
)

var (
	errPasswordsMissmatch error = errors.New("passwords does not match")
)

type CreateProfileView struct {
	*View
	editing             bool
	msgContainer        *widgets.Paragraph
	inputAddr           *components.TextEdit
	inputUsername       *components.TextEdit
	inputPassword       *components.TextEdit
	inputRepeatPassword *components.TextEdit
}

func (v *CreateProfileView) OnUnmount() {
	v.inputUsername.Unfocus()
	v.inputPassword.Unfocus()
}

func (v *CreateProfileView) OnInput(key string) {
	defer v.Render()
	if v.editing {
		// @todo need better way to connect events
		v.inputAddr.OnInput(key)
		v.inputPassword.OnInput(key)
		v.inputUsername.OnInput(key)
		v.inputRepeatPassword.OnInput(key)
	}
	switch key {
	case "A":
		if v.editing {
			return
		}
		v.editing = true
		v.inputAddr.Focus()
	case "U":
		if v.editing {
			return
		}
		v.editing = true
		v.inputUsername.Focus()
	case "P":
		if v.editing {
			return
		}
		v.editing = true
		v.inputPassword.Focus()
	case "R":
		if v.editing {
			return
		}
		v.editing = true
		v.inputRepeatPassword.Focus()
	case "<Escape>", "<Enter>":
		if v.editing {
			v.editing = false
			return
		}
		if err := v.validatePassword(); err != nil {
			v.msgContainer.Text = err.Error()
			v.Render()
			return
		}
		v.createProfile()
	}
}

func (v *CreateProfileView) validatePassword() error {
	if v.inputPassword.Value() != v.inputRepeatPassword.Value() {
		v.inputPassword.Clear()
		v.inputRepeatPassword.Clear()
		return errPasswordsMissmatch
	}
	return nil
}

func (v *CreateProfileView) createProfile() {
	var profile model.Profile
	profile.Tag = v.inputUsername.Value()
	profile.Host = v.inputAddr.Value()
	data, errMarshall := msgpack.Marshal(profile)
	if errMarshall != nil {
		v.msgContainer.Text = errMarshall.Error()
		v.Render()
		return
	}
	key := crypto.KeyFromPassword(v.inputPassword.Value())
	encrypted, errEncrypt := crypto.AesCBCEncrypt(data, key)
	if errEncrypt != nil {
		v.msgContainer.Text = errEncrypt.Error()
		v.Render()
		return
	}

	if errWrite := os.WriteFile(dotProfile, encrypted, os.ModePerm); errWrite != nil {
		v.msgContainer.Text = errWrite.Error()
		v.Render()
		return
	}
	//@todo push profile to daemon
}

func NewCreateProfileView(_ ViewArgs) Mountable {
	var grid *termui.Grid = termui.NewGrid()
	var view CreateProfileView
	// @todo need a better way to do it
	view.View = NewView(grid, view.OnInput, view.OnUnmount)

	view.inputAddr = components.NewTextEdit("<A>DDRESS", false)
	view.inputUsername = components.NewTextEdit("<U>SERNAME", false)
	view.inputPassword = components.NewTextEdit("<P>ASSWORD", true)
	view.inputRepeatPassword = components.NewTextEdit("<R>EPEAT", true)
	view.msgContainer = widgets.NewParagraph()
	view.msgContainer.Border = false
	grid.Set(
		termui.NewRow(1.0/5,
			termui.NewCol(1.0/1, view.inputAddr),
		),
		termui.NewRow(1.0/5,
			termui.NewCol(1.0/1, view.inputUsername),
		),
		termui.NewRow(1.0/5,
			termui.NewCol(1.0/1, view.inputPassword),
		),
		termui.NewRow(1.0/5,
			termui.NewCol(1.0/1, view.inputRepeatPassword),
		),
		termui.NewRow(1.0/5,
			termui.NewCol(1.0/1, view.msgContainer),
		),
	)
	return view
}
