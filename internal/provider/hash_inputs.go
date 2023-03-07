package provider

import (
	"crypto/sha1"
	"encoding/hex"
)

func hashInputs(data *TfSplitPoliciesDataSourceModel) (string, error) {
	if data == nil {
		return "", nil
	}
	// Using SHA1 is safe here as we only want a value that
	// changes whenever an input changes, this isn't used for
	// cryptographic stuff whatsoever
	var hasher = sha1.New()
	var err error

	if _, err = hasher.Write([]byte("split-policies")); err != nil {
		return "", err
	}

	for _, policy := range data.Policies {
		if _, err = hasher.Write([]byte(policy.ValueString() + ":")); err != nil {
			return "", err
		}
	}

	return hex.EncodeToString(hasher.Sum([]byte(""))), nil
}
