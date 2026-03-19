//  Copyright(C) 2026 github.com/hidu  All Rights Reserved.
//  Author: hidu <duv123+git@gmail.com>
//  Date: 2026-03-17

package factory

import (
	"github.com/xanygo/anygo/ds/xsync"

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
