//  Copyright(C) 2026 github.com/hidu  All Rights Reserved.
//  Author: hidu <duv123+git@gmail.com>
//  Date: 2026-03-28

package config

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"

	"github.com/xanygo/anygo"
	"github.com/xanygo/anygo/xattr"
	"github.com/xanygo/anygo/xcfg"
	"github.com/xanygo/anygo/xio/xfs"
	"github.com/xanygo/anygo/xlog"

	"github.com/xanygo/aimux/internal/apigate"
)

// 加载配置文件中的静态配置
func loadStaticAPIServices() {
	value, ok := xattr.AppMain().GetOther("StaticAPIFile")
	xlog.Info(context.Background(), "read AppMain().StaticAPIFile", xlog.Any("value", value), xlog.Bool("ok", ok))
	if !ok {
		return
	}

	filename, ok := value.(string)
	if !ok {
		log.Fatalf("invalid StaticAPIFile=%#v", value)
	}

	var ss apigate.Services
	err := xcfg.Parse(filename, &ss)
	xlog.Info(context.Background(), "parser static_api", xlog.ErrorAttr("err", err))
	if err != nil && errors.Is(err, fs.ErrNotExist) {
		xlog.Info(context.Background(), "static_api file not found, skipped")
	}
	anygo.Must(err)

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
