package index

import (
	"time"
)

type NoteIndexEntry struct {
	Path       string
	Title      string
	Tags       []string
	Links      []string
	ModifiedAt time.Time

	// bookkeeping
	Hash string
}

type VaultIndex struct {
	Version   int
	CreatedAt time.Time
	Notes     map[string]NoteIndexEntry // key = Path
}
