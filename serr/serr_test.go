package serr

import "testing"

func TestIsBadRequest(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "ok",
			args: args{
				err: BadRequest("hello"),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsBadRequest(tt.args.err); got != tt.want {
				t.Errorf("IsBadRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
