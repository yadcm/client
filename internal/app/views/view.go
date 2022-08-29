package views

import (
	"context"
	"log"
	"sync"

	"github.com/gizak/termui/v3"
)

type Mountable interface {
	Mount(context.Context) *PresentView
	Unmount()
}

type ViewArgs []interface{}

type PresentView struct {
	View func(V ViewArgs) Mountable
	Args ViewArgs
}

type View struct {
	termui.Drawable
	ctx       context.Context
	done      context.CancelFunc
	mu        sync.Mutex
	mounted   chan struct{}
	nextView  *PresentView
	OnInput   func(string)
	OnUnmount func()
}

func NewView(
	dravable termui.Drawable,
	onInput func(string),
	onUnmount func(),
) *View {
	return &View{
		Drawable:  dravable,
		mounted:   make(chan struct{}),
		OnInput:   onInput,
		OnUnmount: onUnmount,
	}
}

func (v *View) Render() {
	v.mu.Lock()
	termui.Clear()
	termui.Render(v)
	v.mu.Unlock()
}

func (v *View) Mount(ctx context.Context) *PresentView {
	termWidth, termHeight := termui.TerminalDimensions()
	v.SetRect(0, 0, termWidth, termHeight)
	v.Render()
	v.ctx, v.done = context.WithCancel(ctx)
	close(v.mounted)
	v.handleEvents()
	return v.nextView
}

func (v *View) Unmount() {
	log.Println("Unmount")
	if v.done != nil {
		v.done()
	}
	if v.OnUnmount != nil {
		v.OnUnmount()
	}
}

func (v *View) Navigate(nextView PresentView) {
	v.nextView = &nextView
	log.Println("Navigate")
	v.Unmount()
}

func (v *View) handleEvents() {
	uiEvents := termui.PollEvents()
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "<C-c>":
				return
			case "<Resize>":
				payload := e.Payload.(termui.Resize)
				v.SetRect(0, 0, payload.Width, payload.Height)
				v.Render()
			default:
				if v.OnInput != nil {
					v.OnInput(e.ID)
				}
			}
		case <-v.ctx.Done():
			v.Unmount()
			return
		}
	}
}
