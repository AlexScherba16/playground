package cli

import (
	"flag"
	"fmt"
	"os"
	cnst "playground/internal/constants"
	err "playground/internal/utils/cerror"
	"testing"
)

const (
	DefaultModelParam     = "defaultModelParam"
	DefaultSourceParam    = "defaultSourceParam"
	DefaultAggregateParam = "defaultAggregateParam"
)

func TestNewFlags(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedResult cliParams
		expectedError  bool
		errorStr       string
	}{
		{
			name:           "emptyCli",
			args:           []string{""},
			expectedResult: cliParams{},
			expectedError:  true,
			errorStr:       err.NewCustomError(fmt.Sprintf("%q is required", cnst.CliModelParam)).Error(),
		},
		{
			name:           "emptyModel",
			args:           []string{fmt.Sprintf("-%s", cnst.CliModelParam), ""},
			expectedResult: cliParams{},
			expectedError:  true,
			errorStr:       err.NewCustomError(fmt.Sprintf("%q is required", cnst.CliModelParam)).Error(),
		},
		{
			name:           "modelParameterNameTypo",
			args:           []string{fmt.Sprintf("-%serr", cnst.CliModelParam), ""},
			expectedResult: cliParams{},
			expectedError:  true,
			errorStr:       err.NewCustomError(fmt.Sprintf("%q is required", cnst.CliModelParam)).Error(),
		},
		{
			name: "emptySource",
			args: []string{
				fmt.Sprintf("-%s", cnst.CliModelParam), DefaultModelParam,
				fmt.Sprintf("-%s", cnst.CliSourceParam), "",
			},
			expectedResult: cliParams{},
			expectedError:  true,
			errorStr:       err.NewCustomError(fmt.Sprintf("%q is required", cnst.CliSourceParam)).Error(),
		},
		{
			name: "sourceParameterNameTypo",
			args: []string{
				fmt.Sprintf("-%s", cnst.CliModelParam), DefaultModelParam,
				fmt.Sprintf("-%serr", cnst.CliSourceParam), DefaultSourceParam,
			},
			expectedResult: cliParams{},
			expectedError:  true,
			errorStr:       err.NewCustomError(fmt.Sprintf("%q is required", cnst.CliSourceParam)).Error(),
		},
		{
			name: "emptyAggregate",
			args: []string{
				fmt.Sprintf("-%s", cnst.CliModelParam), DefaultModelParam,
				fmt.Sprintf("-%s", cnst.CliSourceParam), DefaultSourceParam,
				fmt.Sprintf("-%s", cnst.CliAggregateParam), "",
			},
			expectedResult: cliParams{},
			expectedError:  true,
			errorStr:       err.NewCustomError(fmt.Sprintf("%q is required", cnst.CliAggregateParam)).Error(),
		},
		{
			name: "aggregateParameterNameTypo",
			args: []string{
				fmt.Sprintf("-%s", cnst.CliModelParam), DefaultModelParam,
				fmt.Sprintf("-%s", cnst.CliSourceParam), DefaultSourceParam,
				fmt.Sprintf("-%serr", cnst.CliAggregateParam), DefaultAggregateParam,
			},
			expectedResult: cliParams{},
			expectedError:  true,
			errorStr:       err.NewCustomError(fmt.Sprintf("%q is required", cnst.CliAggregateParam)).Error(),
		},
		{
			name: "validParams",
			args: []string{
				fmt.Sprintf("-%s", cnst.CliModelParam), DefaultModelParam,
				fmt.Sprintf("-%s", cnst.CliSourceParam), DefaultSourceParam,
				fmt.Sprintf("-%s", cnst.CliAggregateParam), DefaultAggregateParam,
			},
			expectedResult: cliParams{model: DefaultModelParam, source: DefaultSourceParam, aggregate: DefaultAggregateParam},
			expectedError:  false,
			errorStr:       "",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			/* ARRANGE */

			// Update os.Args to simulate different command line arguments
			os.Args = append([]string{"test_application"}, testCase.args...)

			// Reset the flag values for each test case
			flag.CommandLine = flag.NewFlagSet(testCase.name, flag.ContinueOnError)

			/* ACT */
			flags, err := NewFlags()

			/* ASSERT */
			// Assert expected error string
			if (err != nil) && (err.Error() != testCase.errorStr) {
				t.Fatalf("NewFlags() with args %v: expected error string [%s], got [%s]", testCase.args, testCase.errorStr, err.Error())
			}

			// Assert expected error
			if (err != nil) != testCase.expectedError {
				t.Fatalf("NewFlags() with args %v: expected error %v, got %v", testCase.args, testCase.expectedError, err != nil)
			}

			// Assert result
			if flags != testCase.expectedResult {
				t.Fatalf("NewFlags() with args %v: expected %v, got %v", testCase.args, testCase.expectedResult, flags)
			}
		})
	}
}
