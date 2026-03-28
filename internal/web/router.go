//  Copyright(C) 2025 github.com/hidu  All Rights Reserved.
//  Author: hidu <duv123+git@gmail.com>
//  Date: 2025-09-19

package web

import (
	"net/http"

	"github.com/xanygo/anygo/xhttp"
	"github.com/xanygo/anygo/xhttp/xhandler"
	"github.com/xanygo/anygo/xlog"

	"github.com/xanygo/aimux/internal/apigate"
	"github.com/xanygo/aimux/internal/config"
	"github.com/xanygo/aimux/internal/web/admin"
)

func initAdminRouter(withAPI bool) *xhttp.Router {
	router := xhttp.NewRouter()
	registerMiddleware(router)

	adminGroup := router.Prefix(config.AdminPath())
	admin.Router(adminGroup)

	if withAPI {
		router.Handle("/*", apigate.Default())
	} else {
		router.GetFunc("/", func(w http.ResponseWriter, req *http.Request) {
			http.Redirect(w, req, "/admin/", http.StatusFound)
		})
	}
	return router
}

func initAPIRouter() *xhttp.Router {
	router := xhttp.NewRouter()
	registerMiddleware(router)
	router.Handle("/*", apigate.Default())
	return router
}

func registerMiddleware(router *xhttp.Router) {
	al := &xhandler.AccessLog{
		Logger: xlog.AccessLogger(),
		OnCookies: func(cookies []*http.Cookie) []xlog.Attr {
			return nil
		},
		OnHeaders: func(h http.Header) []xlog.Attr {
			return nil
		},
	}
	router.Use(al.Next)
}
