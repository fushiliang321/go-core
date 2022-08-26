package routers

import (
	"core/router"
)

var routers = router.New()

func Set(config *router.Router) {
	routers = config
}
func Get() *router.Router {
	return routers
}
