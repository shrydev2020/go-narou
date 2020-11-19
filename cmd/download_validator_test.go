package cmd

import "testing"

func Test_validate(t *testing.T) {
	t.Parallel()
	type args struct {
		args []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "missing args", args: args{args: nil}, wantErr: true},
		{name: "too many args", args: args{args: []string{"a", "b"}}, wantErr: true},
		{name: "not narou", args: args{args: []string{"http://dummy"}}, wantErr: true},
		{name: "not narou", args: args{args: []string{"https://ncode.syosetu.com/"}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validate(tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
