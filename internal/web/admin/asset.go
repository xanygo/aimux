//  Copyright(C) 2025 github.com/hidu  All Rights Reserved.
//  Author: hidu <duv123+git@gmail.com>
//  Date: 2025-09-19

//go:build !release

//go:generate anygo-encrypt-zip -token 6b65f8daae3270d839b443cd8327a801 -o asset.ez -go asset_ez.go -var files tpls

package admin

import "embed"

//go:embed tpls/*
var files embed.FS
