//  Copyright(C) 2025 github.com/hidu  All Rights Reserved.
//  Author: hidu <duv123+git@gmail.com>
//  Date: 2025-09-22

package config

import (
	"github.com/xanygo/anygo/xattr"
	"github.com/xanygo/anygo/xhttp"
)

func AdminPath() string {
	return xattr.GetDefault[string]("AdminPath", "/admin")
}

func AdminLink(s string) string {
	return xhttp.PathJoin(AdminPath(), s)
}
