//  Copyright(C) 2026 github.com/hidu  All Rights Reserved.
//  Author: hidu <duv123+git@gmail.com>
//  Date: 2026-04-18

package admin

import "github.com/xanygo/anygo/xhttp"

type formAction struct {
	// FormAction 表单点击 submit 后提交的行为
	FormAction string
}

const actionSaveReload = "save_reload"

func isSaveAndReload(bd *xhttp.Binder) bool {
	var action formAction
	return bd.BinJSON(&action) == nil && action.FormAction == actionSaveReload
}
