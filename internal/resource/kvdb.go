//  Copyright(C) 2026 github.com/hidu  All Rights Reserved.
//  Author: hidu <duv123+git@gmail.com>
//  Date: 2026-03-13

package resource

import (
	"context"
	"sync/atomic"

	"github.com/xanygo/anygo/xkv"
	"github.com/xanygo/anygo/xkv/xkvx"
)

func getKVDB[V any]() xkv.Storage[V] {
	return &xkv.Monitor[V]{
		Store: xkvx.MustLoad[V]("default"),
		After: func(ctx context.Context, dataType xkv.DataType, action string, err error, keys ...string) {
			if !xkv.IsReadAction(dataType, action) {
				needReload.Store(true)
			}
		},
	}
}

func HashDB[T any](key string) xkv.Hash[T] {
	tr := getKVDB[T]()
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
