package main

import (
	"github.com/hotfile/adapters"
	"github.com/hotfile/domain"
	"github.com/hotfile/viewers"
)

func main() {
	repository := adapters.GitRepository{}

	commits, err := repository.GetCommits()
	if err != nil {
		panic(err)
	}

	filesStats := domain.ComputeFilesStats(commits)

	viewers.DisplayFilesStats(filesStats)
}
