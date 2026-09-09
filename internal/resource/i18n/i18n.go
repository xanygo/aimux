//  Copyright(C) 2024 github.com/hidu  All Rights Reserved.
//  Author: hidu <duv123+git@gmail.com>
//  Date: 2024-09-08

package i18n

import (
	"embed"

	"github.com/xanygo/anygo"
	"github.com/xanygo/anygo/xenc"
	"github.com/xanygo/anygo/xi18n"
	"gopkg.in/yaml.v3"
)

//go:embed zh/* en/*
var files embed.FS

var Resource = xi18n.NewBundle(xi18n.LangZh, xi18n.LangEn)

func init() {
	err := xi18n.LoadFS(Resource, files, ".", ".yml", xenc.UnmarshalFunc(yaml.Unmarshal))
	anygo.Must(err)
}
