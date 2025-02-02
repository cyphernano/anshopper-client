package page

import (
	// "anshopper/database"
	"anshopper/icon"
	"log"
	"time"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget/material"
	"gioui.org/x/component"
)

type Page interface {
	Actions() []component.AppBarAction
	Overflow() []component.OverflowAction
	Layout(gtx layout.Context, th *material.Theme) layout.Dimensions
	NavItem() component.NavItem
}

type Router struct {
	page    map[interface{}]Page
	current interface{}
	*component.ModalNavDrawer
	NavAnim component.VisibilityAnimation
	*component.AppBar
	*component.ModalLayer
}

func NewRouter(userID string) Router {
	modal := component.NewModal()
	nav := component.NewNav(userID, "This is your ID")
	modalNav := component.ModalNavFrom(&nav, modal)
	bar := component.NewAppBar(modal)
	bar.NavigationIcon = icon.HomeIcon

	na := component.VisibilityAnimation{
		State:    component.Invisible,
		Duration: time.Millisecond * 250,
	}

	return Router{
		page:           make(map[interface{}]Page),
		ModalLayer:     modal,
		ModalNavDrawer: modalNav,
		AppBar:         bar,
		NavAnim:        na,
	}
}

func (r *Router) Register(tag interface{}, p Page) {
	r.page[tag] = p
	navItem := p.NavItem()
	navItem.Tag = tag
	if r.current == interface{}(nil) {
		r.current = tag
		r.AppBar.Title = navItem.Name
		r.AppBar.SetActions(p.Actions(), p.Overflow())
	}
	r.ModalNavDrawer.AddNavItem(navItem)
}

func (r *Router) SwitchTo(tag interface{}) {
	p, ok := r.page[tag]
	if !ok {
		return
	}
	navItem := p.NavItem()
	r.current = tag
	r.AppBar.Title = navItem.Name
}

func (r *Router) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	for _, event := range r.AppBar.Events(gtx) {
		switch event := event.(type) {
		case component.AppBarNavigationClicked:
			r.ModalNavDrawer.Appear(gtx.Now)
			r.NavAnim.Disappear(gtx.Now)
		case component.AppBarContextMenuDismissed:
			log.Printf("Context menu dismissed: %v", event)
		}
	}
	if r.ModalNavDrawer.NavDestinationChanged() {
		r.SwitchTo(r.ModalNavDrawer.CurrentNavDestination())
	}
	paint.Fill(gtx.Ops, th.Palette.Bg)
	content := layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Max.X /= 3
				return r.NavDrawer.Layout(gtx, th, &r.NavAnim)
			}),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return r.page[r.current].Layout(gtx, th)
			}),
		)
	})
	bar := layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return r.AppBar.Layout(gtx, th, "Menu", "Actions")
	})
	flex := layout.Flex{Axis: layout.Vertical}
	flex.Layout(gtx, bar, content)

	r.ModalLayer.Layout(gtx, th)
	return layout.Dimensions{Size: gtx.Constraints.Max}
}
