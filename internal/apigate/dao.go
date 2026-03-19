//  Copyright(C) 2026 github.com/hidu  All Rights Reserved.
//  Author: hidu <duv123+git@gmail.com>
//  Date: 2026-03-17

package apigate

import (
	"context"
	"slices"

	"github.com/xanygo/anygo/ds/xmap"
	"github.com/xanygo/anygo/store/xkv"
)

type Dao struct {
	KVDB xkv.Hash[*Service]
}

func (d *Dao) Set(ctx context.Context, key string, value *Service) error {
	return d.KVDB.HSet(ctx, key, value)
}

func (d *Dao) Get(ctx context.Context, key string) (*Service, error) {
	v, _, err := d.KVDB.HGet(ctx, key)
	return v, err
}

func (d *Dao) GetAll(ctx context.Context) (Services, error) {
	vs, err := d.KVDB.HGetAll(ctx)
	if err != nil {
		return nil, err
	}
	return xmap.Values(vs), nil
}

func (d *Dao) GetAllActive(ctx context.Context) (Services, error) {
	ss, err := d.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	ss = slices.DeleteFunc(ss, func(s *Service) bool {
		return s.Disabled
	})
	return ss, nil
}
