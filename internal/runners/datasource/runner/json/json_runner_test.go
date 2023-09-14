package json

import (
	c "context"
	"encoding/json"
	"fmt"
	"os"
	tp "playground/internal/types"
	"playground/internal/utils/cerror"
	"playground/internal/utils/parser"
	"reflect"
	s "sync"
	"testing"
	"time"
)

const (
	NoExFile        = "no_Json_no_ExistFile"
	InvalidJsonFile = "tmp.inv.abc.*.json"
	ValidJsonFile   = "tmp.*.json"
)

func createTempJSON(fileName string, data interface{}) (*os.File, error) {
	f, err := os.CreateTemp("", fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to create temp sjon file: %w", err)
	}

	// Encode data to json and store to file
	encoder := json.NewEncoder(f)
	if err := encoder.Encode(data); err != nil {
		return nil, err
	}

	// Close json file
	if err := f.Close(); err != nil {
		return nil, fmt.Errorf("failed to close the temp json file: %w", err)
	}
	return f, nil
}

func fieldPerUser(json *tp.JsonFileData) {
	json.Ltv1 = json.Ltv1 / float64(json.Users)
	json.Ltv2 = json.Ltv2 / float64(json.Users)
	json.Ltv3 = json.Ltv3 / float64(json.Users)
	json.Ltv4 = json.Ltv4 / float64(json.Users)
	json.Ltv5 = json.Ltv5 / float64(json.Users)
	json.Ltv6 = json.Ltv6 / float64(json.Users)
	json.Ltv7 = json.Ltv7 / float64(json.Users)
}

type inputParameters struct {
	ctx  c.Context
	wg   *s.WaitGroup
	path string
	rCh  tp.RecordChannel
	eCh  tp.ErrorChannel
}

type newDataSourceResult struct {
	dataSource *jsonDataSourceRunner
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
				t.Fatalf("NewDataSourceRunner() : expected error %v, got %v", testCase.expectedError, err != nil)
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

func TestCsvDataSource_RunInvalidJsonFileOpening(t *testing.T) {
	/* ARRANGE */
	errorStr := cerror.NewCustomError(fmt.Sprintf("failed to read json file %q", NoExFile)).Error()
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

func TestNewDataSource_RunReadValidJsonFile(t *testing.T) {
	/* ARRANGE */
	// Prepare valid json data
	jsonData := []tp.JsonFileData{
		{
			CampaignId: "9566c74d-1003-4c4d-bbbb-0407d1e2c649", Country: "TR",
			Ltv1: 1.9542502880389025, Ltv2: 1.994132946978472, Ltv3: 3.0126373791241345, Ltv4: 3.113804897018578,
			Ltv5: 3.201461265181941, Ltv6: 3.796798675112415, Ltv7: 4.321961161757773, Users: 93,
		},
		{
			CampaignId: "6694d2c4-22ac-4208-a007-2939487f6999", Country: "IT",
			Ltv1: 0.46401632345650307, Ltv2: 0.7080665558155662, Ltv3: 0.9479043587807372, Ltv4: 1.3855588020049658,
			Ltv5: 1.812878842576647, Ltv6: 2.423993387880591, Ltv7: 3.3931016433043153, Users: 97,
		},
	}

	f, err := createTempJSON(ValidJsonFile, jsonData)
	if err != nil {
		t.Fatalf("Failed to create file [%s]", err.Error())
	}
	defer os.Remove(f.Name())

	// Prepare expected data
	expectedRecords := []*tp.Record{}
	for _, json := range jsonData {
		fieldPerUser(&json)
		rec := parser.NewRecordFromJsonStruct(&json)
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

func TestNewDataSource_RunCancelReadingJsonFile(t *testing.T) {
	/* ARRANGE */
	// Prepare valid json data
	jsonData := []tp.JsonFileData{
		{
			CampaignId: "9566c74d-1003-4c4d-bbbb-0407d1e2c649", Country: "TR",
			Ltv1: 1.9542502880389025, Ltv2: 1.994132946978472, Ltv3: 3.0126373791241345, Ltv4: 3.113804897018578,
			Ltv5: 3.201461265181941, Ltv6: 3.796798675112415, Ltv7: 4.321961161757773, Users: 93,
		},
		{
			CampaignId: "6694d2c4-22ac-4208-a007-2939487f6999", Country: "IT",
			Ltv1: 0.46401632345650307, Ltv2: 0.7080665558155662, Ltv3: 0.9479043587807372, Ltv4: 1.3855588020049658,
			Ltv5: 1.812878842576647, Ltv6: 2.423993387880591, Ltv7: 3.3931016433043153, Users: 97,
		},
	}
	f, err := createTempJSON(ValidJsonFile, jsonData)
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

func TestNewDataSource_RunReadCorruptedJsonFileContent(t *testing.T) {
	/* ARRANGE */
	// Prepare valid json data
	jsonData := []map[string]interface{}{
		{
			"CampaignId": "9566c74d-1003-4c4d-bbbb-0407d1e2c649", "Country": "TR",
			"Ltv1": "1.9542502880389025", "Ltv2": "1.994132946978472", "Ltv3": "3.0126373791241345", "Ltv4": "3.113804897018578",
			"Ltv5": "3.201461265181941", "Ltv6": "3.796798675112415", "Ltv7": "4.321961161757773", "Users": "93",
		},
		{
			"CampaignId": "9566c74d-1003-4c4d-bbbb-0407d1e2c649", "Country": "TR",
			"Ltv1": "1.9542502880389025", "Ltv2": "1.994HELLO978472", "Ltv3": "3.0126373791241345", "Ltv4": "3.113804897018578",
			"Ltv5": "3.201461265181941", "Ltv6": "3.796798675112415", "Ltv7": "4.321961161757773", "Users": "93",
		},
	}
	f, err := createTempJSON(InvalidJsonFile, jsonData)
	if err != nil {
		t.Fatalf("Failed to create file [%s]", err.Error())
	}
	defer os.Remove(f.Name())
	errorStr := cerror.NewCustomError(fmt.Sprintf("failed to unmarchall json data %q", f.Name())).Error()

	// Prepare expected data
	json := tp.JsonFileData{
		CampaignId: "9566c74d-1003-4c4d-bbbb-0407d1e2c649", Country: "TR",
		Ltv1: 1.9542502880389025, Ltv2: 1.994132946978472, Ltv3: 3.0126373791241345, Ltv4: 3.113804897018578,
		Ltv5: 3.201461265181941, Ltv6: 3.796798675112415, Ltv7: 4.321961161757773, Users: 93,
	}
	fieldPerUser(&json)
	expectedRecords := []*tp.Record{parser.NewRecordFromJsonStruct(&json)}

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
