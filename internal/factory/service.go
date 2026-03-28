//  Copyright(C) 2026 github.com/hidu  All Rights Reserved.
//  Author: hidu <duv123+git@gmail.com>
//  Date: 2026-03-17

package factory

import (
	"context"
	"time"

	"github.com/xanygo/anygo"
	"github.com/xanygo/anygo/ds/xsync"
	"github.com/xanygo/anygo/xlog"

	"github.com/xanygo/aimux/internal/apigate"
	"github.com/xanygo/aimux/internal/resource"
)

var srvDao = &xsync.OnceInit[*apigate.Dao]{
	New: func() *apigate.Dao {
		return &apigate.Dao{
			KVDB: resource.HashDB[*apigate.Service]("services"),
		}
	},
}

func ServiceDao() *apigate.Dao {
	return srvDao.Load()
}

func MustLoadFromDB() {
	var err error
	var ss apigate.Services
	for i := 0; i < 3; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		ss, err = ServiceDao().GetAllActive(ctx)
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
