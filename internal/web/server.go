//  Copyright(C) 2025 github.com/hidu  All Rights Reserved.
//  Author: hidu <duv123+git@gmail.com>
//  Date: 2025-09-19

package web

import (
	"context"
	"net/http"

	"github.com/xanygo/anygo/ds/xmime"
	"github.com/xanygo/anygo/ds/xsync"
	"github.com/xanygo/anygo/xattr"
	"github.com/xanygo/anygo/xlog"

	"github.com/xanygo/aimux/internal/factory"
)

func Run() error {
	xmime.Register()

	factory.MustLoadFromDB()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	apiListen := xattr.AppMain().GetListen("api")
	// 若没有配置独立的 api gateway 端口，则和 admin 使用同一个端口
	if apiListen == "" {
		return runAdminServer(ctx, true)
	}

	var wg xsync.WaitGroup
	wg.GoCtxErr(ctx, func(ctx context.Context) error {
		return runAdminServer(ctx, false)
	})
	wg.GoCtxErr(ctx, func(ctx context.Context) error {
		return runAPIServer(ctx, apiListen)
	})
	return wg.Wait()
}

func runAdminServer(ctx context.Context, withAPI bool) error {
	listen := xattr.AppMain().GetListen("admin")
	xlog.Info(ctx, "admin server Listen", xlog.String("Listen", listen))
	ser := &http.Server{
		Handler: initAdminRouter(withAPI),
		Addr:    listen,
	}
	ser.SetKeepAlivesEnabled(true)
	err := ser.ListenAndServe()
	xlog.Warn(ctx, "admin server exit", xlog.ErrorAttr("error", err))
	return err
}

func runAPIServer(ctx context.Context, listen string) error {
	xlog.Info(ctx, "api server Listen", xlog.String("Listen", listen))
	ser := &http.Server{
		Handler: initAPIRouter(),
		Addr:    listen,
	}
	ser.SetKeepAlivesEnabled(true)
	err := ser.ListenAndServe()
	xlog.Warn(ctx, "api server exit", xlog.ErrorAttr("error", err))
	return err
}
