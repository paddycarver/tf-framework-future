package framework

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type TypeFactory interface {
	NewInstance() Type
	NilInstance() Type
}

func TypeConformanceTests(factory TypeFactory) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		t.Run("test-roundtrip", typeConformanceTestsRoundtrip(factory))
		t.Run("test-nil", typeConformanceTestsNil(factory))
		t.Run("test-unknown", typeConformanceTestsUnknown(factory))
	}
}

func typeConformanceTestsRoundtrip(factory TypeFactory) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		typ := factory.NewInstance()
		i, err := typ.ToTerraformValue()
		if err != nil {
			t.Errorf("Got unexpected error converting %T to Terraform value: %s", typ, err)
		}
		err = tftypes.ValidateValue(typ.UnderlyingType(), i)
		if err != nil {
			t.Errorf("Got unexpected error validating %T's Terraform value (%+v) against its underlying type (%s): %s", typ, i, typ.UnderlyingType(), err)
		}
		val := tftypes.NewValue(typ.UnderlyingType(), i)
		resTyp := factory.NewInstance()
		err = resTyp.FromTerraformValue(val)
		if err != nil {
			t.Errorf("Got unexpected error converting tftypes.Value built from %T (%s) back to %T: %s", typ, val, resTyp, err)
		}
		if !cmp.Equal(typ, resTyp) {
			t.Errorf("Type %T not reported as equal after round-tripping through ToTerraformValue and FromTerraformValue. Differences: %s", typ, cmp.Diff(typ, resTyp))
		}
	}
}

func typeConformanceTestsNil(factory TypeFactory) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		typ := factory.NilInstance()
		val := tftypes.NewValue(typ.UnderlyingType(), nil)
		resTyp := factory.NewInstance()
		err := resTyp.FromTerraformValue(val)
		if err != nil {
			t.Errorf("Got unexpected error converting to %T from %s: %s", typ, val, err)
		}
		i, err := resTyp.ToTerraformValue()
		if err != nil {
			t.Errorf("Got unexpected error converting %T to Terraform value: %s", resTyp, err)
		}
		if i != nil {
			t.Errorf("Expected nil, got %+v from %T when converting a Terraform value", i, resTyp)
		}
	}
}

func typeConformanceTestsUnknown(factory TypeFactory) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		typ := factory.NewInstance()
		val := tftypes.NewValue(typ.UnderlyingType(), tftypes.UnknownValue)
		resTyp := factory.NewInstance()
		err := resTyp.FromTerraformValue(val)
		if err != nil {
			t.Errorf("Got unexpected error converting to %T from %s: %s", typ, val, err)
		}
		i, err := resTyp.ToTerraformValue()
		if err != nil {
			t.Errorf("Got unexpected error converting %T to Terraform value: %s", resTyp, err)
		}
		if i != tftypes.UnknownValue {
			t.Errorf("Expected tftypes.UnknownValue, got %+v from %T when converting a Terraform value", i, resTyp)
		}
	}
}
