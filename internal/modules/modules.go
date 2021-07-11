package modules

import (
	_ "github.com/RindangRamadhan/mnc-interview/internal/modules/api"
	"github.com/RindangRamadhan/mnc-interview/internal/router"
)

func init() {
	router.MiddlewareInit()
}
