package api

import (
	"testing"
)

func Test_censured(t *testing.T) {
	type args struct {
		str  string
		subs []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "one",
			args: args{
				str:  "qwerty йцукен zxvbnm",
				subs: []string{"qwerty", "йцукен", "zxvbnm"},
			},
			want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := censured(tt.args.str, tt.args.subs...); got != tt.want {
				t.Errorf("censured() = %v, want %v", got, tt.want)
			}
		})
	}
}
