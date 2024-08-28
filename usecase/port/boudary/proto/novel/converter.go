package novel

import (
	"time"

	"narou/domain/metadata"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func Convert2ProtoBuf(lst []metadata.Novel) *Novels {
	ret := Novels{}
	novels := make([]*Novel, 0, len(lst))
	for _, n := range lst {
		novels = append(novels, &Novel{
			Id:              uint64(n.ID),
			Author:          n.Author,
			Title:           n.Title,
			FileTitle:       n.FileTitle,
			TopUrl:          string(n.TopUrl),
			SiteName:        n.SiteName,
			Story:           n.Story,
			NovelType:       Novel_Kind(n.NovelType),
			End:             n.End,
			LastUpdate:      convTime2Stamp(n.LastUpdate),
			NewArrivalsDate: convTime2Stamp(n.NewArrivalsDate),
			UseSubdirectory: n.UseSubdirectory,
			GeneralFirstUp:  convTime2Stamp(n.GeneralFirstUp),
			NovelUpdatedAt:  convTime2Stamp(n.NovelUpdatedAt),
			GeneralKastUp:   convTime2Stamp(n.GeneralKastUp),
			Length:          uint64(n.Length),
			Suspend:         n.Suspend,
			GeneralAllNo:    uint64(n.GeneralAllNo),
			LastCheckAt:     convTime2Stamp(n.LastCheckAt),
			Subs:            convSub2PbSub(n.Subs),
		})
	}
	ret.Novels = novels
	return &ret
}

func convSub2PbSub(subs []metadata.Sub) []*Sub {
	ret := make([]*Sub, 0, len(subs))
	for _, sub := range subs {
		ret = append(ret, &Sub{
			NovelId:      uint64(sub.NovelID),
			IndexId:      uint64(sub.IndexID),
			Href:         string(sub.Href),
			Chapter:      sub.Chapter,
			Subtitle:     sub.Subtitle,
			SubDate:      convTime2Stamp(sub.SubDate),
			SubUpdatedAt: convTime2Stamp(sub.SubUpDatedAt),
			DownloadAt:   convTime2Stamp(sub.DownloadAt),
		})
	}
	return ret
}

func convTime2Stamp(time time.Time) *timestamppb.Timestamp {
	return timestamppb.New(time)
}
