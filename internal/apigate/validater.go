//  Copyright(C) 2026 github.com/hidu  All Rights Reserved.
//  Author: hidu <duv123+git@gmail.com>
//  Date: 2026-03-18

package apigate

import (
	"regexp"
	"strings"
)

var pathRegexp = regexp.MustCompile(`^/[a-zA-Z0-9/_\-\.]*$`)

func ValidatePath(route string) bool {
	if !pathRegexp.MatchString(route) || strings.Contains(route, "//") || strings.Contains(route, "..") {
		return false
	}
	return true
}
