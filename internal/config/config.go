//  Copyright(C) 2025 github.com/hidu  All Rights Reserved.
//  Author: hidu <duv123+git@gmail.com>
//  Date: 2025-09-22

package config

import (
	"github.com/xanygo/anygo/xattr"
)

func AdminPath() string {
	return xattr.GetDefault[string]("AdminPath", "/admin/")
}

func SecretKey() string {
	return xattr.GetDefault[string]("SecretKey", "e60734fa48b72bb766aa639060352efd")
}
