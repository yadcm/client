package main

import (
	"context"
	"log"
	"yadcmc/internal/app/views"

	ui "github.com/gizak/termui/v3"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() //@todo
	var nextView *views.PresentView = &views.PresentView{
		View: views.NewInitialView,
	}
	for {
		if nextView == nil {
			break
		}
		if nextView.View == nil {
			break
		}
		if view := nextView.View(nextView.Args); view != nil {
			nextView = view.Mount(ctx)
		}
	}
}
