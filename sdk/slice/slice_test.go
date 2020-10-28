package slice

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSortStrings(t *testing.T) {
	type args struct {
		s []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "test",
			args: args{s: []string{
				"1 1.html",
				"100 1.html",
				"42 33.html",
				"44 後日談（後）.html",
				"43 後日談（前）.html"}},
			want: []string{
				"1 1.html",
				"42 33.html",
				"43 後日談（前）.html",
				"44 後日談（後）.html",
				"100 1.html",
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SortStrings(tt.args.s)
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("Hogefunc differs: (-got +want)\n%s", diff)
			}
		})
	}
}
