package benchmark_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/globusdigital/logbench/benchmark"
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

		require.Equal(
			t,
			1,
			fl.Type.NumIn(),
			"Setup.%s must accept io.ReadWriter as its first argument",
			fl.Name,
		)

		require.Equal(
			t,
			1,
			fl.Type.NumIn(),
			"Setup.%s must accept io.ReadWriter as its first argument",
			fl.Name,
		)

		in1 := fl.Type.In(0)
		require.True(
			t,
			fmt.Sprintf("%s.%s", in1.PkgPath(), in1.Name()) == "io.ReadWriter",
			"Setup.%s must accept io.ReadWriter as its first argument",
			fl.Name,
		)

		require.Equal(
			t,
			2,
			fl.Type.NumOut(),
			"Setup.%s must return (func(...), error)",
			fl.Name,
		)

		out1 := fl.Type.Out(0)
		require.True(
			t,
			reflect.Func == out1.Kind() && out1.NumOut() == 0,
			"Setup.%s must return (func(...), error)",
			fl.Name,
		)

		out2 := fl.Type.Out(1)
		require.True(
			t,
			out2.Kind() == reflect.Interface && out2.Name() == "error",
			"Setup.%s must return (func(...), error)",
			fl.Name,
		)
	}
}
