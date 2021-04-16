package diags

import "github.com/hashicorp/terraform-plugin-go/tftypes"

type Diagnostics []Diagnostic

type Diagnostic struct {
	Level       DiagnosticLevel
	Summary     string
	Description string
	Attribute   *tftypes.AttributePath
}

type DiagnosticLevel uint8

const (
	Invalid DiagnosticLevel = 0
	Warning DiagnosticLevel = iota
	Error   DiagnosticLevel = iota
)
