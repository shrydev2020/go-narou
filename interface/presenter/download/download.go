package download

import "narou/domain/metadata"

type download struct {
}

func NewOutPutBoundary() *download {
	return &download{}
}

//  OutPUtBoundary ダウンロード後に表示するデータ
func (download) OutPUtBoundary(novels []metadata.Novel) []string {
	var ret []string
	for _, v := range novels {
		ret = append(ret, v.Title)
	}
	return ret
}
