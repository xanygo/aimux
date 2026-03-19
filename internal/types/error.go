//  Copyright(C) 2026 github.com/hidu  All Rights Reserved.
//  Author: hidu <duv123+git@gmail.com>
//  Date: 2026-03-13

package types

import (
	"net/http"

	"github.com/xanygo/anygo/xhttp"
)

func NewRequestError(err error) OpenAIError {
	return OpenAIError{
		Error: OpenAIErrorBody{
			Message: err.Error(),
			Type:    "invalid_request_error",
		},
	}
}

type OpenAIError struct {
	Error OpenAIErrorBody `json:"error"`
}

func (oe OpenAIError) Write(w http.ResponseWriter) {
	xhttp.WriteJSONStatus(w, http.StatusBadRequest, oe)
}

type OpenAIErrorBody struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   string `json:"param,omitempty"`
	Code    string `json:"code,omitempty"`
}
