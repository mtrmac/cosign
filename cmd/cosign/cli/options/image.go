package options

import (
	"context"
	"fmt"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/sigstore/cosign/v2/internal/ui"
	ociremote "github.com/sigstore/cosign/v2/pkg/oci/remote"
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

// ResolveReference parses the user-provided imageRef, and resolves it to a digest reference.
// It returns both the parsed user input, and the digest reference; canonically the return values are
// (userImageInput, resolvedImageDigest).
//
// The parsed input reflects user intent;
// the digest reference avoid a race where we use a tag multiple times, and it potentially points to different things at each access.
func (o *CriticalImageOptions) ResolveReference(ctx context.Context, imageRef string, nameOpts []name.Option, ociremoteOpts []ociremote.Option) (name.Reference, name.Digest, error) {
	ref, err := o.parseReference(ctx, imageRef, nameOpts)
	if err != nil {
		return nil, name.Digest{}, err
	}
	digest, err := ociremote.ResolveDigest(ref, ociremoteOpts...)
	if err != nil {
		return nil, name.Digest{}, err
	}
	return ref, digest, nil
}

func (o *CriticalImageOptions) parseReference(ctx context.Context, imageRef string, nameOpts []name.Option) (name.Reference, error) {
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
