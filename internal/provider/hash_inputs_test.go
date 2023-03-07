package provider

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func getEmptyTestInputs() *TfSplitPoliciesDataSourceModel {
	return &TfSplitPoliciesDataSourceModel{
		Policies: []types.String{},
	}
}
func getTestInputs() *TfSplitPoliciesDataSourceModel {
	return &TfSplitPoliciesDataSourceModel{
		Policies: []types.String{
			types.StringValue("foo"),
			types.StringValue("bar"),
			types.StringValue("baz"),
		},
	}
}

func TestHashInputs(t *testing.T) {
	var value, err = hashInputs(getTestInputs())
	assert.NoError(t, err)
	if diff := cmp.Diff("89160751b83d1e981c46598c914b50c870d141f3", value); diff != "" {
		t.Errorf("Got invalid hash:\n%s", diff)
	}
}

func TestEmptyHashInputs(t *testing.T) {
	var value, err = hashInputs(getEmptyTestInputs())
	assert.NoError(t, err)
	if diff := cmp.Diff("e7607a63a2af79a17a07d6aec52471436b72a363", value); diff != "" {
		t.Errorf("Got invalid hash:\n%s", diff)
	}
}

func TestEnsureDifferentFromEmpty(t *testing.T) {
	value1, err := hashInputs(getTestInputs())
	assert.NoError(t, err)

	value2, err := hashInputs(getEmptyTestInputs())
	assert.NoError(t, err)

	if value1 == value2 {
		t.Errorf("The hash results are the same. Clearly something has gone wrong.")
	}
}
