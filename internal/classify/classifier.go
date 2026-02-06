/*
 * Copyright (c) 2026. Sascha Kiefer.
 * Licensed under the MIT license. See LICENSE file in the project root for details.
 */

package classify

import "context"

type ItemType string

const (
	ItemTodo    ItemType = "todo"
	ItemNote    ItemType = "note"
	ItemIdea    ItemType = "idea"
	ItemProject ItemType = "project"
	ItemData    ItemType = "data"
)

type Item struct {
	Type       ItemType `json:"type"`
	Text       string   `json:"text"`
	Level      int      `json:"level"`
	Confidence float64  `json:"confidence"`
}

type Classifier interface {
	Classify(ctx context.Context, chunks []string) ([]Item, error)
}
