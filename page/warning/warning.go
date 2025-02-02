package warning

import (
	"anshopper/icon"
	"anshopper/page"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

type Page struct {
	widget.List
	*page.Router
}

var uID string
var rstr = ""

func New(router *page.Router, userID string) *Page {
	uID = userID
	return &Page{
		Router: router,
	}
}

var _ page.Page = &Page{}

func (p *Page) Actions() []component.AppBarAction {
	return []component.AppBarAction{}
}

func (p *Page) Overflow() []component.OverflowAction {
	return []component.OverflowAction{}
}

func (p *Page) NavItem() component.NavItem {
	return component.NavItem{
		Name: "Warnings",
		Icon: icon.OtherIcon,
	}
}

var list widget.List

func (p *Page) Layout(gtx C, th *material.Theme) D {

	p.List.Axis = layout.Vertical
	return material.List(th, &p.List).Layout(gtx, 1, func(gtx C, _ int) D {

		return layout.Center.Layout(gtx, func(gtx C) D {
			return layout.Flex{
				Axis:      layout.Vertical,
				Alignment: layout.Start,
			}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Axis:      layout.Vertical,
						Alignment: layout.Start,
					}.
						Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								ptxt := material.Label(th, unit.Sp(18), "1. We not accept extra larger packages")
								return layout.Inset{
									Top: unit.Dp(14),
								}.Layout(gtx, ptxt.Layout)
							}),
						)
				}),
			)
		})
	})
}
