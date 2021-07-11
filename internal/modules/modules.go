package modules

import (
	_ "github.com/RindangRamadhan/mnc-interview/internal/modules/example"
	"github.com/RindangRamadhan/mnc-interview/internal/router"
)

func init() {
	router.MiddlewareInit()
}
