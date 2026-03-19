//  Copyright(C) 2026 github.com/hidu  All Rights Reserved.
//  Author: hidu <duv123+git@gmail.com>
//  Date: 2026-03-18

package apigate

import "testing"

func TestValidatePath(t *testing.T) {
	type args struct {
		route string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "case 1",
			args: args{
				route: "/api",
			},
			want: true,
		},
		{
			name: "case 2",
			args: args{
				route: "/v1/api",
			},
			want: true,
		},
		{
			name: "case 3",
			args: args{
				route: "v1/api",
			},
			want: false,
		},
		{
			name: "case 4",
			args: args{
				route: "/v1//api",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidatePath(tt.args.route); got != tt.want {
				t.Errorf("ValidatePath()  want= %v", tt.want)
			}
		})
	}
}
