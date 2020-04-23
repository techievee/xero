package commands

import (
	"github.com/techievee/xero/database"
	"github.com/techievee/xero/xeroLog/debugcore"
)

type ProductsCmds struct {
	DB     *database.DB
	Logger debugcore.Logger
}
