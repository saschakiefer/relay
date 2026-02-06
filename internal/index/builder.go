package index

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/saschakiefer/relay/internal/vault"
)

type Builder struct {
	Reader vault.VaultReader
}

func NewBuilder(r vault.VaultReader) *Builder {
	return &Builder{Reader: r}
}

func (b *Builder) Build(ctx context.Context) (*VaultIndex, error) {
	notes, err := b.Reader.ListNotes(ctx)
	if err != nil {
		return nil, err
	}

	idx := &VaultIndex{
		Version:   1,
		CreatedAt: time.Now(),
		Notes:     make(map[string]NoteIndexEntry),
	}

	for _, n := range notes {
		h := hashNoteMeta(n)

		idx.Notes[n.Path] = NoteIndexEntry{
			Path:       n.Path,
			Title:      n.Title,
			Tags:       n.Tags,
			Links:      n.Links,
			ModifiedAt: n.ModifiedAt,
			Hash:       h,
		}
	}

	return idx, nil
}

func hashNoteMeta(n vault.NoteMeta) string {
	h := sha256.New()
	h.Write([]byte(n.Path))
	h.Write([]byte(n.Title))
	for _, t := range n.Tags {
		h.Write([]byte(t))
	}
	for _, l := range n.Links {
		h.Write([]byte(l))
	}
	return hex.EncodeToString(h.Sum(nil))
}
