//
// Copyright 2021 The Sigstore Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package options

import (
	"time"

	"github.com/spf13/cobra"
)

// SignBlobOptions is the top level wrapper for the sign-blob command.
// The new output-certificate flag is only in use when COSIGN_EXPERIMENTAL is enabled
type SignBlobOptions struct {
	Key               string
	Base64Output      bool
	Output            string // deprecated: TODO remove when the output flag is fully deprecated
	OutputSignature   string // TODO: this should be the root output file arg.
	OutputCertificate string
	SecurityKey       SecurityKeyOptions
	Fulcio            FulcioOptions
	Rekor             RekorOptions
	OIDC              OIDCOptions
	Registry          RegistryOptions
	Timeout           time.Duration
	BundlePath        string
}

var _ Interface = (*SignBlobOptions)(nil)

// AddFlags implements Interface
func (o *SignBlobOptions) AddFlags(cmd *cobra.Command) {
	o.SecurityKey.AddFlags(cmd)
	o.Fulcio.AddFlags(cmd)
	o.Rekor.AddFlags(cmd)
	o.OIDC.AddFlags(cmd)
	o.Registry.AddFlags(cmd)

	cmd.Flags().StringVar(&o.Key, "key", "",
		"path to the private key file, KMS URI or Kubernetes Secret")

	cmd.Flags().BoolVar(&o.Base64Output, "b64", true,
		"whether to base64 encode the output")

	cmd.Flags().StringVar(&o.OutputSignature, "output-signature", "",
		"write the signature to FILE")

	// TODO: remove when output flag is fully deprecated
	cmd.Flags().StringVar(&o.Output, "output", "", "write the signature to FILE")

	cmd.Flags().StringVar(&o.OutputCertificate, "output-certificate", "",
		"write the certificate to FILE")

	cmd.Flags().DurationVar(&o.Timeout, "timeout", time.Second*30,
		"HTTP Timeout defaults to 30 seconds")

	cmd.Flags().StringVar(&o.BundlePath, "bundle", "",
		"write everything required to verify the blob to a FILE")
}
