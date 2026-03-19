//  Copyright(C) 2025 github.com/hidu  All Rights Reserved.
//  Author: hidu <duv123+git@gmail.com>
//  Date: 2025-09-19

package admin

import (
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, dashboard.PathPrefix+"/api/", http.StatusFound)
	// values := map[string]any{
	//	"User":     nvwa.UserFormContext(r.Context()),
	//	"SubTitle": "首页",
	// }
	// dashboard.RenderWithLayout(r.Context(), w, r, "index.html", values)
}
