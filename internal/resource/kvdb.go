//  Copyright(C) 2026 github.com/hidu  All Rights Reserved.
//  Author: hidu <duv123+git@gmail.com>
//  Date: 2026-03-13

package resource

import (
	"path/filepath"
	"strings"

	"github.com/xanygo/anygo/ds/xsync"
	"github.com/xanygo/anygo/store/xkv"
	"github.com/xanygo/anygo/store/xkv/xkvx"
	"github.com/xanygo/anygo/store/xredis"
	"github.com/xanygo/anygo/xattr"
	"github.com/xanygo/anygo/xcodec"

	"github.com/xanygo/aimux/internal/config"
)

var kvDB = xsync.OnceInit[xkv.Storage[string]]{
	New: func() xkv.StringStorage {
		dbType := xattr.GetDefault[string]("KVDB", "file")
		pre, suf, found := strings.Cut(dbType, ":")
		switch pre {
		case "memory":
			return xkv.NewMemoryStore()
		case "file":
			return &xkv.FileStore{
				DataDir: filepath.Join(xattr.DataDir(), "xkv_db"),
			}
		case "redis":
			if !found || suf == "" {
				panic("invalid redis db type:" + dbType)
			}
			return &xkvx.RedisStorage{
				KeyPrefix: "kxcms|",
				Client:    xredis.NewClient(suf),
			}
		default:
			panic("not support KVDB type: " + dbType)
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

func kvdb() xkv.StringStorage {
	return kvDB.Load()
}

func HashDB[T any](key string) xkv.Hash[T] {
	tr := &xkv.Transformer[T]{
		Storage: kvdb(),
		Codec:   coder.Load(),
	}
	return tr.Hash(key)
}
