package types

import (
	"testing"
	"time"

	framework "paddy.dev/tf-framework-future"
)

func TestTimestampConformance(t *testing.T) {
	framework.TypeConformanceTests(
		// type we're testing
		Timestamp{},
		// sample value of that type
		TimestampValue{
			Time: time.Now(),
		})(t) // run the tests with t
}
