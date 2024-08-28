package download

import "narou/domain/metadata"

type Download struct {
}

func NewOutPutBoundary() *Download {
	return &Download{}
}

// OutPutBoundary ダウンロード後に表示するデータ
func (Download) OutPutBoundary(novels []metadata.Novel) []string {
	var ret []string
	for _, v := range novels {
		ret = append(ret, v.Title)
	}
	return ret
}
