package framework

import (
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"paddy.dev/tf-framework-future/diags"
)

type Type interface {
	TerraformType() tftypes.Type
	Validate(tftypes.Value) diags.Diagnostics
	Description(format TextFormat) string
	ValueFromTerraform(tftypes.Value) (Value, error)
}

type PlanModifyingType interface {
	ModifyPlan(state, planned tftypes.Value) (tftypes.Value, error)
}

type Value interface {
	ToTerraformValue() (interface{}, error)
	Equal(Value) bool
}
