package link

import "github.com/qlova/seed"

type Widget struct {
	seed.Seed
}

func New(url ...string) Widget {
	widget := seed.New()
	widget.SetTag("a")
	
	if len(url) > 0 {
		widget.SetAttributes("href='"+url[0]+"'")
	} else {
		widget.SetAttributes("href='#'")
	}

	return Widget{widget}
}

func AddTo(parent seed.Interface, url ...string) Widget {
	var widget = New(url...)
	parent.Root().Add(widget)
	return widget
}
