package options

import (
	"context"
	"testing"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/sigstore/cosign/v2/internal/ui"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CriticalImageOptions_ParseReference(t *testing.T) {
	var tests = []struct {
		ref                         string
		expectedUserImageInput      string
		expectedResolvedImageDigest string // "" for nil
		expectedWarning             string
	}{
		{
			"image:bytag",
			"image:bytag",
			"",
			"WARNING: Image reference image:bytag uses a tag, not a digest",
		},
		{
			"image@sha256:be5d77c62dbe7fedfb0a4e5ec2f91078080800ab1f18358e5f31fcc8faa023c4",
			"image@sha256:be5d77c62dbe7fedfb0a4e5ec2f91078080800ab1f18358e5f31fcc8faa023c4",
			"image@sha256:be5d77c62dbe7fedfb0a4e5ec2f91078080800ab1f18358e5f31fcc8faa023c4",
			"",
		},
	}
	for _, tt := range tests {
		var userImageInput name.Reference
		var resolvedImageDigest *name.Digest
		var err error
		opts := CriticalImageOptions{}
		stderr := ui.RunWithTestCtx(func(ctx context.Context, write ui.WriteFunc) {
			userImageInput, resolvedImageDigest, err = opts.parseReference(ctx, tt.ref, nil)
		})
		require.NoError(t, err)
		if len(tt.expectedWarning) > 0 {
			assert.Contains(t, stderr, tt.expectedWarning, stderr, "bad warning message")
		} else {
			assert.Empty(t, stderr, "expected no warning")
		}
		expectedUI, err := name.ParseReference(tt.expectedUserImageInput)
		require.NoError(t, err)
		assert.Equal(t, expectedUI.Name(), userImageInput.Name())
		if tt.expectedResolvedImageDigest == "" {
			assert.Nil(t, resolvedImageDigest)
		} else {
			require.NotNil(t, resolvedImageDigest)
			expectedRID, err := name.NewDigest(tt.expectedResolvedImageDigest)
			require.NoError(t, err)
			assert.Equal(t, expectedRID.Name(), resolvedImageDigest.Name())
		}
	}
}
