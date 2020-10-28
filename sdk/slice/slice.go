package slice

import (
	"github.com/mattn/natural"
)

func SortStrings(s []string) []string {
	natural.Sort(s)
	return s
}
