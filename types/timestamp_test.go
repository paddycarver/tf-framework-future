package types

import (
	"testing"
	"time"

	framework "paddy.dev/tf-framework-future"
)

type timestampFactory struct{}

func (f timestampFactory) NilInstance() framework.Type {
	var ts *Timestamp
	return ts
}

func (f timestampFactory) NewInstance() framework.Type {
	return &Timestamp{
		Time: time.Now().Round(time.Second),
	}
}

func TestTimestampConformance(t *testing.T) {
	framework.TypeConformanceTests(timestampFactory{})(t)
}
