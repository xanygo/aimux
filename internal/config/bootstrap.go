//  Copyright(C) 2026 github.com/hidu  All Rights Reserved.
//  Author: hidu <duv123+git@gmail.com>
//  Date: 2026-03-13

package config

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/xanygo/anygo"
	"github.com/xanygo/anygo/xattr"
	"github.com/xanygo/anygo/xcodec"
	"github.com/xanygo/anygo/xio/xfs"
	"github.com/xanygo/anygo/xlog"

	"github.com/xanygo/aimux/internal/apigate"
)

func Bootstrap() {
	initLogger()
	parserServices()
	initRPCDump()
}

func parserServices() {
	services, ok := xattr.AppMain().GetOther("Services")
	xlog.Info(context.Background(), "read AppMain().Services", xlog.Bool("exists", ok))
	if ok {
		var ss apigate.Services
		anygo.Must(xcodec.Convert(services, &ss))
		xlog.Info(context.Background(), "AppMain().Services", xlog.Int("len", len(ss)))
		var num = 1
		for index := range ss {
			for {
				id := fmt.Sprintf("static_%d", num)
				if s, _ := ss.FIndByID(id); s == nil {
					ss[index].ID = id
					break
				}
				num++
			}
		}
		apigate.Static = ss
		apigate.Default().MustRegisterStatic(ss)
	}
}

func initRPCDump() {
	doDump := xattr.GetDefault[bool]("RPCDump", false)
	if !doDump {
		return
	}
	w := &xfs.Rotator{
		Path: filepath.Join(xattr.LogDir(), "rpcdump", "dump.txt"),
	}
	anygo.Must(w.Init())
	apigate.SetDumpWriter(w)
}

func initLogger() {
	logLevelStr := xattr.GetDefault[string]("LogLevel", "INFO")

	logLevel := xlog.ParserLevel(logLevelStr)
	xlog.DefaultLevel = logLevel
	xlog.InitAllDefaultLogger()
}
