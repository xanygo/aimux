//  Copyright(C) 2025 github.com/hidu  All Rights Reserved.
//  Author: hidu <duv123+git@gmail.com>
//  Date: 2025-09-19

package web

import (
	"context"
	"net/http"

	"github.com/xanygo/anygo/ds/xmime"
	"github.com/xanygo/anygo/xattr"
	"github.com/xanygo/anygo/xlog"
)

func Run() {
	routerRegister()
	xmime.Register()

	listen := xattr.AppMain().GetListen("main")
	xlog.Info(context.Background(), "server Listen", xlog.String("Listen", listen))
	ser := &http.Server{
		Handler: router,
		Addr:    listen,
	}
	ser.SetKeepAlivesEnabled(true)
	err := ser.ListenAndServe()
	xlog.Warn(context.Background(), "server exit", xlog.ErrorAttr("error", err))
}
