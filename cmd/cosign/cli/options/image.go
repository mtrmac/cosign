package options

import (
	"context"
	"fmt"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/sigstore/cosign/v2/internal/ui"
	"github.com/spf13/cobra"
)

// CriticalImageOptions allows specifying the expected image digest to operate on.
type CriticalImageOptions struct {
}

var _ Interface = (*CriticalImageOptions)(nil)

// AddFlags implements Interface
func (*CriticalImageOptions) AddFlags(cmd *cobra.Command) {
	// Nothing
}

func (o *CriticalImageOptions) ParseReference(ctx context.Context, imageRef string, nameOpts []name.Option) (name.Reference, error) {
	ref, err := name.ParseReference(imageRef, nameOpts...)
	if err != nil {
		return nil, fmt.Errorf("parsing image name %s: %w", imageRef, err)
	}
	if _, ok := ref.(name.Digest); !ok {
		msg := fmt.Sprintf(ui.TagReferenceMessage, imageRef)
		ui.Warnf(ctx, msg)
	}
	return ref, nil
}
