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
	ref, userInputDigest, err := o.parseReference(ctx, imageRef, nameOpts)
	if err != nil {
		return nil, name.Digest{}, err
	}
	var digest name.Digest
	if userInputDigest != nil {
		digest = *userInputDigest
	} else {
		d, err := ociremote.ResolveDigest(ref, ociremoteOpts...)
		if err != nil {
			return nil, name.Digest{}, err
		}
		digest = d
	}
	return ref, digest, nil
}

// parseReference parses the user-provided imageRef, and returns it. If it is a digest reference, it returns that as well.
func (o *CriticalImageOptions) parseReference(ctx context.Context, imageRef string, nameOpts []name.Option) (name.Reference, *name.Digest, error) {
	ref, err := name.ParseReference(imageRef, nameOpts...)
	if err != nil {
		return nil, nil, fmt.Errorf("parsing image name %s: %w", imageRef, err)
	}
	digest := (*name.Digest)(nil)
	if d, ok := ref.(name.Digest); ok {
		digest = &d
	} else {
		msg := fmt.Sprintf(ui.TagReferenceMessage, imageRef)
		ui.Warnf(ctx, msg)
	}
	return ref, digest, nil
}
