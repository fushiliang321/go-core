package routers

import (
	"github.com/fushiliang321/go-core/router"
)

var routers = router.New()

func Set(config *router.Router) {
	routers = config
}
func Get() *router.Router {
	return routers
}
