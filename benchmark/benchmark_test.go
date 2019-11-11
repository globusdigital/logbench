package benchmark_test

import (
	"logbench/benchmark"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetupFunctionTypes(t *testing.T) {
	tp := reflect.TypeOf(benchmark.Setup{})
	for i := 0; i < tp.NumField(); i++ {
		fl := tp.Field(i)
		require.Equal(
			t,
			reflect.Func,
			fl.Type.Kind(),
			"Setup.%s type is not a function",
			fl.Name,
		)
	}
}
