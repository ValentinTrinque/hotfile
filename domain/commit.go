package domain

import "time"

type Repository interface {
	GetCommits() ([]Commit, error)
}

type Commit struct {
	Changes       []Change
	SignatureTime time.Time
}

type changeAction string

const (
	Addition = changeAction("addition")
	Deletion = changeAction("deletion")
	Rename   = changeAction("rename")
	Update   = changeAction("update")
)

type Change struct {
	action       changeAction
	targetFile   string
	originalFile string
}

func NewAddition(name string) Change {
	return Change{
		targetFile: name,
		action:     Addition,
	}
}

func NewDeletion(name string) Change {
	return Change{
		targetFile: name,
		action:     Deletion,
	}
}

func NewUpdate(name string) Change {
	return Change{
		targetFile: name,
		action:     Update,
	}
}

func NewRename(from, to string) Change {
	return Change{
		targetFile:   to,
		originalFile: from,
		action:       Rename,
	}
}

func (c Change) IsRename() bool {
	return c.action == Rename
}

func (c Change) IsAddition() bool {
	return c.action == Addition
}

func (c Change) IsUpdate() bool {
	return c.action == Update
}

func (c Change) IsDeletion() bool {
	return c.action == Deletion
}

func  (c Change) File() string {
	return c.targetFile
}

func  (c Change) OriginalFile() string {
	return c.originalFile
}
