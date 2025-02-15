package orders

import (
	"anshopper/database"
	"anshopper/icon"
	"anshopper/page"
	"image/color"
	"log"
	"strconv"
	"strings"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/supabase-community/postgrest-go"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

type Page struct {
	db        database.Export
	idb       *postgrest.Client
	syncBtn   widget.Clickable
	txid      []component.TextField
	submitBtn []widget.Clickable
	widget.List
	*page.Router
}

var uID string
var rstr = ""

func (p *Page) ActSelectDB(userID string) string {
	if userID == "" {
		uID = p.db.GetUserID()
	} else {
		uID = userID
	}

	p.idb = p.db.InitDB("34.170.73.36", "3000", "public")
	r := p.db.SelectDB(
		p.idb,
		"orders",
		uID,
		[]string{
			"id",
			"link",
			"description",
			"delivery_address",
			"crypto",
			"txid",
		},
	)
	return r
}

func (p *Page) ActUpdateDB(id string, toUp map[string]string) {
	if uID == "" {
		uID = p.db.GetUserID()
	}
	p.idb = p.db.InitDB("34.170.73.36", "3000", "public")
	_, e := p.db.UpdateDB(
		p.idb,
		"orders",
		uID,
		id,
		toUp,
	)
	if e != nil {
		log.Printf("error from clientPostgrest: %s\n", e)
	}
}

func New(router *page.Router, userID, ords string) *Page {
	uID = userID
	rstr = ords
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
		Name: "Orders",
		Icon: icon.AccountBoxIcon,
	}
}

func splitKV(data string) string {
	_, v, _ := strings.Cut(data, ":")
	return v
}

var list widget.List
var matrix [][]string
var afterSubmit string = "Submit"
var mtx []string

func ordformat(data string) []string {
	i := strings.ReplaceAll(data, " \n ", "")
	ii := strings.ReplaceAll(i, "[{", "")
	iii := strings.ReplaceAll(ii, "}]", "")
	iiii := strings.ReplaceAll(iii, "{", "")
	val := strings.Split(iiii, "},")
	return val
}

func (p *Page) populateOrders(data string) {

	val := ordformat(data)

	rows, cols := len(val), 6
	var sa [][]string

	var sbtn = make([]widget.Clickable, rows)
	var stxid = make([]component.TextField, rows)

	for i := 0; i < rows; i++ {
		row := make([]string, cols)
		sa = append(sa, row)
	}

	for j := range len(sa) {
		js := strings.Split(val[j], ",")
		for jj := 0; jj < len(js); jj++ {
			_, v, _ := strings.Cut(js[jj], ":")
			sa[j][jj] = strings.ReplaceAll(v, "\"", "")
			stxid[j] = component.TextField{}
			sbtn[j] = widget.Clickable{}
		}
	}
	matrix = sa
	p.txid = append(p.txid, stxid...)
	p.submitBtn = append(p.submitBtn, sbtn...)
}

func (p *Page) Layout(gtx C, th *material.Theme) D {

	if rstr != "" {
		p.populateOrders(rstr)
	}

	p.List.Axis = layout.Vertical
	return material.List(th, &p.List).Layout(gtx, 1, func(gtx C, _ int) D {

		if p.syncBtn.Clicked(gtx) {
			rstr = p.ActSelectDB(uID)
			if rstr != "" {
				p.populateOrders(rstr)
			}
		}

		if afterSubmit == "Check" {
			r := p.db.CheckUpdateTxidDB(p.idb, mtx)
			if r == true {
				ti, _ := strconv.Atoi(mtx[0])
				for i := range ti {
					p.txid[i].Clear()
				}
				afterSubmit = "Done"
				rstr = p.ActSelectDB(uID)
				if rstr != "" {
					p.populateOrders(rstr)
				}
			} else {
				afterSubmit = "Fail"
			}
		}

		return layout.Center.Layout(gtx, func(gtx C) D {
			gtx.Constraints.Max.X = gtx.Dp(unit.Dp(300))
			hdclr := color.NRGBA{R: 99, G: 35, B: 210, A: 255}
			return layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Axis:      layout.Horizontal,
						Alignment: layout.End,
					}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							order := material.H6(th, "List of orders")
							order.Font.Weight = 200
							order.Color = hdclr
							order.TextSize = unit.Sp(24)
							return layout.Inset{
								Top: unit.Dp(20),
							}.Layout(gtx, order.Layout)
						}),
						layout.Rigid(func(gtx C) D {
							sbtn := material.IconButton(th, &p.syncBtn, icon.PlusIcon, "sync")
							sbtn.Size = unit.Dp(1)
							return layout.Inset{
								Top:  unit.Dp(20),
								Left: unit.Dp(50),
							}.Layout(gtx, sbtn.Layout)
						}),
					)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					list.Alignment = layout.Start
					list.Axis = layout.Vertical
					screenLaoyut := material.List(th, &list).
						Layout(gtx, len(matrix), func(gtx C, i int) D {
							idx := (len(matrix) - 1) - i
							return layout.Flex{
								Axis: layout.Vertical,
							}.
								Layout(gtx,
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										var order material.LabelStyle
										var val = matrix[idx][1]
										if val != "" {
											order = material.Body1(th, "link: "+val)
										}
										return layout.Inset{
											Top: unit.Dp(16),
										}.Layout(gtx,
											order.Layout,
										)
									}),
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										var order material.LabelStyle
										var val = matrix[idx][2]
										if val != "" {
											order = material.Body1(th, "description: "+val)
										}
										return order.Layout(gtx)
									}),
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										var order material.LabelStyle
										var val = matrix[idx][3]
										if val != "" {
											order = material.Body1(th, "delivery_address: "+val)
										}
										return order.Layout(gtx)
									}),
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										var order material.LabelStyle
										var val = matrix[idx][4]
										if val != "" {
											order = material.Body1(th, "crypto: "+val)
										}
										return order.Layout(gtx)
									}),
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										v := matrix[idx][5]
										runetxid := []rune(v)
										var order material.LabelStyle
										var nothing layout.Dimensions
										if len(runetxid) > 5 {
											order = material.Body1(th, "txid: "+v)
										} else {
											order = material.Body1(th, "")
										}
										if order.Text != "" {
											return order.Layout(gtx)
										} else {
											return nothing
										}
									}),
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										var txidWidget layout.Dimensions
										v := matrix[idx][5]
										runetxid := []rune(v)
										if len(runetxid) < 5 {
											if p.submitBtn[idx].Clicked(gtx) {
												afterSubmit = "Check"
												up := p.txid[idx].Text()
												toUp := map[string]string{"txid": up}
												p.ActUpdateDB(matrix[idx][0], toUp)

												ti := len(matrix)
												ni := strconv.Itoa(ti)
												mtx = []string{ni, uID, up}
											}

											txidWidget = layout.Flex{
												Axis:      layout.Horizontal,
												Alignment: layout.Start,
											}.
												Layout(gtx,
													layout.Rigid(func(gtx C) D {
														gtx.Constraints.Max.X = gtx.Dp(unit.Dp(200))
														return p.txid[idx].Layout(gtx, th, "Transaction ID")
													}),
													layout.Rigid(func(gtx C) D {
														if afterSubmit == "Check" {
															gtx.Disabled()
														}

														sbtn := material.Button(th, &p.submitBtn[idx], afterSubmit)
														return layout.Inset{
															Top:  unit.Dp(10),
															Left: unit.Dp(2),
														}.Layout(gtx,
															sbtn.Layout,
														)
													}),
												)
										} else {
											txidWidget = layout.Dimensions{}
										}
										return txidWidget
									}),
								) // Layout.Flex
						}) // Layout list
					return screenLaoyut
				}), // top parent Rigid
			)
		})
	})
}
