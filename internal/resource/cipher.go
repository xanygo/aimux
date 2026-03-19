//  Copyright(C) 2025 github.com/hidu  All Rights Reserved.
//  Author: hidu <duv123+git@gmail.com>
//  Date: 2025-10-29

package resource

import (
	"github.com/xanygo/anygo/ds/xsync"
	"github.com/xanygo/anygo/xcodec"
	"github.com/xanygo/anygo/xcodec/xbase"

	"github.com/xanygo/aimux/internal/config"
)

var cipher = xsync.OnceInit[xcodec.Cipher]{
	New: func() xcodec.Cipher {
		return &xcodec.AesOFB{
			Key: config.SecretKey(),
		}
	},
}

func Encrypt(text string) (string, error) {
	bf, err := cipher.Load().Encrypt([]byte(text))
	if err != nil {
		return "", err
	}
	return xbase.Base58.EncodeToString(bf), nil
}

func Decrypt(text string) (string, error) {
	bf, err := xbase.Base58.DecodeString(text)
	if err != nil {
		return "", err
	}
	val, err := cipher.Load().Decrypt(bf)
	if err != nil {
		return "", err
	}
	return string(val), nil
}
