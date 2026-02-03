/*
 * Copyright (c) 2026. Sascha Kiefer.
 * Licensed under the MIT license. See LICENSE file in the project root for details.
 */

package classify

import "context"

type FakeClassifier struct {
	Items []Item
	Err   error
}

func (f *FakeClassifier) Classify(
	ctx context.Context,
	_ []string,
) ([]Item, error) {
	return f.Items, f.Err
}
