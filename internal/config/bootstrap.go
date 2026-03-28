//  Copyright(C) 2026 github.com/hidu  All Rights Reserved.
//  Author: hidu <duv123+git@gmail.com>
//  Date: 2026-03-13

package config

import (
	"context"
	"path/filepath"

	"github.com/xanygo/anygo"
	"github.com/xanygo/anygo/xattr"
	"github.com/xanygo/anygo/xlog"
	"github.com/xanygo/anygo/xnet/xrpc"
	"github.com/xanygo/anygo/xnet/xservice"
)

func Bootstrap() {
	initFramework()

	{
		initRPCDump()
		loadStaticAPIServices()
	}
}

// 依据配置，初始化框架
func initFramework() {
	// 可选，初始化日志配置
	{
		logLevelStr := xattr.GetDefault[string]("LogLevel", "INFO")
		logLevel := xlog.ParserLevel(logLevelStr)
		xlog.DefaultLevel = logLevel
		xlog.InitAllDefaultLogger()
	}

	// 可选：加载 service 配置
	{
		err := xservice.LoadDir(context.Background(), filepath.Join(xattr.ConfDir(), "service", "*.yml"))
		anygo.Must(err)
	}

	// 可选：给 RPC client 注册日志中间件，用于打印日志
	{
		rl := &xrpc.Logger{}
		xrpc.RegisterTCPIT(rl.Interceptor())
	}
}
