package types

import (
	"time"

	"github.com/hashicorp/terraform-plugin-go/tftypes"
	framework "paddy.dev/tf-framework-future"
	"paddy.dev/tf-framework-future/diags"
)

var _ framework.Type = &Timestamp{}

type Timestamp struct {
	Time    time.Time
	Unknown bool
	Null    bool
}

func (t *Timestamp) UnderlyingType() tftypes.Type {
	return tftypes.String
}

func (t *Timestamp) ToTerraformValue() (interface{}, error) {
	if t == nil || t.Null {
		return nil, nil
	}
	return t.Time.Format(time.RFC3339), nil
}

func (t *Timestamp) FromTerraformValue(in tftypes.Value) error {
	if in.IsNull() {
		t.Null = true
		return nil
	}
	if !in.IsKnown() {
		t.Unknown = true
		return nil
	}
	var val string
	err := in.As(&val)
	if err != nil {
		return err
	}
	parsed, err := time.Parse(time.RFC3339, val)
	if err != nil {
		return err
	}
	t.Time = parsed
	return nil
}

func (t *Timestamp) Validate(in tftypes.Value) diags.Diagnostics {
	if in.IsNull() {
		return nil
	}
	if !in.IsKnown() {
		return nil
	}
	var val string
	err := in.As(&val)
	if err != nil {
		return diags.Diagnostics{
			{
				Level:       diags.Error,
				Summary:     "Invalid timestamp.",
				Description: "Timestamps must be parseable as strings. The following error occurred trying to parse this timestamp as a string:\n\n" + err.Error(),
			},
		}
	}
	_, err = time.Parse(time.RFC3339, val)
	if err != nil {
		return diags.Diagnostics{
			{
				Level:       diags.Error,
				Summary:     "Invalid timestamp format.",
				Description: "Timestamps must be in the following format: \"2006-01-02T15:04:05Z07:00\". The following error was encountered parsing the timestamp:\n\n" + err.Error(),
			},
		}
	}
	return nil
}

func (t *Timestamp) Description(format framework.TextFormat) string {
	switch format {
	default:
		return "Timestamps represent a specific point in time, formatted as 2006-01-02T15:04:05Z07:00 (2006 is the year, 01 is the month with leading zero if necessary, 02 is the day with leading zero if necessary, 15 is the 24-hour formatted hour with leading zero if necessary, 04 is the minute with leading zero if necessary, 05 is the second with leading zero if necessary, and 07:00 if the timezone offset), commonly known as RFC3339 formatting. Times are semantically equivalent and do not show up in diffs if they're the same point in time but in different timezones."
	}
}

func (t *Timestamp) ModifyPlan(state, planned tftypes.Value) (tftypes.Value, error) {
	var stateTime, plannedTime Timestamp
	err := stateTime.FromTerraformValue(state)
	if err != nil {
		return planned, err
	}
	err = plannedTime.FromTerraformValue(planned)
	if err != nil {
		return planned, err
	}
	if stateTime.Null || plannedTime.Null {
		return planned, nil
	}
	if stateTime.Unknown || plannedTime.Unknown {
		return planned, nil
	}
	if plannedTime.Time.Equal(stateTime.Time) {
		return state, nil
	}
	return planned, nil
}
