//  Copyright(C) 2026 github.com/hidu  All Rights Reserved.
//  Author: hidu <duv123+git@gmail.com>
//  Date: 2026-03-18

package metric

import (
	"time"

	"github.com/xanygo/anygo/ds/xcounter"
)

var loginFailCounter = xcounter.NewSlidingWindow(30*time.Minute, 10*time.Second)

func CanLogin() bool {
	return loginFailCounter.WindowTotal() < 100
}

func LoginFailed() {
	loginFailCounter.Incr()
}

func CounterStatus() any {
	return map[string]any{
		"LoginFailed(登录失败)": map[string]any{
			"All":   loginFailCounter.LifetimeTotal(),
			"30min": loginFailCounter.WindowTotal(),
			"1min":  loginFailCounter.Count(time.Minute),
			"5min":  loginFailCounter.Count(5 * time.Minute),
		},
	}
}
