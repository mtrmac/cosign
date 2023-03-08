package options

import (
	"context"
	"testing"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/sigstore/cosign/v2/internal/ui"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ParseCriticalImageReference(t *testing.T) {
	var tests = []struct {
		ref             string
		expectedWarning string
	}{
		{"image:bytag", "WARNING: Image reference image:bytag uses a tag, not a digest"},
		{"image@sha256:be5d77c62dbe7fedfb0a4e5ec2f91078080800ab1f18358e5f31fcc8faa023c4", ""},
	}
	for _, tt := range tests {
		var parsedRef name.Reference
		var err error
		stderr := ui.RunWithTestCtx(func(ctx context.Context, write ui.WriteFunc) {
			parsedRef, err = ParseCriticalImageReference(ctx, tt.ref, nil)
		})
		require.NoError(t, err)
		if len(tt.expectedWarning) > 0 {
			assert.Contains(t, stderr, tt.expectedWarning, stderr, "bad warning message")
		} else {
			assert.Empty(t, stderr, "expected no warning")
		}
		expectedRef, err := name.ParseReference(tt.ref)
		require.NoError(t, err)
		assert.Equal(t, expectedRef.Name(), parsedRef.Name())
	}
}
