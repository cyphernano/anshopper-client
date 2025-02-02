package applayout

import (
	"gioui.org/layout"
	"gioui.org/unit"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

type DetailRow struct {
	PrimaryWith float32
	layout.Inset
}

var DefaultInset = layout.UniformInset(unit.Dp(8))

func (d DetailRow) Layout(gtx C, primary, detail layout.Widget) D {
	if d.PrimaryWith == 0 {
		d.PrimaryWith = 0.3
	}
	if d.Inset == (layout.Inset{}) {
		d.Inset = DefaultInset
	}
	return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Flexed(d.PrimaryWith, func(gtx C) D {
			return d.Inset.Layout(gtx, primary)
		}),
		layout.Flexed(1-d.PrimaryWith, func(gtx C) D {
			return d.Inset.Layout(gtx, detail)
		}),
	)
}
