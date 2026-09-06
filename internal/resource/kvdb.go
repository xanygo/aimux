//  Copyright(C) 2026 github.com/hidu  All Rights Reserved.
//  Author: hidu <duv123+git@gmail.com>
//  Date: 2026-03-13

package resource

import (
	"context"
	"log"
	"sync/atomic"

	"github.com/xanygo/anygo/ds/xsync"
	"github.com/xanygo/anygo/store/xkv"
	"github.com/xanygo/anygo/store/xkv/xkvx"
	"github.com/xanygo/anygo/xcodec"

	"github.com/xanygo/aimux/internal/config"
)

var kvDB = xsync.OnceInit[xkv.Storage[string]]{
	New: func() xkv.StringStorage {
		return &xkv.Monitor[string]{
			Store: xkvx.MustLoad[string]("default"),
			After: func(ctx context.Context, dataType xkv.DataType, action string, err error, keys ...string) {
				if !xkv.IsReadAction(dataType, action) {
					log.Println("kxdb action:", dataType, action)
					needReload.Store(true)
				}
			},
		}
	},
}

var coder = &xsync.OnceInit[xcodec.Codec]{
	New: func() xcodec.Codec {
		aes := &xcodec.AesOFB{
			Key: config.SecretKey(),
		}
		return xcodec.CodecWithCipher(xcodec.JSON, aes)
	},
}

func HashDB[T any](key string) xkv.Hash[T] {
	tr := &xkv.Transformer[T]{
		Storage: kvDB.Load(),
		Codec:   coder.Load(),
	}
	return tr.Hash(key)
}

var needReload atomic.Bool

// NeedReload 数据已变化，是否需要重新加载
func NeedReload() bool {
	return needReload.Load()
}

func ResetNeedReload() {
	needReload.Store(false)
}
