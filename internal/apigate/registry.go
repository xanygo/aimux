//  Copyright(C) 2026 github.com/hidu  All Rights Reserved.
//  Author: hidu <duv123+git@gmail.com>
//  Date: 2026-03-13

package apigate

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/xanygo/anygo/xhttp"
	"github.com/xanygo/anygo/xlog"
)

var _ http.Handler = (*Registry)(nil)

type Registry struct {
	groups map[string]Services
	mux    sync.RWMutex
	router *xhttp.Router
}

func (r *Registry) MustRegisterStatic(ss Services) {
	err := r.register("static", ss)
	if err == nil {
		return
	}
	panic(fmt.Sprintf("register static service failed: %v", err))
}

func (r *Registry) RegisterDny(ss Services) error {
	return r.register("dny", ss)
}

func (r *Registry) register(group string, ss Services) error {
	ss = ss.CloneEnabled()
	r.mux.Lock()
	defer r.mux.Unlock()
	for g, gs := range r.groups {
		if g == group {
			continue
		}
		if err := gs.CheckRoutes(ss); err != nil {
			return err
		}
	}
	r.groups[group] = ss

	return r.rebuildRouter()
}

func (r *Registry) rebuildRouter() error {
	router := xhttp.NewRouter()
	for group, gs := range r.groups {
		for _, srv := range gs {
			if err := router.Handle(srv.RouterPattern(), srv); err != nil {
				return fmt.Errorf("register group %q : %w", group, err)
			}
		}
	}

	r.router = router
	return nil
}

func (r *Registry) FindByRoute(route string) *Service {
	r.mux.RLock()
	defer r.mux.RUnlock()
	for _, gs := range r.groups {
		for _, srv := range gs {
			if srv.Route == route {
				return srv
			}
		}
	}
	return nil
}

func (r *Registry) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	xlog.Info(req.Context(), "Registry.ServeHTTP", xlog.String("req", req.RequestURI))
	r.mux.RLock()
	defer r.mux.RUnlock()
	r.router.ServeHTTP(w, req)
}

var defaultRegistry = &Registry{
	groups: make(map[string]Services),
	router: xhttp.NewRouter(),
}

func Default() *Registry {
	return defaultRegistry
}
