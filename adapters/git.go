package adapters

import (
	"fmt"
	"os"
	"sort"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/utils/merkletrie"

	"github.com/hotfile/domain"
)

type GitRepository struct {
}

// GetCommits returns all the commits, sorted from the root to the last one.
func (g *GitRepository) GetCommits() ([]domain.Commit, error) {
	currentDirectory, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	cIter, err := getLogIterator(currentDirectory)
	if err != nil {
		return nil, err
	}

	commits, err := collectCommits(err, cIter)
	if err != nil {
		return nil, err
	}

	sort.Sort(bySignatureTime(commits))

	return commits, nil
}

func collectCommits(err error, cIter object.CommitIter) ([]domain.Commit, error) {
	var commits []domain.Commit

	err = cIter.ForEach(func(c *object.Commit) error {
		commit := domain.Commit{
			SignatureTime: c.Committer.When,
		}

		parent, err := c.Parent(0)
		if err != nil {
			files, err := c.Files()
			if err != nil {
				return err
			}
			err1 := files.ForEach(func(file *object.File) error {
				commit.Changes = append(commit.Changes, domain.NewAddition(file.Name))
				return nil
			})
			commits = append(commits, commit)
			return err1
		}

		parentTree, err := parent.Tree()
		if err != nil {
			return err
		}

		currentTree, err := c.Tree()
		if err != nil {
			return err
		}

		diff, err := parentTree.Diff(currentTree)
		if err != nil {
			return err
		}

		for _, change := range diff {

			commit.Changes = append(commit.Changes, toChange(change))
		}

		commits = append(commits, commit)

		return nil
	})

	return commits, err
}

func getLogIterator(currentDirectory string) (object.CommitIter, error) {
	r, err := git.PlainOpenWithOptions(currentDirectory, &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		return nil, err
	}

	ref, err := r.Head()
	if err != nil {
		return nil, err
	}

	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return nil, err
	}
	return cIter, nil
}

func toChange(c *object.Change) domain.Change {
	action, err := c.Action()
	if err != nil {
		panic(fmt.Errorf("failed to retrieve the change action: %v", err))
	}

	if action == merkletrie.Modify && c.From.Name != c.To.Name {
		return domain.NewRename(c.From.Name, c.To.Name)
	} else if action == merkletrie.Insert {
		return domain.NewAddition(c.To.Name)
	} else if action == merkletrie.Delete {
		return domain.NewDeletion(c.From.Name)
	} else {
		return domain.NewUpdate(c.From.Name)
	}
}

// bySignatureTime implements sort.Interface to sort commits in chronological
// order.
type bySignatureTime []domain.Commit

func (a bySignatureTime) Len() int {
	return len(a)
}

func (a bySignatureTime) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a bySignatureTime) Less(i, j int) bool {
	return a[i].SignatureTime.UnixNano() < a[j].SignatureTime.UnixNano()
}
