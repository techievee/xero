package ctls

import (
	productServiceCmds "github.com/techievee/xero/productService/commands"
	"github.com/techievee/xero/xeroLog/debugcore"
)

type ProductsCtl struct {
	ServiceCommands *productServiceCmds.ProductsCmds
	Logger          debugcore.Logger
}
