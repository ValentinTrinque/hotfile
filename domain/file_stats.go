package domain

import (
	"fmt"
	"sort"
	"time"
)

type FileStats struct {
	Filename       string
	TotalOfChanges uint64
	LastChangeTime time.Time
}

func (s *FileStats) Add(n uint64) {
	s.TotalOfChanges += n
}

func ComputeFilesStats(commits []Commit) []FileStats {
	filesStats := map[string]FileStats{}

	for i, commit := range commits {
		for _, change := range commit.Changes {
			fileStats, ok := filesStats[change.File()]
			if !ok {
				fileStats = newEmptyFileStats(change)
			}

			if change.IsAddition() || change.IsUpdate() {
				fileStats.TotalOfChanges += 1
				fileStats.LastChangeTime = commit.SignatureTime
				filesStats[change.File()] = fileStats
			} else if change.IsRename() {
				originalFileStats, ok := filesStats[change.OriginalFile()]
				if !ok {
					panic(fmt.Errorf("file stats were not found for \"%s\"\n%d\n%v", change.OriginalFile(), i, fileStats))
				}
				fileStats.TotalOfChanges += originalFileStats.TotalOfChanges
				fileStats.LastChangeTime = commit.SignatureTime
				filesStats[change.File()] = fileStats
				delete(filesStats, change.OriginalFile())
			} else {
				delete(filesStats, change.File())
			}
		}
	}

	filesStatsSlice := fromMapToSlice(filesStats)

	reversedFilesStatsSlice := sort.Reverse(byTotalOfChange(filesStatsSlice))
	sort.Sort(reversedFilesStatsSlice)
	return filesStatsSlice
}

func newEmptyFileStats(change Change) FileStats {
	return FileStats{
		Filename:       change.File(),
		TotalOfChanges: 0,
	}
}

func fromMapToSlice(filesStats map[string]FileStats) []FileStats {
	var s []FileStats
	for _, fileStats := range filesStats {
		s = append(s, fileStats)
	}
	return s
}

// byTotalOfChange implements sort.Interface to sort stats from highest to
// lowest total of change. If total of changes are equal then it sorts by
// alphabetical order, for predictability.
type byTotalOfChange []FileStats

func (a byTotalOfChange) Len() int {
	return len(a)
}

func (a byTotalOfChange) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a byTotalOfChange) Less(i, j int) bool {
	if a[i].TotalOfChanges == a[j].TotalOfChanges {
		return a[i].Filename < a[j].Filename
	}
	return a[i].TotalOfChanges < a[j].TotalOfChanges
}
