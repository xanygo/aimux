//  Copyright(C) 2026 github.com/hidu  All Rights Reserved.
//  Author: hidu <duv123+git@gmail.com>
//  Date: 2026-03-16

package admin

import (
	"net/http"

	"github.com/xanygo/anygo/xhttp"
)

var _ xhttp.GroupHandler = (*userHandler)(nil)

type userHandler struct {
}

func (u userHandler) GroupHandler() map[string]xhttp.PatternHandler {
	return nil
}

func (u userHandler) Index(w http.ResponseWriter, r *http.Request) {}
