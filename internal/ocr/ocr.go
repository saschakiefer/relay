/*
 * Copyright (c) 2026. Sascha Kiefer.
 * Licensed under the MIT license. See LICENSE file in the project root for details.
 */

package ocr

import "context"

type Engine interface {
	Extract(ctx context.Context, inputPath string) (string, error)
}
