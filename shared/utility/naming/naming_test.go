package naming

import (
	"testing"
)

func TestHumpNaming(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{name: "aa_bb_cc"},
			want: "AaBbCc",
		},
		{
			name: "2",
			args: args{name: "aa_Bb_cc"},
			want: "AaBbCc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HumpNaming(tt.args.name); got != tt.want {
				t.Errorf("HumpNaming() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnderlineNaming(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{name: "AaBbCc"},
			want: "aa_bb_cc",
		},
		{
			name: "2",
			args: args{name: "AaBBCc"},
			want: "aa_bb_cc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UnderlineNaming(tt.args.name); got != tt.want {
				t.Errorf("UnderlineNaming() = %v, want %v", got, tt.want)
			}
		})
	}
}
