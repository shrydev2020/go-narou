package convert

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_addCSSClass(t *testing.T) {
	type args struct {
		body string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "japanese", args: args{`今日は2020年12月31日月曜日`},
			want: `今日は2020年<span class="text-combine">12</span>月<span class="text-combine">31</span>日月曜日`},
		{
			name: "eng", args: args{`today 12/31 sunday`},
			want: `today <span class="text-combine">12</span>/<span class="text-combine">31</span> sunday`},
		{
			name: "?/!/!!/!?/！/！？", args: args{`
a?b
a!b
a!!b
a!?b
a！b
a！？b
`},
			want: `
a<span class="text-combine">?</span>b
a<span class="text-combine">!</span>b
a<span class="text-combine">!!</span>b
a<span class="text-combine">!?</span>b
a<span class="text-combine">！</span>b
a<span class="text-combine">！？</span>b
`}, {
			name: "?/!/!!/!?/！/！？ 日本語", args: args{`
「いいい!?　ろろろろ！　はははは!!」
`},
			want: `
「いいい<span class="text-combine">!?</span>　ろろろろ<span class="text-combine">！</span>　はははは<span class="text-combine">!!</span>」
`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := addCSSClass(tt.args.body)
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("addCSSClass differs: (-got +want)\n%s", diff)
			}
		})
	}
}
