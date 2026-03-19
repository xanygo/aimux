//  Copyright(C) 2025 github.com/hidu  All Rights Reserved.
//  Author: hidu <duv123+git@gmail.com>
//  Date: 2025-09-19

package main

import (
	"flag"

	"github.com/xanygo/anygo/xattr"
	"github.com/xanygo/anygo/xcfg"
	"github.com/xanygo/ext"

	"github.com/xanygo/aimux/internal/config"
	"github.com/xanygo/aimux/internal/web"
)

func init() {
	ext.Init()
}

var c = flag.String("conf", "conf/app.yml", "app main config file")

func main() {
	flag.Parse()
	xattr.MustInitAppMain(*c, xcfg.Parse)
	config.Bootstrap()
	web.Run()
}
