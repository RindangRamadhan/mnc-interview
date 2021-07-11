package modules

import (
	_ "github.com/RindangRamadhan/mnc-interview/internal/modules/api/example"
	"github.com/RindangRamadhan/mnc-interview/internal/router"
)

func init() {
	router.MiddlewareInit()
}
