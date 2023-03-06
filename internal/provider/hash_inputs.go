package provider

import (
	"crypto/sha1"
	"encoding/hex"
)

func hashInputs(data *TfSplitPoliciesDataSourceModel) string {
	if data == nil {
		return ""
	}
	// Using SHA1 is safe here as we only want a value that
	// changes whenever an input changes, this isn't used for
	// cryptographic stuff whatsoever
	var hasher = sha1.New()
	hasher.Write([]byte("tf-split-policies"))
	for _, policy := range data.Policies {
		hasher.Write([]byte(policy.ValueString()))
		hasher.Write([]byte(":"))
	}

	return hex.EncodeToString(hasher.Sum([]byte("")))
}
