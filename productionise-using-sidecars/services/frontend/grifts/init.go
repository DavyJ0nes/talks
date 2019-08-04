package grifts

import (
	"github.com/davyj0nes/talks/productionise-using-sidecars/frontend/actions"
	"github.com/gobuffalo/buffalo"
)

func init() {
	buffalo.Grifts(actions.App())
}
