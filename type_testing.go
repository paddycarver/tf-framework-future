package framework

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TypeConformanceTests(typ Type, val Value) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		t.Run("test-roundtrip", typeConformanceTestsRoundtrip(typ, val))
		t.Run("test-nil", typeConformanceTestsNil(typ))
		t.Run("test-unknown", typeConformanceTestsUnknown(typ))
	}
}

func typeConformanceTestsRoundtrip(typ Type, startVal Value) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		// convert from a framework.Value to a primitive that can be used for a tftypes.Value
		underlyingValue, err := startVal.ToTerraformValue()
		if err != nil {
			t.Errorf("Error generating starting value from %T: %s", typ, err)
		}

		// make sure the underlying value they're giving us matches the underlying type they're giving us
		err = tftypes.ValidateValue(typ.TerraformType(), underlyingValue)
		if err != nil {
			t.Errorf("Got unexpected error validating %T's underlying value (%+v) against its underlying type (%s): %s", typ, underlyingValue, typ.TerraformType(), err)
		}

		// convert back from a tftypes.Value to a framework.Value
		endVal, err := typ.ValueFromTerraform(tftypes.NewValue(typ.TerraformType(), underlyingValue))
		if err != nil {
			t.Errorf("Got unexpected error creating Value for %T from Terraform value %s: %s", typ, startVal, err)
		}

		// should get the same underlying type
		if !cmp.Equal(startVal, endVal) {
			t.Errorf("Value of %T was not equal after round-tripping through ToTerraformValue and ValueFromTerraform. Diff: %s", typ, cmp.Diff(startVal, endVal))
		}
	}
}

func typeConformanceTestsNil(typ Type) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		val, err := typ.ValueFromTerraform(tftypes.NewValue(typ.TerraformType(), nil))
		if err != nil {
			t.Errorf("Got unexpected error creating Value for %T from Terraform value %s: %s", typ, tftypes.NewValue(typ.TerraformType(), nil), err)
		}
		i, err := val.ToTerraformValue()
		if err != nil {
			t.Errorf("Got unexpected error converting %T to Terraform value: %s", val, err)
		}
		if i != nil {
			t.Errorf("Expected nil, got %+v from %T when converting a Terraform value", i, val)
		}
	}
}

func typeConformanceTestsUnknown(typ Type) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		val, err := typ.ValueFromTerraform(tftypes.NewValue(typ.TerraformType(), tftypes.UnknownValue))
		if err != nil {
			t.Errorf("Got unexpected error creating Value for %T from Terraform value %s: %s", typ, tftypes.NewValue(typ.TerraformType(), tftypes.UnknownValue), err)
		}
		i, err := val.ToTerraformValue()
		if err != nil {
			t.Errorf("Got unexpected error converting %T to Terraform value: %s", val, err)
		}
		if i != tftypes.UnknownValue {
			t.Errorf("Expected tftypes.UnknownValue, got %+v from %T when converting a Terraform value", i, val)
		}
	}
}
