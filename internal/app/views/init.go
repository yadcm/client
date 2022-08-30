package views

import (
	"os"

	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

const dotProfile string = ".profile"

type InitialView struct {
	*View
	p *widgets.Paragraph
}

func NewInitialView(_ ViewArgs) Mountable {
	var grid *termui.Grid = termui.NewGrid()
	var view InitialView = InitialView{
		View: NewView(grid, nil, nil),
	}

	view.p = widgets.NewParagraph()
	view.p.Title = "LOADING"
	grid.Set(termui.NewRow(1.0, view.p))
	go view.getLocalProfile()
	return view
}

func (v *InitialView) getLocalProfile() {
	<-v.mounted
	v.p.Text = "Trying to open .profile"
	v.Render()
	data, err := os.ReadFile(dotProfile)
	if err != nil {
		v.p.Text = err.Error()
		v.Render()
		v.Navigate(PresentView{
			View: NewCreateProfileView,
		})
		return
	}
	v.Navigate(PresentView{
		View: NewUnlockProfileView,
		Args: []interface{}{data},
	})
}
