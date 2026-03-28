//  Copyright(C) 2026 github.com/hidu  All Rights Reserved.
//  Author: hidu <duv123+git@gmail.com>
//  Date: 2026-03-16

package admin

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/xanygo/anygo"
	"github.com/xanygo/anygo/ds/xslice"
	"github.com/xanygo/anygo/ds/xstr"
	"github.com/xanygo/anygo/xhttp"
	"github.com/xanygo/anygo/xi18n"
	"github.com/xanygo/webr"

	"github.com/xanygo/aimux/internal/apigate"
	"github.com/xanygo/aimux/internal/factory"
)

var _ xhttp.GroupHandler = (*apiHandler)(nil)

type apiHandler struct {
}

func (a apiHandler) GroupHandler() map[string]xhttp.PatternHandler {
	return nil
}

func (a apiHandler) Index(w http.ResponseWriter, req *http.Request) {
	dao := factory.ServiceDao()
	items, err := dao.GetAll(req.Context())
	values := map[string]any{
		"Title":  anygo.Must1(xi18n.RB(req.Context(), "API 列表", "layout@menu_api")),
		"Static": apigate.Static,
		"Dyn":    items,
		"Error":  err,
	}
	dashboard.RenderWithLayout(req.Context(), w, req, "api_list.html", values)
}

func (a apiHandler) getService(req *http.Request) (*apigate.Service, bool, error) {
	query := req.URL.Query()
	typ := query.Get("type")
	id := query.Get("id")
	switch typ {
	case "static":
		srv, err := apigate.Static.FIndByID(id)
		return srv, false, err
	case "dyn":
		srv, err := factory.ServiceDao().Get(req.Context(), id)
		return srv, true, err
	default:
		err := fmt.Errorf("invalie type: %q", typ)
		return nil, false, err
	}
}

func (a apiHandler) View(w http.ResponseWriter, req *http.Request) {
	srv, _, err := a.getService(req)
	if err != nil {
		showError(w, req, err.Error())
		return
	}

	// 只保留有效的配置内容
	if srv != nil && req.URL.Query().Get("clear") != "" {
		srv = srv.CloneEnabled()
	}
	values := map[string]any{
		"Title":     anygo.Must1(xi18n.RB(req.Context(), " 查看API", "layout@view_api")) + "-" + srv.Name,
		"Srv":       srv,
		"WriteAble": false,
		"Methods":   httpMethods,
	}
	dashboard.RenderWithLayout(req.Context(), w, req, "api_edit.html", values)
}

var httpMethods = []string{"ANY", "GET", "POST", "DELETE", "PUT"}

func (a apiHandler) Edit(w http.ResponseWriter, req *http.Request) {
	srv, writeAble, err := a.getService(req)
	if err != nil {
		showError(w, req, err.Error())
		return
	}
	values := map[string]any{
		"Title":     anygo.Must1(xi18n.RB(req.Context(), "编辑 API", "layout@edit_api")) + "-" + srv.Name,
		"Srv":       srv,
		"WriteAble": writeAble,
		"Methods":   httpMethods,
	}
	dashboard.RenderWithLayout(req.Context(), w, req, "api_edit.html", values)
}

func (a apiHandler) New(w http.ResponseWriter, req *http.Request) {
	srv := &apigate.Service{
		ID:      xstr.RandIdentN(12),
		Methods: []string{"GET", "POST"},
	}
	values := map[string]any{
		"Title":     anygo.Must1(xi18n.RB(req.Context(), "创建API", "layout@create_api")),
		"Srv":       srv,
		"WriteAble": true,
		"Methods":   httpMethods,
	}
	dashboard.RenderWithLayout(req.Context(), w, req, "api_edit.html", values)
}

func (a apiHandler) Save(w http.ResponseWriter, req *http.Request) {
	var srv *apigate.Service
	err := xhttp.Bind(req, &srv)
	if err == nil {
		err = a.checkService(srv)
	}
	if err != nil {
		webr.WriteJSONAuto(w, err)
		return
	}

	for _, s := range apigate.Static {
		if s.Name == srv.Name {
			webr.WriteJSONAuto(w, fmt.Errorf("duplicate Name %q", srv.Name))
			return
		}
		if s.Route == srv.Route {
			webr.WriteJSONAuto(w, fmt.Errorf("duplicate route %q", srv.Route))
			return
		}
	}

	dao := factory.ServiceDao()
	ss, err := dao.GetAll(req.Context())
	if err != nil {
		webr.WriteJSONAuto(w, err)
		return
	}
	for _, s := range ss {
		if s.ID == srv.ID {
			continue
		}
		if s.Route == srv.Route {
			webr.WriteJSONAuto(w, fmt.Errorf("duplicate route %q", srv.Route))
			return
		}
	}
	err = dao.Set(req.Context(), srv.ID, srv)
	if err != nil {
		webr.WriteJSONAuto(w, err)
		return
	}

	rr := webr.Response{
		Jump: dashboard.AbsLink("api/Edit?type=dyn&id=" + srv.ID),
	}
	rr.WriteJSON(w)
}

func (a apiHandler) checkService(srv *apigate.Service) error {
	if srv.ID == "" {
		return errors.New("id is required")
	}
	if len(srv.Methods) == 0 {
		return errors.New("methods is empty")
	}
	if m, ok := xslice.AllContains(httpMethods, srv.Methods); !ok {
		return fmt.Errorf("invalid methods: %q", m)
	}
	srv.Route = strings.TrimSpace(srv.Route)
	if !apigate.ValidatePath(srv.Route) {
		return fmt.Errorf("invalid route: %q", srv.Route)
	}
	return nil
}

func (a apiHandler) PostReload(w http.ResponseWriter, req *http.Request) {
	dao := factory.ServiceDao()
	err := apigate.LoadFromDB(req.Context(), dao)
	webr.WriteJSONAuto(w, err)
}
