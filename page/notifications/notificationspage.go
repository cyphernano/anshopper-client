package notifications

import (
	"anshopper/database"
	"anshopper/icon"
	"anshopper/page"
	"bytes"
	"fmt"
	"image"
	"io"
	"log"
	"strings"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/supabase-community/postgrest-go"
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

type Page struct {
	db      database.Export
	idb     *postgrest.Client
	syncBtn widget.Clickable
	sel     []widget.Selectable
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
		"notifications",
		uID,
		[]string{
			"ref_order",
			"state",
			"amount",
			"pay",
		},
	)
	return r
}

func New(router *page.Router, userID, notfs string) *Page {

	uID = userID
	rstr = notfs
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
		Name: "Notifications",
		Icon: icon.VisibilityIcon,
	}
}

func splitKV(data string) string {
	_, v, _ := strings.Cut(data, ":")
	return v
}

var list widget.List
var matrix [][]string

type nopCloser struct {
	io.Writer
}

func (nopCloser) Close() error { return nil }

func genereteQr(qrdata string, gtx C) widget.Image {
	qrc, err := qrcode.NewWith(qrdata,
		qrcode.WithEncodingMode(qrcode.EncModeByte),
		qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionMedium),
	)

	if err != nil {
		panic(err)
	}

	buf := bytes.NewBuffer(nil)
	wr := nopCloser{Writer: buf}
	w2 := standard.NewWithWriter(wr, standard.WithQRWidth(5))
	if err = qrc.Save(w2); err != nil {
		panic(err)
	}
	img, _, err := image.Decode(bytes.NewReader(buf.Bytes()))
	if err != nil {
		log.Println(err)
	}
	im := widget.Image{Src: paint.NewImageOp(img)}
	im.Src.Filter = paint.FilterNearest
	im.Scale = float32(0.8)
	im.Src.Add(gtx.Ops)
	return im
}

func notformat(data string) []string {
	i := strings.ReplaceAll(data, " \n ", "")
	ii := strings.ReplaceAll(i, "[{", "")
	iii := strings.ReplaceAll(ii, "}]", "")
	iiii := strings.ReplaceAll(iii, "{", "")
	iiiii := strings.ReplaceAll(iiii, ", ", "; \n")
	val := strings.Split(iiiii, "},")
	return val
}

func (p *Page) populateNotifications(data string) {

	val := notformat(data)

	rows, cols := len(val), 4
	var sa [][]string

	var slct = make([]widget.Selectable, rows)

	for i := 0; i < rows; i++ {
		row := make([]string, cols)
		sa = append(sa, row)
	}

	for j := range len(sa) {
		js := strings.Split(val[j], ",")
		for jj := 0; jj < len(js); jj++ {
			_, v, _ := strings.Cut(js[jj], ":")
			sa[j][jj] = strings.ReplaceAll(v, "\"", "")
			slct[j] = widget.Selectable{}
		}
	}
	matrix = sa
	p.sel = append(p.sel, slct...)
}

func insetTextOfCard(gtx C, txt layout.Widget) D {
	return layout.Inset{Top: unit.Dp(2), Left: unit.Dp(35)}.
		Layout(gtx, txt)
}

func (p *Page) Layout(gtx C, th *material.Theme) D {

	if rstr != "" {
		p.populateNotifications(rstr)
	}

	p.List.Axis = layout.Vertical
	return material.List(th, &p.List).Layout(gtx, 1, func(gtx C, _ int) D {

		if p.syncBtn.Clicked(gtx) {
			rstr = p.ActSelectDB(uID)
			if rstr != "" {
				p.populateNotifications(rstr)
			}
		}

		return layout.Center.Layout(gtx, func(gtx C) D {
			return layout.Flex{
				Axis:      layout.Vertical,
				Alignment: layout.Middle,
			}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Axis:      layout.Horizontal,
						Alignment: layout.Middle,
						Spacing:   25}.
						Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								order := material.H6(th, "List of notifications")
								return layout.Inset{
									Top:    unit.Dp(25),
									Left:   unit.Dp(100),
									Bottom: unit.Dp(10),
									Right:  unit.Dp(100)}.
									Layout(gtx, order.Layout)
							}),
							layout.Rigid(func(gtx C) D {
								sbtn := material.IconButton(th, &p.syncBtn, icon.PlusIcon, "sync")
								sbtn.Size = unit.Dp(3)
								return layout.Inset{
									Top:    unit.Dp(20),
									Left:   unit.Dp(50),
									Bottom: unit.Dp(10),
									Right:  unit.Dp(2)}.
									Layout(gtx,
										sbtn.Layout,
									)
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
								Axis:      layout.Vertical,
								Alignment: layout.Start,
							}.
								Layout(gtx,
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										var order material.LabelStyle
										var val string = matrix[idx][0]
										val = strings.ReplaceAll(val, ";", "")
										if val != "" {
											order = material.Body1(th, "Order: "+val)
										}
										return layout.Inset{
											Top:  unit.Dp(16),
											Left: unit.Dp(35),
										}.
											Layout(gtx, order.Layout)
									}),
									layout.Rigid(func(gtx C) D {
										var order material.LabelStyle
										var val string = matrix[idx][1]
										if val != "" {
											order = material.Body1(th, "State: "+val)
										}
										return insetTextOfCard(gtx, order.Layout)
									}),
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										var order material.LabelStyle
										var val string = matrix[idx][2]
										val = strings.ReplaceAll(val, ";", "")
										if val != "" {
											order = material.Body1(th, "Amount: "+val)
										}
										return insetTextOfCard(gtx, order.Layout)
									}),
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										var order material.LabelStyle
										var val string = matrix[idx][3]
										if val != "" {
											nv := fmt.Sprint("To pay: ", val)
											order = material.Body1(th, nv)
											order.State = &p.sel[idx]
										}
										return insetTextOfCard(gtx, order.Layout)
									}),
									layout.Rigid(func(gtx C) D {
										var img widget.Image
										var val string = matrix[idx][3]
										if val != "" {
											img = genereteQr(val, gtx)
										}
										return img.Layout(gtx)
									}),
								) // Layout.Flex
						}) // Layout list
					return screenLaoyut
				}), // top parent Rigid
			)
		})
	})
}
