package domain_test

import (
	"testing"
	"time"

	"github.com/hotfile/domain"
	"github.com/stretchr/testify/assert"
)

func Test_ShouldCountAddedFiles(t *testing.T) {
	// given
	commits := []domain.Commit{
		{
			Changes: []domain.Change{
				domain.NewAddition("test1.1"),
				domain.NewAddition("test1.2"),
			},
			SignatureTime: time.Unix(10, 0),
		},
		{
			Changes: []domain.Change{
				domain.NewAddition("test2"),
			},
			SignatureTime: time.Unix(20, 0),
		},
	}

	// when
	filesStats := domain.ComputeFilesStats(commits)

	// then
	assert.Equal(t, []domain.FileStats{
		{
			Filename:       "test2",
			TotalOfChanges: 1,
			LastChangeTime: time.Unix(20, 0),
		}, {
			Filename:       "test1.2",
			TotalOfChanges: 1,
			LastChangeTime: time.Unix(10, 0),
		}, {
			Filename:       "test1.1",
			TotalOfChanges: 1,
			LastChangeTime: time.Unix(10, 0),
		},
	}, filesStats)
}

func Test_ShouldCountUpdatedFiles(t *testing.T) {
	// given
	commits := []domain.Commit{
		{
			Changes: []domain.Change{
				domain.NewAddition("test1.1"),
				domain.NewAddition("test1.2"),
			},
			SignatureTime: time.Unix(10, 0),
		},
		{
			Changes: []domain.Change{
				domain.NewUpdate("test1.1"),
				domain.NewUpdate("test1.2"),
				domain.NewAddition("test2"),
			},
			SignatureTime: time.Unix(20, 0),
		},
	}

	// when
	filesStats := domain.ComputeFilesStats(commits)

	// then
	assert.Equal(t, []domain.FileStats{
		{
			Filename:       "test1.2",
			TotalOfChanges: 2,
			LastChangeTime: time.Unix(20, 0),
		}, {
			Filename:       "test1.1",
			TotalOfChanges: 2,
			LastChangeTime: time.Unix(20, 0),
		}, {
			Filename:       "test2",
			TotalOfChanges: 1,
			LastChangeTime: time.Unix(20, 0),
		},
	}, filesStats)
}

func Test_ShouldCountRenamedFiles(t *testing.T) {
	// given
	commits := []domain.Commit{
		{
			Changes: []domain.Change{
				domain.NewAddition("test1.1"),
				domain.NewAddition("test1.2"),
			},
			SignatureTime: time.Unix(10, 0),
		}, {
			Changes: []domain.Change{
				domain.NewUpdate("test1.1"),
				domain.NewUpdate("test1.2"),
			},
			SignatureTime: time.Unix(20, 0),
		}, {
			Changes: []domain.Change{
				domain.NewRename("test1.1", "test2.1"),
				domain.NewRename("test1.2", "test2.2"),
				domain.NewAddition("test2.3"),
			},
			SignatureTime: time.Unix(30, 0),
		},
	}

	// when
	filesStats := domain.ComputeFilesStats(commits)

	// then
	assert.Equal(t, []domain.FileStats{
		{
			Filename:       "test2.2",
			TotalOfChanges: 2,
			LastChangeTime: time.Unix(30, 0),
		}, {
			Filename:       "test2.1",
			TotalOfChanges: 2,
			LastChangeTime: time.Unix(30, 0),
		}, {
			Filename:       "test2.3",
			TotalOfChanges: 1,
			LastChangeTime: time.Unix(30, 0),
		},
	}, filesStats)
}

func Test_ShouldNoCountDeletedFiles(t *testing.T) {
	// given
	commits := []domain.Commit{
		{
			Changes: []domain.Change{
				domain.NewAddition("test1.1"),
				domain.NewAddition("test1.2"),
			},
			SignatureTime: time.Unix(10, 0),
		}, {
			Changes: []domain.Change{
				domain.NewUpdate("test1.1"),
				domain.NewUpdate("test1.2"),
			},
			SignatureTime: time.Unix(20, 0),
		}, {
			Changes: []domain.Change{
				domain.NewDeletion("test1.1"),
				domain.NewDeletion("test1.2"),
				domain.NewAddition("test2.3"),
			},
			SignatureTime: time.Unix(30, 0),
		},
	}

	// when
	filesStats := domain.ComputeFilesStats(commits)

	// then
	assert.Equal(t, []domain.FileStats{
		{
			Filename:       "test2.3",
			TotalOfChanges: 1,
			LastChangeTime: time.Unix(30, 0),
		},
	}, filesStats)
}
