package download

import "narou/domain/metadata"

/*
 * Input Port
 *  └─ Interactor で実装、Controller で使用される
 */
type RequestParams struct {
	URI metadata.URI
}

/*
 * Output Port
 *  └─ Presenter で実装、Interactor で使用される
 */
type OutputPorter interface {
	OutPUtBoundary(novels []metadata.Novel) []string
}
