// Package zerowrap provides the setup for zerowrap/log
package zerowrap

import (
	"fmt"
	"io"

	"github.com/globusdigital/logbench/benchmark"
	"github.com/globusdigital/logbench/zerowrap/log"
)

func newInfo(out io.ReadWriter) (benchmark.FnInfo, error) {
	l := log.NewLog(out)
	return func(msg string) {
		l.Info(msg)
	}, nil
}

func newInfoFmt(out io.ReadWriter) (benchmark.FnInfoFmt, error) {
	l := log.NewLog(out)
	return func(msg string, data int) {
		l.Info(fmt.Sprintf(msg, data))
	}, nil
}

func newInfoWithErrorStack(out io.ReadWriter) (
	benchmark.FnInfoWithErrorStack, error,
) {
	l := log.NewLog(out)
	return func(msg string, err error) {
		l.WithErr(err).Info(msg)
	}, nil
}

func newError(out io.ReadWriter) (benchmark.FnError, error) {
	l := log.NewLog(out)
	return func(msg string) {
		l.Error(msg)
	}, nil
}

func newInfoWith3(out io.ReadWriter) (benchmark.FnInfoWith3, error) {
	l := log.NewLog(out)
	return func(msg string, fields *benchmark.Fields3) {
		l.
			Str(fields.Name1, fields.Value1).
			Int(fields.Name2, fields.Value2).
			Float64(fields.Name3, fields.Value3).
			Info(msg)
	}, nil
}

func newInfoWith10(out io.ReadWriter) (benchmark.FnInfoWith10, error) {
	l := log.NewLog(out)
	return func(msg string, fields *benchmark.Fields10) {
		l.
			Str(fields.Name1, fields.Value1).
			Str(fields.Name2, fields.Value2).
			Str(fields.Name3, fields.Value3).
			Bool(fields.Name4, fields.Value4).
			Str(fields.Name5, fields.Value5).
			Int(fields.Name6, fields.Value6).
			Float64(fields.Name7, fields.Value7).
			Strs(fields.Name8, fields.Value8).
			Ints(fields.Name9, fields.Value9).
			Float64s(fields.Name10, fields.Value10).
			Info(msg)
	}, nil
}

func newInfoWith10Exist(out io.ReadWriter) (
	benchmark.FnInfoWith10Exist,
	error,
) {
	fields := benchmark.NewFields10()
	l := log.NewLog(out).
		Str(fields.Name1, fields.Value1).
		Str(fields.Name2, fields.Value2).
		Str(fields.Name3, fields.Value3).
		Bool(fields.Name4, fields.Value4).
		Str(fields.Name5, fields.Value5).
		Int(fields.Name6, fields.Value6).
		Float64(fields.Name7, fields.Value7).
		Strs(fields.Name8, fields.Value8).
		Ints(fields.Name9, fields.Value9).
		Float64s(fields.Name10, fields.Value10)
	return func(msg string) {
		l.Info(msg)
	}, nil
}

// Setup initializes the zerolog based logger
func Setup() benchmark.Setup {
	return benchmark.Setup{
		Info:               newInfo,
		InfoFmt:            newInfoFmt,
		InfoWithErrorStack: newInfoWithErrorStack,
		Error:              newError,
		InfoWith3:          newInfoWith3,
		InfoWith10:         newInfoWith10,
		InfoWith10Exist:    newInfoWith10Exist,
	}
}
