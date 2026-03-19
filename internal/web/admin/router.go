//  Copyright(C) 2026 github.com/hidu  All Rights Reserved.
//  Author: hidu <duv123+git@gmail.com>
//  Date: 2026-03-13

package admin

import (
	"github.com/xanygo/anygo/store/xsession"
	"github.com/xanygo/anygo/xhttp"
	"github.com/xanygo/anygo/xi18n"

	"github.com/xanygo/aimux/internal/resource"
	"github.com/xanygo/aimux/internal/resource/i18n"
)

func Router(router *xhttp.Router) {
	sh := &xsession.HTTPHandler{
		NewStorage: resource.SessionStorage(),
	}
	router.Use(sh.Next)

	xih := xi18n.HTTPHandler{
		Bundle: i18n.Resource,
	}
	router.Use(xih.Next)

	dashboard.Router = router
	dashboard.InitOnce()

	router.GetFunc("/", index)
	xhttp.RegisterGroup(router, "/api/", &apiHandler{})
	xhttp.RegisterGroup(router, "/user/", &userHandler{})
}
