/*
 * Copyright (c) 2026. Sascha Kiefer.
 * Licensed under the MIT license. See LICENSE file in the project root for details.
 */

package ocr

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGoogleVisionEngine_MissingCredentialsPath(t *testing.T) {
	engine := &GoogleVisionEngine{
		CredentialsPath: "",
	}

	_, err := engine.Extract(context.Background(), "dummy.png")

	require.Error(t, err)
	require.Contains(t, err.Error(), "credentials not configured")
}

func TestGoogleVisionEngine_CredentialsFileNotFound(t *testing.T) {
	engine := &GoogleVisionEngine{
		CredentialsPath: "/this/path/does/not/exist.json",
	}

	_, err := engine.Extract(context.Background(), "dummy.png")

	require.Error(t, err)
	require.Contains(t, err.Error(), "credentials file not found")
}
