package datasource_factory

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"playground/internal/types"
	"playground/internal/utils/cerror"
	"sync"
	"testing"
)

const (
	NoExFile           = "NoExistFile"
	UnsupportedFileExt = "*.acb"
	ValidCsvFile       = "tmp.*.csv"
	ValidJsonFile      = "tmp.*.json"
)

func TestNewDataSource(t *testing.T) {
	// Prepare test data
	validCsvFile, err := os.CreateTemp("", ValidCsvFile)
	if err != nil {
		t.Fatalf("Failed to create tmp csv file data [%s]", err.Error())
	}
	defer os.Remove(validCsvFile.Name())

	validJsonFile, err := os.CreateTemp("", ValidJsonFile)
	if err != nil {
		t.Fatalf("Failed to create tmp csv file data [%s]", err.Error())
	}
	defer os.Remove(validJsonFile.Name())

	unsupportedFile, _ := os.CreateTemp("", UnsupportedFileExt)
	if err != nil {
		t.Fatalf("Failed to create unsupported tmp file data [%s]", err.Error())
	}
	defer os.Remove(unsupportedFile.Name())
	ext := filepath.Ext(UnsupportedFileExt)

	tests := []struct {
		name          string
		filePath      string
		expectedError bool
		errorStr      string
	}{
		{
			name:          "FileDoesNotExist",
			filePath:      NoExFile,
			expectedError: true,
			errorStr:      cerror.NewCustomError(fmt.Sprintf("%q no such file", NoExFile)).Error(),
		},
		{
			name:          "UnsupportedFileExtension",
			filePath:      unsupportedFile.Name(),
			expectedError: true,
			errorStr:      cerror.NewCustomError(fmt.Sprintf("%q invalid data source type extension", ext)).Error(),
		},
		{
			name:     "ValidCsvFile",
			filePath: validCsvFile.Name(),
		},
		{
			name:     "ValidJsonFile",
			filePath: validJsonFile.Name(),
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			/* ARRANGE */
			// Prepare input parameters
			ctx := context.Background()
			wg := &sync.WaitGroup{}
			recordCh := types.NewRecordChannel(0)
			errorCh := types.NewErrorChannel(0)

			/* ACT */
			_, err := NewRunner(ctx, wg, testCase.filePath, recordCh, errorCh)

			/* ASSERT */
			// Assert expected error string
			if (err != nil) && (err.Error() != testCase.errorStr) {
				t.Fatalf("NewRunner() : expected error string [%s], got [%s]", testCase.errorStr, err.Error())
			}

			// Assert expected error
			if (err != nil) != testCase.expectedError {
				t.Fatalf("NewRunner() : expected error %v, got %v", testCase.expectedError, err != nil)
			}
		})
	}
}
