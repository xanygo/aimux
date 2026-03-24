//  Copyright(C) 2025 github.com/hidu  All Rights Reserved.
//  Author: hidu <duv123+git@gmail.com>
//  Date: 2025-09-19

package web

import (
	"context"
	"net/http"
	"time"

	"github.com/xanygo/anygo"
	"github.com/xanygo/anygo/xhttp"
	"github.com/xanygo/anygo/xhttp/xhandler"
	"github.com/xanygo/anygo/xlog"
	"github.com/xanygo/anygo/xnet/xservice"

	"github.com/xanygo/aimux/internal/apigate"
	"github.com/xanygo/aimux/internal/factory"
	"github.com/xanygo/aimux/internal/web/admin"
)

var router = xhttp.NewRouter()

func routerRegister() {
	doInit()

	adminGroup := router.Prefix("/admin/")
	admin.Router(adminGroup)

	router.Handle("/*", apigate.Default())
}

func doInit() {
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

	loadApiGate()

	xservice.DefaultRegistry().Register(xservice.DefaultDummyService())
}

func loadApiGate() {
	var err error
	var ss apigate.Services
	for i := 0; i < 3; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		ss, err = factory.ServiceDao().GetAllActive(ctx)
		xlog.Info(ctx, "apigate.GetAllActive", xlog.ErrorAttr("error", err), xlog.Int("ss.len", len(ss)))
		cancel()
		if err == nil {
			err = apigate.Default().RegisterDny(ss)
			xlog.Info(ctx, "apigate.RegisterDny", xlog.ErrorAttr("error", err))
			anygo.Must(err)
			break
		}
	}
	anygo.Must(err)
}
