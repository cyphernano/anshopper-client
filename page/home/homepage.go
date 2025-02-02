package home

import (
	"anshopper/database"
	"anshopper/icon"
	page "anshopper/page"
	"fmt"
	"image/color"
	"log"
	s "strings"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"

	postgrest "github.com/supabase-community/postgrest-go"
)

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var p = fmt.Printf

var uID string = ""

type (
	C = layout.Context
	D = layout.Dimensions
)

type Page struct {
	linkField     component.TextField
	DeliveryField component.TextField
	DescField     component.TextField
	submitButton  widget.Clickable
	db            database.Export
	idb           *postgrest.Client
	swt           widget.Bool
	widget.List
	*page.Router
}

func New(router *page.Router, userID string) *Page {
	uID = userID
	return &Page{
		Router: router,
	}
}

var afterSubmit string = "Insert"
var _ page.Page = &Page{}

func (p *Page) Actions() []component.AppBarAction {
	return []component.AppBarAction{}
}

func (p *Page) Overflow() []component.OverflowAction {
	return []component.OverflowAction{}
}

func (p *Page) NavItem() component.NavItem {
	return component.NavItem{
		Name: "Home",
		Icon: icon.HomeIcon,
	}
}

func (p *Page) Layout(gtx C, th *material.Theme) D {
	p.List.Axis = layout.Vertical

	if uID == "" {
		uID = p.db.GetUserID()
	}

	if p.idb == nil {
		p.idb = p.db.InitDB("34.170.73.36", "3000", "public")
	}

	hdclr := color.NRGBA{R: 99, G: 35, B: 210, A: 255}
	noclr := color.NRGBA{R: 175, G: 175, B: 175, A: 255}
	var linkT, DeliveryT, DescT, crypto = "", "", "", ""

	var str string

	if p.submitButton.Clicked(gtx) {
		afterSubmit = "Check"
		if lt := p.linkField.Text(); lt != "" {
			linkT = lt
		}
		if dyt := p.DeliveryField.Text(); dyt != "" {
			dyt2 := s.ReplaceAll(dyt, ",", " -")
			DeliveryT = dyt2
		}
		if dt := p.DescField.Text(); dt != "" {
			dt2 := s.ReplaceAll(dt, ",", " -")
			DescT = dt2
		}
		if !p.swt.Value {
			crypto = "nano"
		} else {
			crypto = "monero"
		}

		str = ("uuid: " + uID + ", " +
			"link: " + linkT + ", " +
			"description: " + DescT + ", " +
			"delivery_address: " + DeliveryT + ", " +
			"crypto: " + crypto + ", " +
			"txid: ")

		p.db.InsertDB(p.idb, str)
	}

	if afterSubmit == "Check" {
		r := p.db.CheckPostDB(p.idb, str)
		if r == false {
			p.linkField.Clear()
			p.DeliveryField.Clear()
			p.DescField.Clear()
			afterSubmit = "Done"
		} else {
			afterSubmit = "Fail"
		}
	}
	return material.List(th, &p.List).Layout(gtx, 1, func(gtx C, _ int) D {

		return layout.Flex{
			Alignment: layout.Middle,
			Spacing:   8,
			Axis:      layout.Vertical,
		}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Center.Layout(gtx, func(gtx C) D {
					gtx.Constraints.Max.X = gtx.Dp(unit.Dp(300))
					return layout.Flex{
						Axis:      layout.Vertical,
						Spacing:   8,
						Alignment: layout.Middle,
					}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							r := material.Label(th, unit.Sp(8), "")
							return layout.Inset{
								Top:    unit.Dp(50),
								Bottom: unit.Dp(50),
							}.Layout(gtx,
								r.Layout,
							)
						}),
						layout.Rigid(func(gtx C) D {
							title := material.Body1(th, "Submit your order")
							title.Color = hdclr
							title.TextSize = unit.Sp(32)
							return layout.Inset{Bottom: 14}.
								Layout(gtx,
									title.Layout,
								)
						}),
						layout.Rigid(func(gtx C) D {
							return p.linkField.Layout(gtx, th, "Link of product")
						}),
						layout.Rigid(func(gtx C) D {
							return p.DeliveryField.Layout(gtx, th, "Delivery address")
						}),
						layout.Rigid(func(gtx C) D {
							return p.DescField.Layout(gtx, th, "Description")
						}),

						layout.Rigid(func(gtx C) D {
							return layout.Flex{
								Axis:    layout.Horizontal,
								Spacing: 8}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									crypto := material.Body2(th, "Nano")
									if !p.swt.Value {
										crypto.Color = hdclr
									} else {
										crypto.Color = noclr
									}
									return layout.Inset{
										Left: unit.Dp(8),
										Top:  unit.Dp(18),
									}.Layout(gtx,
										crypto.Layout,
									)
								}),
								layout.Rigid(func(gtx C) D {
									return layout.Inset{
										Left: unit.Dp(8),
										Top:  unit.Dp(18),
									}.Layout(gtx,
										material.Switch(th, &p.swt, "Select crypto").Layout,
									)
								}),
								layout.Rigid(func(gtx C) D {
									crypto := material.Body2(th, "Monero")
									if p.swt.Value {
										crypto.Color = hdclr
									} else {
										crypto.Color = noclr
									}
									return layout.Inset{
										Left: unit.Dp(8),
										Top:  unit.Dp(18),
									}.Layout(gtx,
										crypto.Layout,
									)
								}),
								layout.Rigid(func(gtx C) D {
									var buttonWithResponse layout.Dimensions
									if afterSubmit == "Check" {
										gtx.Disabled()
									}
									if lt := p.linkField.Text(); lt != "" {
										afterSubmit = "Insert"
									}
									buttonWithResponse = layout.Inset{
										Left: unit.Dp(60),
										Top:  unit.Dp(8),
									}.Layout(gtx,
										material.Button(th, &p.submitButton, afterSubmit).Layout)
									return buttonWithResponse
								}),
							)
						}),
					)
				})
			}),
		)
	})

}
