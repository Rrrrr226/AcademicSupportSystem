package web

import (
	"HelpStudent/pkg/utils/page"
	"github.com/flamego/flamego"
)

func InjectPaginate() flamego.Handler {
	return func(r flamego.Render, c flamego.Context) {
		var req page.Paginate
		req.Current = c.QueryInt("page")
		req.PageSize = c.QueryInt("pageSize")
		c.Map(req)
	}
}
