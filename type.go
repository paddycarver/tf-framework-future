package framework

import (
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"paddy.dev/tf-framework-future/diags"
)

type Type interface {
	UnderlyingType() tftypes.Type
	ToTerraformValue() (interface{}, error)
	FromTerraformValue(tftypes.Value) error
	Validate(tftypes.Value) diags.Diagnostics
	Description(format TextFormat) string
}

type PlanModifyingType interface {
	ModifyPlan(state, planned tftypes.Value) (tftypes.Value, error)
}
