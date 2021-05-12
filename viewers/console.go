package viewers

import (
	"fmt"
	"sort"

	"github.com/hotfile/domain"
)

func DisplayFilesStats(filesStats []domain.FileStats) {
	sortedFilesStats := filesStats
	sort.Sort(sort.Reverse(byMostChanges(sortedFilesStats)))

	for _, fileStats := range sortedFilesStats {
		fmt.Printf("%s \t%d\t%v\n",
			fileStats.Filename,
			fileStats.TotalOfChanges,
			fileStats.LastChangeTime,
		)
	}
}

type byMostChanges []domain.FileStats

func (a byMostChanges) Len() int {
	return len(a)
}

func (a byMostChanges) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a byMostChanges) Less(i, j int) bool {
	return a[i].TotalOfChanges < a[j].TotalOfChanges
}
