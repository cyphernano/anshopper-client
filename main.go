package main

import (
	"anshopper/database"
	page "anshopper/page"
	"anshopper/page/home"
	"anshopper/page/notifications"
	"anshopper/page/orders"
	"flag"
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget/material"
	"log"
	"os"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

var db database.Export
var orderP orders.Page
var notfsP notifications.Page

func main() {
	flag.Parse()
	go func() {
		w := new(app.Window)
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func loop(w *app.Window) error {
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	var ops op.Ops

	w.Option(app.Title("AnShopper"))
	userID := db.GetUserID()
	var notfs, ords string
	if userID == "f" {
		userID = db.GetUserID()
	}
	ords = orderP.ActSelectDB(userID)
	if ords != "" {
		notfs = notfsP.ActSelectDB(userID)
	}
	router := page.NewRouter(userID)
	router.Register(0, home.New(&router, userID))
	router.Register(1, orders.New(&router, userID, ords))
	router.Register(2, notifications.New(&router, userID, notfs))

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			router.Layout(gtx, th)
			e.Frame(gtx.Ops)
		}
	}
}
