package csv

import (
	c "context"
	"fmt"
	"os"
	tp "playground/internal/types"
	"playground/internal/utils/cerror"
	"playground/internal/utils/parser"
	"reflect"
	"strings"
	s "sync"
	"testing"
	"time"
)

const (
	NoExFile       = "no_no_no_ExistFile"
	InvalidCsvFile = "tmp.inv.abc.*.csv"
	ValidCsvFile   = "tmp.*.csv"
)

func createTempCSV(fileName string, data []string) (*os.File, error) {
	f, err := os.CreateTemp("", fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to create temp csv file: %w", err)
	}

	// Write data to the temporary file
	for _, line := range data {
		if _, err := f.WriteString(line); err != nil {
			return nil, fmt.Errorf("failed to write to temp csv file: %w", err)
		}
	}

	// Close csv file
	if err := f.Close(); err != nil {
		return nil, fmt.Errorf("failed to close the temp csv file: %w", err)
	}
	return f, nil
}

type inputParameters struct {
	ctx  c.Context
	wg   *s.WaitGroup
	path string
	rCh  tp.RecordChannel
	eCh  tp.ErrorChannel
}

type newDataSourceResult struct {
	dataSource *csvDataSourceRunner
	err        error
}

func getResult(
	ctx c.Context,
	wg *s.WaitGroup,
	path string,
	rCh tp.RecordChannel,
	eCh tp.ErrorChannel) newDataSourceResult {
	ds, err := NewDataSourceRunner(ctx, wg, path, rCh, eCh)
	return newDataSourceResult{dataSource: ds, err: err}
}

func TestNewDataSource_InvalidInputParams(t *testing.T) {
	tests := []struct {
		name                                 string
		input                                inputParameters
		expectedResult                       newDataSourceResult
		expectedError                        bool
		createExpectedResultUsingInputParams bool
	}{
		{
			name:           "noContext",
			input:          inputParameters{nil, &s.WaitGroup{}, "", tp.NewRecordChannel(0), tp.NewErrorChannel(0)},
			expectedResult: newDataSourceResult{dataSource: nil, err: cerror.NewCustomError("invalid context")},
			expectedError:  true,
		},
		{
			name:           "noWaitGroup",
			input:          inputParameters{c.Background(), nil, "", tp.NewRecordChannel(0), tp.NewErrorChannel(0)},
			expectedResult: newDataSourceResult{dataSource: nil, err: cerror.NewCustomError("invalid wait group")},
			expectedError:  true,
		},
		{
			name:           "noRecordChannel",
			input:          inputParameters{c.Background(), &s.WaitGroup{}, "", nil, tp.NewErrorChannel(0)},
			expectedResult: newDataSourceResult{dataSource: nil, err: cerror.NewCustomError("invalid record channel")},
			expectedError:  true,
		},
		{
			name:           "noErrorChannel",
			input:          inputParameters{c.Background(), &s.WaitGroup{}, "", tp.NewRecordChannel(0), nil},
			expectedResult: newDataSourceResult{dataSource: nil, err: cerror.NewCustomError("invalid error channel")},
			expectedError:  true,
		},
		{
			name:                                 "noExistFile",
			input:                                inputParameters{c.Background(), &s.WaitGroup{}, NoExFile, tp.NewRecordChannel(0), tp.NewErrorChannel(0)},
			expectedError:                        false,
			createExpectedResultUsingInputParams: true,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			/* ARRANGE */
			// Yeah not so good, but I don't know how to hack reflect.DeepEqual
			if testCase.createExpectedResultUsingInputParams {
				testCase.expectedResult = getResult(testCase.input.ctx, testCase.input.wg, testCase.input.path, testCase.input.rCh, testCase.input.eCh)
			}

			/* ACT */
			result, err := NewDataSourceRunner(testCase.input.ctx, testCase.input.wg, testCase.input.path, testCase.input.rCh, testCase.input.eCh)

			/* ASSERT */
			// Assert expected error
			if (err != nil) != testCase.expectedError {
				t.Fatalf("NewDataSource() : expected error %v, got %v", testCase.expectedError, err != nil)
			}

			errorStr := ""
			if testCase.expectedResult.err != nil {
				errorStr = testCase.expectedResult.err.Error()
			}
			// Assert expected error string
			if (err != nil) && (err.Error() != errorStr) {
				t.Fatalf("NewDataSourceRunner() : expected error string [%s], got [%s]", errorStr, err.Error())
			}

			// Assert result
			if !reflect.DeepEqual(result, testCase.expectedResult.dataSource) {
				t.Fatalf("NewDataSourceRunner() exp: %+v\ngot: %+v", testCase.expectedResult.dataSource, result)
			}
		})
	}
}

func TestCsvDataSource_RunInvalidCsvFileOpening(t *testing.T) {
	/* ARRANGE */
	errorStr := cerror.NewCustomError(fmt.Sprintf("failed to open csv file %q", NoExFile)).Error()
	in := inputParameters{c.Background(), &s.WaitGroup{}, NoExFile, tp.NewRecordChannel(0), tp.NewErrorChannel(0)}
	in.wg.Add(1)
	source, _ := NewDataSourceRunner(in.ctx, in.wg, in.path, in.rCh, in.eCh)

	/* ACT */
	go source.Run()

	/* ASSERT */
	for {
		select {
		// Assert unexpected record data
		case _, ok := <-in.rCh:
			if ok {
				t.Fatalf("Run() with params %v: unexpected record channel value", in)
			}
			// Assert expected error data
		case err, ok := <-in.eCh:
			if ok {
				if err.Error() != errorStr {
					t.Fatalf("Run() : expected error string [%s], got [%s]", errorStr, err.Error())
				}
				return
			}
			// Assert potential hang situation
		case <-time.After(1 * time.Second):
			t.Fatalf("Run() : timeout")
		}
	}
}

func TestNewDataSource_RunReadCorruptedCsvFileHeader(t *testing.T) {
	/* ARRANGE */
	errorStr := cerror.NewCustomError(fmt.Sprintf("failed to read csv %q", "header")).Error()

	// Prepare quoted string, to fail data source reader
	csvData := []string{`one,two,th"ree,four`}
	f, err := createTempCSV(InvalidCsvFile, csvData)
	if err != nil {
		t.Fatalf("Failed to create file [%s]", err.Error())
	}
	defer os.Remove(f.Name())

	in := inputParameters{c.Background(), &s.WaitGroup{}, f.Name(), tp.NewRecordChannel(0), tp.NewErrorChannel(0)}
	in.wg.Add(1)
	source, _ := NewDataSourceRunner(in.ctx, in.wg, in.path, in.rCh, in.eCh)

	/* ACT */
	go source.Run()

	/* ASSERT */
	for {
		select {
		// Assert unexpected record data
		case _, ok := <-in.rCh:
			if ok {
				t.Fatalf("Run() with params %v: unexpected record channel value", in)
			}
			// Assert expected error data
		case err, ok := <-in.eCh:
			if ok {
				if err.Error() != errorStr {
					t.Fatalf("Run() : expected error string [%s], got [%s]", errorStr, err.Error())
				}
				return
			}
			// Assert potential hang situation
		case <-time.After(1 * time.Second):
			t.Fatalf("Run() : timeout")
		}
	}
}

func TestNewDataSource_RunReadValidCsvFile(t *testing.T) {
	/* ARRANGE */
	// Prepare valid csv data
	csvData := []string{
		"UserId,CampaignId,Country,Ltv1,Ltv2,Ltv3,Ltv4,Ltv5,Ltv6,Ltv7\n",
		"6,9566c74d-1003-4c4d-bbbb-0407d1e2c649,JP,1.73305638789404,1.7684248856061633,2.781764692566589,0,0,0,0\n",
		"8,6325253f-ec73-4dd7-a9e2-8bf921119c16,US,1.9466884664338124,3.166483202629052,4.892883942338033,0,0,0,0\n",
		"9,680b4e7c-8b76-4a1b-9d49-d4955c848621,DE,1.281468676817884,1.5047392622480078,1.7456670792496436,0,0,0,0\n",
	}
	f, err := createTempCSV(ValidCsvFile, csvData)
	if err != nil {
		t.Fatalf("Failed to create file [%s]", err.Error())
	}
	defer os.Remove(f.Name())

	// Prepare expected data
	expectedRecords := []*tp.Record{}
	for _, csv := range csvData[1:] {
		strs := strings.Split(strings.ReplaceAll(csv, "\n", ""), ",")
		rec, err := parser.NewRecordFromCsvStrings(strs)
		if err != nil {
			t.Fatalf("Failed to parse tmp csv file data [%s]", err.Error())
		}
		expectedRecords = append(expectedRecords, rec)
	}

	in := inputParameters{c.Background(), &s.WaitGroup{}, f.Name(), tp.NewRecordChannel(0), tp.NewErrorChannel(0)}
	in.wg.Add(1)
	source, _ := NewDataSourceRunner(in.ctx, in.wg, in.path, in.rCh, in.eCh)

	/* ACT */
	go source.Run()

	/* ASSERT */
	for {
		select {
		// Assert expected record data
		case result, ok := <-in.rCh:
			if ok {
				// Assert result
				expected := expectedRecords[0]
				if !reflect.DeepEqual(expected, result) {
					t.Fatalf("Run() exp: %+v\ngot: %+v", expected, result)
				}
				// Remove 1 element, slice as a queue )
				expectedRecords = expectedRecords[1:]

			} else {
				// Assert empty expected records list
				if len(expectedRecords) != 0 {
					t.Fatalf("Run() unexpected records slice len exp: %+v\ngot: %+v", 0, len(expectedRecords))
				}
				return
			}
			// Assert unexpected error data
		case err, ok := <-in.eCh:
			if ok {
				t.Fatalf("Run() with params %v: unexpected error channel value [%s]", in, err.Error())
			} else {
				// record channel closed, skip select case
				in.eCh = nil
			}
			// Assert potential hang situation
		case <-time.After(1 * time.Second):
			t.Fatalf("Run() : timeout")
		}
	}
}

func TestNewDataSource_RunCancelReadingCsvFile(t *testing.T) {
	/* ARRANGE */
	// Prepare valid csv data
	csvData := []string{
		"UserId,CampaignId,Country,Ltv1,Ltv2,Ltv3,Ltv4,Ltv5,Ltv6,Ltv7\n",
		"6,9566c74d-1003-4c4d-bbbb-0407d1e2c649,JP,1.73305638789404,1.7684248856061633,2.781764692566589,0,0,0,0\n",
		"8,6325253f-ec73-4dd7-a9e2-8bf921119c16,US,1.9466884664338124,3.166483202629052,4.892883942338033,0,0,0,0\n",
		"9,680b4e7c-8b76-4a1b-9d49-d4955c848621,DE,1.281468676817884,1.5047392622480078,1.7456670792496436,0,0,0,0\n",
	}
	f, err := createTempCSV(ValidCsvFile, csvData)
	if err != nil {
		t.Fatalf("Failed to create file [%s]", err.Error())
	}
	defer os.Remove(f.Name())

	// Prepare expected data
	expectedRecords := []*tp.Record{nil}
	in := inputParameters{c.Background(), &s.WaitGroup{}, f.Name(), tp.NewRecordChannel(0), tp.NewErrorChannel(0)}
	in.wg.Add(1)

	// Set cancel context
	ctx, cancel := c.WithCancel(in.ctx)
	source, _ := NewDataSourceRunner(ctx, in.wg, in.path, in.rCh, in.eCh)
	// Invoke cancel
	cancel()

	/* ACT */
	go source.Run()

	/* ASSERT */
	for {
		select {
		// Assert expected record data
		case result, ok := <-in.rCh:
			if ok {
				// Assert result
				expected := expectedRecords[0]
				if !reflect.DeepEqual(expected, result) {
					t.Fatalf("Run() exp: %+v\ngot: %+v", expected, result)
				}
				// Remove 1 element, slice as a queue )
				expectedRecords = expectedRecords[1:]

			} else {
				// Assert empty expected records list
				if len(expectedRecords) != 0 {
					t.Fatalf("Run() unexpected records slice len exp: %+v\ngot: %+v", 0, len(expectedRecords))
				}
				return
			}
			// Assert unexpected error data
		case err, ok := <-in.eCh:
			if ok {
				t.Fatalf("Run() with params %v: unexpected error channel value [%s]", in, err.Error())
			} else {
				// record channel closed, skip select case
				in.eCh = nil
			}
			// Assert potential hang situation
		case <-time.After(1 * time.Second):
			t.Fatalf("Run() : timeout")
		}
	}
}

func TestNewDataSource_RunReadCorruptedLtvDataFromCsvFileData(t *testing.T) {
	/* ARRANGE */
	errorStr := cerror.NewCustomError("failed to convert ltv data").Error()
	// Prepare invalid csv data
	csvData := []string{
		"UserId,CampaignId,Country,Ltv1,Ltv2,Ltv3,Ltv4,Ltv5,Ltv6,Ltv7\n",
		"6,9566c74d-1003-4c4d-bbbb-0407d1e2c649,JP,1.73305638789404,1.7684248856061633,2.781764692566589,0,0,0,0\n",
		"8,6325253f-ec73-4dd7-a9e2-8bf921119c16,US,1.9466884664338124,3.166483202629052,4.892883942338033,0,0,0,0\n",
		"9,680b4e7c-8b76-4a1b-9d49-d4955c848621,DE,1.281468676817884,HELLO,1.7456670792496436,0,0,0,0\n",
	}
	f, err := createTempCSV(InvalidCsvFile, csvData)
	if err != nil {
		t.Fatalf("Failed to create file [%s]", err.Error())
	}
	defer os.Remove(f.Name())

	// Prepare expected data, skip last record (corrupted data)
	expectedRecords := []*tp.Record{}
	for _, csv := range csvData[1 : len(csvData)-1] {
		strs := strings.Split(strings.ReplaceAll(csv, "\n", ""), ",")
		rec, err := parser.NewRecordFromCsvStrings(strs)
		if err != nil {
			t.Fatalf("Failed to parse tmp csv file data [%s]", err.Error())
		}
		expectedRecords = append(expectedRecords, rec)
	}

	in := inputParameters{c.Background(), &s.WaitGroup{}, f.Name(), tp.NewRecordChannel(0), tp.NewErrorChannel(0)}
	in.wg.Add(1)
	source, _ := NewDataSourceRunner(in.ctx, in.wg, in.path, in.rCh, in.eCh)

	/* ACT */
	go source.Run()

	/* ASSERT */
	for {
		select {
		// Assert expected record data
		case result, ok := <-in.rCh:
			if ok {
				// Assert result
				expected := expectedRecords[0]
				if !reflect.DeepEqual(expected, result) {
					t.Fatalf("Run() exp: %+v\ngot: %+v", expected, result)
				}
				// Remove 1 element, slice as a queue )
				expectedRecords = expectedRecords[1:]

			} else {
				// Assert empty expected records list
				if len(expectedRecords) != 0 {
					t.Fatalf("Run() unexpected records slice len exp: %+v\ngot: %+v", 0, len(expectedRecords))
				}
				return
			}
			// Assert expected error data
		case err, ok := <-in.eCh:
			if ok {
				if err.Error() != errorStr {
					t.Fatalf("Run() : expected error string [%s], got [%s]", errorStr, err.Error())
				}
				return
			} else {
				in.eCh = nil
			}
			// Assert potential hang situation
		case <-time.After(1 * time.Second):
			t.Fatalf("Run() : timeout")
		}
	}
}

func TestNewDataSource_RunReadCorruptedCsvFileContent(t *testing.T) {
	/* ARRANGE */
	errorStr := cerror.NewCustomError(fmt.Sprintf("failed to read csv %q", "line")).Error()
	// Prepare invalid csv data
	csvData := []string{
		"UserId,CampaignId,Country,Ltv1,Ltv2,Ltv3,Ltv4,Ltv5,Ltv6,Ltv7\n",
		"6,9566c74d-1003-4c4d-bbbb-0407d1e2c649,JP,1.73305638789404,1.7684248856061633,2.781764692566589,0,0,0,0\n",
		`one,two,th"ree,four,f'M'"F"\n`,
	}
	f, err := createTempCSV(InvalidCsvFile, csvData)
	if err != nil {
		t.Fatalf("Failed to create file [%s]", err.Error())
	}
	defer os.Remove(f.Name())

	// Prepare expected data, skip last record (corrupted data)
	expectedRecords := []*tp.Record{}
	for _, csv := range csvData[1 : len(csvData)-1] {
		strs := strings.Split(strings.ReplaceAll(csv, "\n", ""), ",")
		rec, err := parser.NewRecordFromCsvStrings(strs)
		if err != nil {
			t.Fatalf("Failed to parse tmp csv file data [%s]", err.Error())
		}
		expectedRecords = append(expectedRecords, rec)
	}

	in := inputParameters{c.Background(), &s.WaitGroup{}, f.Name(), tp.NewRecordChannel(0), tp.NewErrorChannel(0)}
	in.wg.Add(1)
	source, _ := NewDataSourceRunner(in.ctx, in.wg, in.path, in.rCh, in.eCh)

	/* ACT */
	go source.Run()

	/* ASSERT */
	for {
		select {
		// Assert expected record data
		case result, ok := <-in.rCh:
			if ok {
				// Assert result
				expected := expectedRecords[0]
				if !reflect.DeepEqual(expected, result) {
					t.Fatalf("Run() exp: %+v\ngot: %+v", expected, result)
				}
				// Remove 1 element, slice as a queue )
				expectedRecords = expectedRecords[1:]

			} else {
				// Assert empty expected records list
				if len(expectedRecords) != 0 {
					t.Fatalf("Run() unexpected records slice len exp: %+v\ngot: %+v", 0, len(expectedRecords))
				}
				return
			}
			// Assert expected error data
		case err, ok := <-in.eCh:
			if ok {
				if err.Error() != errorStr {
					t.Fatalf("Run() : expected error string [%s], got [%s]", errorStr, err.Error())
				}
				return
			} else {
				in.eCh = nil
			}
			// Assert potential hang situation
		case <-time.After(1 * time.Second):
			t.Fatalf("Run() : timeout")
		}
	}
}
