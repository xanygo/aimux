//  Copyright(C) 2026 github.com/hidu  All Rights Reserved.
//  Author: hidu <duv123+git@gmail.com>
//  Date: 2026-03-13

package apigate

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"slices"
	"strings"

	"github.com/xanygo/anygo/ds/xslice"
	"github.com/xanygo/anygo/xcodec"
	"github.com/xanygo/anygo/xhttp/xhttpc"
	"github.com/xanygo/anygo/xlog"

	"github.com/xanygo/aimux/internal/types"
)

type Services []*Service

func (ns Services) Clone() Services {
	if len(ns) == 0 {
		return nil
	}
	clone := make(Services, len(ns))
	for i, n := range ns {
		clone[i] = n.Clone()
	}
	return clone
}

func (ns Services) CloneEnabled() Services {
	if len(ns) == 0 {
		return nil
	}
	clone := make(Services, 0, len(ns))
	for _, n := range ns {
		if n.Disabled {
			continue
		}
		nns := n.Clone()
		nns.Auths = xslice.Filter(nns.Auths, func(index int, item *Auth, okTotal int) bool {
			return !item.Disabled
		})
		nns.Nodes = xslice.Filter(nns.Nodes, func(index int, item *Node, okTotal int) bool {
			item.Auths = xslice.Filter(item.Auths, func(index int, item *Auth, okTotal int) bool {
				return !item.Disabled
			})
			return !item.Disabled
		})
		clone = append(clone, nns)
	}
	return clone
}

func (ns Services) FindByName(name string) (*Service, error) {
	for _, n := range ns {
		if n.Name == name {
			return n, nil
		}
	}
	return nil, fmt.Errorf("service Name=%q not found", name)
}

func (ns Services) FIndByID(id string) (*Service, error) {
	for _, n := range ns {
		if n.ID == id {
			return n, nil
		}
	}
	return nil, fmt.Errorf("service ID=%q not found", id)
}

func (ns Services) Routes() []string {
	return xslice.FilterAs(ns, func(index int, item *Service, ok int) (string, bool) {
		return item.Route, true
	})
}

func (ns Services) CheckRoutes(other Services) error {
	if len(ns) == 0 || len(other) == 0 {
		return nil
	}
	nsr := xslice.ToMap(ns.Routes(), true)
	for _, name := range other.Routes() {
		if nsr[name] {
			return fmt.Errorf("service Route %q already exists", name)
		}
	}
	return nil
}

// Service 表示一组下游服务的集合
type Service struct {
	ID       string   `json:"id" yaml:"id"`
	Name     string   `json:"name" yaml:"name"`     // 组名称，必填
	Remark   string   `json:"remark" yaml:"remark"` // 备注信息
	Methods  []string `json:"methods" yaml:"methods"`
	Route    string   `json:"route" yaml:"route"`       // 路由地址，例如 "/api/v1/group"
	Nodes    Nodes    `json:"nodes" yaml:"nodes"`       // 下游服务列表
	Auths    Auths    `json:"auths" yaml:"auths"`       // API 秘钥列表
	Disabled bool     `json:"disabled" yaml:"disabled"` // 是否禁用
	Model    string   `json:"model" yaml:"model"`       // 传入的模型名称，可选
}

func (s *Service) Clone() *Service {
	return &Service{
		Name:     s.Name,
		Remark:   s.Remark,
		Methods:  slices.Clone(s.Methods),
		Route:    s.Route,
		Nodes:    s.Nodes.Clone(),
		Auths:    s.Auths.Clone(),
		Disabled: s.Disabled,
		Model:    s.Model,
	}
}

var _ http.Handler = (*Service)(nil)

// RouterPattern 转换为 xhttp.Router 所用的 Pattern
func (s *Service) RouterPattern() string {
	var b strings.Builder
	if len(s.Methods) != 0 {
		b.WriteString(strings.Join(s.Methods, ","))
		b.WriteString(" ")
	}
	b.WriteString(s.Route)
	return b.String()
}

func (s *Service) getAuthToken(req *http.Request) (string, error) {
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header missing")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 {
		return "", errors.New("invalid authorization header format")
	}
	if parts[0] != "Bearer" {
		return "", errors.New("invalid authorization type")
	}
	return parts[1], nil
}

func (s *Service) checkModel(model string) error {
	if s.Model != "" && s.Model != model {
		return fmt.Errorf("model does not match, got=%q", model)
	}
	return nil
}

// ServeHTTP 处理 HTTP 请求，Service 以及相关信息都已经剔除掉 Disabled 状态的数据
// 并且注册到 HTTP Router 里的是一份拷贝的数据，所以不需要锁
func (s *Service) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	xlog.Info(req.Context(), "Service.ServeHTTP")

	err := s.checkToken(req)
	if err != nil {
		ret := types.NewRequestError(err)
		ret.Write(w)
		return
	}
	node, err := s.oneNodeByWeight()
	if err != nil {
		ret := types.NewRequestError(err)
		ret.Write(w)
		return
	}
	node.serveHTTP(w, req, s)
}

var errNoAuths = errors.New("system error: miss authorization config")

func (s *Service) checkToken(req *http.Request) error {
	if len(s.Auths) == 0 {
		return errNoAuths
	}
	token, err := s.getAuthToken(req)
	if err != nil {
		return err
	}
	for _, node := range s.Auths {
		if node.APIKey == token {
			return nil
		}
	}
	return errors.New("invalid token")
}

func (s *Service) Enable() bool {
	return !s.Disabled
}

func (s *Service) oneNodeByWeight() (*Node, error) {
	if len(s.Nodes) == 0 {
		return nil, errors.New("no enabled service nodes available")
	}
	var totalWeight int

	for i := range s.Nodes {
		node := s.Nodes[i]
		totalWeight += int(node.Weight)
	}

	if totalWeight > 0 {
		r := rand.Intn(totalWeight)
		for _, node := range s.Nodes {
			r -= int(node.Weight)
			if r < 0 {
				return node, nil
			}
		}
	}
	r := rand.Intn(len(s.Nodes))
	return s.Nodes[r], nil
}

type Nodes []*Node

func (ns Nodes) Clone() Nodes {
	if len(ns) == 0 {
		return nil
	}
	clone := make(Nodes, len(ns))
	for i, n := range ns {
		clone[i] = n.Clone()
	}
	return clone
}

// Node 表示单个下游服务节点
type Node struct {
	Name     string `json:"name" yaml:"name"`         // 服务名称
	Endpoint string `json:"endpoint" yaml:"endpoint"` // 服务访问地址
	Remark   string `json:"remark" yaml:"remark"`     // 备注信息
	Weight   int16  `json:"weight" yaml:"weight"`     // 访问权重
	Disabled bool   `json:"disabled" yaml:"disabled"` // 是否禁用
	Auths    Auths  `json:"auths" yaml:"auths"`       // API 秘钥列表
	Models   Models `json:"models" yaml:"models"`     // 可选模型列表
}

func (node *Node) Clone() *Node {
	return &Node{
		Name:     node.Name,
		Endpoint: node.Endpoint,
		Remark:   node.Remark,
		Weight:   node.Weight,
		Disabled: node.Disabled,
		Auths:    node.Auths.Clone(),
		Models:   node.Models.Clone(),
	}
}

var client = &xhttpc.Client{}

func (node *Node) serveHTTP(w http.ResponseWriter, req *http.Request, s *Service) {
	cred, err := node.oneCredentialByWeight()

	var mod *Model
	if err == nil {
		mod, err = node.oneModelByWeight()
	}
	var body []byte
	if err == nil {
		body, err = io.ReadAll(req.Body)
	}

	if err != nil {
		ret := types.NewRequestError(err)
		ret.Write(w)
		return
	}

	if s.Model != "" || mod != nil {
		data := make(map[string]any)
		if err == nil {
			err = xcodec.JSON.Decode(body, &data)
		}
		inputModel, ok := data["model"].(string)
		if err == nil && !ok {
			err = errors.New("no model field")
		}
		if err == nil {
			if s.Model != "" && s.Model != inputModel {
				err = fmt.Errorf("invalid model %q", inputModel)
			}
			if err == nil && mod != nil {
				data["model"] = mod.ID
				body, err = xcodec.JSON.Encode(data)
			}
		}
	}

	if err != nil {
		ret := types.NewRequestError(err)
		ret.Write(w)
		return
	}
	nr, err := http.NewRequestWithContext(req.Context(), req.Method, node.Endpoint, nil)
	if err != nil {
		ret := types.NewRequestError(err)
		ret.Write(w)
		return
	}

	if cred != nil {
		if cred.APIKey != "" {
			nr.Header.Set("Authorization", "Bearer "+cred.APIKey)
		}
	}

	if value := req.Header.Get("User-Agent"); value != "" {
		nr.Header.Set("User-Agent", value)
	}

	if err != nil {
		ret := types.NewRequestError(err)
		ret.Write(w)
		return
	}
	nr.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader(body)), nil
	}

	resp, err := client.RoundTrip(nr)
	if err != nil {
		ret := types.NewRequestError(err)
		ret.Write(w)
		return
	}
	defer resp.Body.Close()
	if value := resp.Header.Get("Content-Type"); value != "" {
		nr.Header.Set("Content-Type", value)
	}
	if value := resp.Header.Get("Content-Encoding"); value != "" {
		nr.Header.Set("Content-Encoding", value)
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (node *Node) oneCredentialByWeight() (*Auth, error) {
	if len(node.Auths) == 0 {
		return nil, nil
	}
	var totalWeight int
	for _, cred := range node.Auths {
		totalWeight += int(cred.Weight)
	}

	if totalWeight > 0 {
		r := rand.Intn(totalWeight)
		for _, cred := range node.Auths {
			r -= int(cred.Weight)
			if r < 0 {
				return cred, nil
			}
		}
	}
	r := rand.Intn(len(node.Auths))
	return node.Auths[r], nil
}

func (node *Node) oneModelByWeight() (*Model, error) {
	if len(node.Models) == 0 {
		return nil, nil // 允许为空
	}
	var totalWeight int

	for _, mod := range node.Models {
		totalWeight += int(mod.Weight)
	}

	if totalWeight > 0 {
		r := rand.Intn(totalWeight)
		for _, mod := range node.Models {
			r -= int(mod.Weight)
			if r < 0 {
				return mod, nil
			}
		}
	}
	r := rand.Intn(len(node.Models))
	return node.Models[r], nil
}

type Models []*Model

func (ms Models) Clone() Models {
	if len(ms) == 0 {
		return nil
	}
	clone := make(Models, len(ms))
	for i, m := range ms {
		clone[i] = m.Clone()
	}
	return clone
}

type Model struct {
	ID       string `json:"id" yaml:"id"`         // 模型的id，如 qwen-image-2.0-pro
	Remark   string `json:"remark" yaml:"remark"` // 备注信息
	Disabled bool   `json:"disabled" yaml:"disabled"`
	Weight   int64  `json:"weight" yaml:"weight"` // 权重，可用于轮询或优先级
}

func (m *Model) Clone() *Model {
	return &Model{
		ID:       m.ID,
		Remark:   m.Remark,
		Disabled: m.Disabled,
		Weight:   m.Weight,
	}
}

type Auths []*Auth

func (mc Auths) Clone() Auths {
	if len(mc) == 0 {
		return nil
	}
	clone := make(Auths, len(mc))
	for i, m := range mc {
		clone[i] = m.Clone()
	}
	return clone
}

// Auth 表示服务访问所需的认证信息
type Auth struct {
	APIKey   string `json:"apikey" yaml:"apikey"`     // API Key 值
	Remark   string `json:"remark" yaml:"remark"`     // 备注信息
	Disabled bool   `json:"disabled" yaml:"disabled"` // 是否禁用
	Weight   int64  `json:"weight" yaml:"weight"`     // 权重，可用于轮询或优先级
}

func (c *Auth) Clone() *Auth {
	return &Auth{
		APIKey:   c.APIKey,
		Remark:   c.Remark,
		Disabled: c.Disabled,
		Weight:   c.Weight,
	}
}
