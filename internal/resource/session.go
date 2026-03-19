//  Copyright(C) 2025 github.com/hidu  All Rights Reserved.
//  Author: hidu <duv123+git@gmail.com>
//  Date: 2025-09-25

package resource

import (
	"net/http"
	"path/filepath"
	"time"

	"github.com/xanygo/anygo/ds/xsync"
	"github.com/xanygo/anygo/store/xsession"
	"github.com/xanygo/anygo/xattr"
	"github.com/xanygo/anygo/xcodec"

	"github.com/xanygo/aimux/internal/config"
)

type SessionNewStorageFunc func(writer http.ResponseWriter, request *http.Request) xsession.Storage

var sessionStore = xsync.OnceInit[SessionNewStorageFunc]{
	New: func() SessionNewStorageFunc {
		sessionType := xattr.GetDefault[string]("Session", "cookie")
		ttl := 7 * 24 * time.Hour
		switch sessionType {
		case "kvdb":
			db := &xsession.KVStore{
				DB:  kvdb(),
				TTL: ttl,
			}
			return func(writer http.ResponseWriter, request *http.Request) xsession.Storage {
				return db
			}
		case "cookie":
			return func(writer http.ResponseWriter, request *http.Request) xsession.Storage {
				return &xsession.CookieStore{
					Writer:  writer,
					Request: request,
					Cipher: &xcodec.AesOFB{
						Key: config.SecretKey(),
					},
				}
			}
		case "memory":
			db := xsession.NewMemoryStore(1024, ttl)
			return func(writer http.ResponseWriter, request *http.Request) xsession.Storage {
				return db
			}
		default:
			db := xsession.NewFileStore(filepath.Join(xattr.TempDir(), "session"), ttl)
			return func(writer http.ResponseWriter, request *http.Request) xsession.Storage {
				return db
			}
		}
	},
}

func SessionStorage() SessionNewStorageFunc {
	return sessionStore.Load()
}
