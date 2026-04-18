//  Copyright(C) 2026 github.com/hidu  All Rights Reserved.
//  Author: hidu <duv123+git@gmail.com>
//  Date: 2026-03-14

package admin

import (
	"context"
	"errors"
	"io/fs"
	"net/http"
	"path"
	"text/template"

	"github.com/xanygo/anygo"
	"github.com/xanygo/nvwa"

	"github.com/xanygo/aimux/internal/config"
	"github.com/xanygo/aimux/internal/metric"
	"github.com/xanygo/aimux/internal/resource"
	"github.com/xanygo/aimux/internal/resource/i18n"
)

var errUserNotFound = errors.New("user not found")

var dashboard = &nvwa.Dashboard{
	PathPrefix:    config.AdminPath(),
	AssetPrefix:   path.Join(config.AdminPath(), "asset"),
	RegisterAsset: "/asset/",
	SecretKey:     config.SecretKey(),
	TemplateFS:    anygo.Must1(fs.Sub(files, "tpls")),
	Bundle:        i18n.Resource,
	UserFinder: func(ctx context.Context, name string) (*nvwa.User, error) {
		u := config.FindUser(name)
		if u == nil {
			return nil, errUserNotFound
		}
		return &nvwa.User{
			Username: u.Username,
			AuthCode: u.Password,
		}, nil
	},
	UserCheckLogin: func(ctx context.Context, req *http.Request, name string, psw string) (*nvwa.User, error) {
		if !metric.CanLogin() {
			return nil, errors.New("denied, please wait a minute")
		}
		u := config.FindUser(name)
		if u == nil {
			metric.LoginFailed()
			return nil, errUserNotFound
		}
		if u.Username == name && u.Password == psw {
			return &nvwa.User{
				Username: u.Username,
				AuthCode: u.Password,
			}, nil
		}
		metric.LoginFailed()
		return nil, errUserNotFound
	},
	FuncMap: template.FuncMap{
		"zNeedReload": func() bool {
			return resource.NeedReload()
		},
	},
}

func showError(w http.ResponseWriter, req *http.Request, err string) {
	values := map[string]any{
		"error": err,
	}
	dashboard.RenderWithLayout(req.Context(), w, req, "error.html", values)
}
