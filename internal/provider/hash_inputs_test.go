package provider

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/types"
)


func getEmptyTestInputs() *TfSplitPoliciesDataSourceModel {
	return &TfSplitPoliciesDataSourceModel {
		Policies: []types.String {},
	}
}
func getTestInputs() *TfSplitPoliciesDataSourceModel {
	return &TfSplitPoliciesDataSourceModel {
		Policies: []types.String {
			types.StringValue("foo"),
			types.StringValue("bar"),
			types.StringValue("baz"),
		},
	}
}

func TestHashInputs(t *testing.T) {
	var value = hashInputs(getTestInputs())
	if diff := cmp.Diff("6147455429151f8f85ddef8ff11cc2db51487778", value); diff != "" {
		t.Errorf("Got invalid hash:\n%s", diff)
	}
}

func TestEmptyHashInputs(t *testing.T) {
	var value = hashInputs(getEmptyTestInputs())
	if diff := cmp.Diff("17ac149523008b6a682017bc07987f0d06e4585b", value); diff != "" {
		t.Errorf("Got invalid hash:\n%s", diff)
	}
}

func TestEnsureDifferentFromEmpty(t *testing.T) {
	var value1 = hashInputs(getTestInputs())
	var value2 = hashInputs(getEmptyTestInputs())
	if value1 == value2 {
		t.Errorf("The hash results are the same. Clearly something has gone wrong.")
	}
}
