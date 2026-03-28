//  Copyright(C) 2025 github.com/hidu  All Rights Reserved.
//  Author: hidu <duv123+git@gmail.com>
//  Date: 2025-09-19

package config

import (
	"sync"

	"github.com/xanygo/anygo/xcfg"
)

type User struct {
	Username string `yaml:"Username"`
	Password string `yaml:"Password"`
}

var allUsers []User

var adminUserOnce sync.Once

func loadAdminUser() {
	type appConfig struct {
		Users []User `yaml:"Users"`
	}
	cfg := &appConfig{}
	xcfg.MustParse("app", cfg)
	allUsers = cfg.Users
}

func FindUser(name string) *User {
	adminUserOnce.Do(loadAdminUser)
	for _, user := range allUsers {
		if user.Username == name {
			return &user
		}
	}
	return nil
}

func Users() (users []User) {
	adminUserOnce.Do(loadAdminUser)
	return allUsers
}
