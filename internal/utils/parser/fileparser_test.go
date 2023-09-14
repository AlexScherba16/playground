package parser

import (
	cnst "playground/internal/constants"
	"playground/internal/types"
	"playground/internal/utils/cerror"
	"reflect"
	"testing"
)

const (
	UserIdStr     = "123"
	CampaignIdStr = "81855ad8-681d-4d86-91e9-1e00167939cb"
	CountryStr    = "BY"

	Ltv1Str   = "1.5499697874482206"
	Ltv1Float = 1.5499697874482206

	Ltv2Str   = "2.252663605698363"
	Ltv2Float = 2.252663605698363

	Ltv3Str   = "2.2986363323452683"
	Ltv3Float = 2.2986363323452683

	Ltv4Str   = "2.8840086432719603"
	Ltv4Float = 2.8840086432719603

	Ltv5Str   = "3.696001808588305"
	Ltv5Float = 3.696001808588305

	Ltv6Str   = "5.714436511023778"
	Ltv6Float = 5.714436511023778

	Ltv7Str   = "9.7414418135349954"
	Ltv7Float = 9.7414418135349954

	Users = 93
)

func TestNewRecordFromCsvStrings_InvalidInputData(t *testing.T) {
	tests := []struct {
		name           string
		input          []string
		expectedResult *types.Record
		expectedError  bool
		errorStr       string
	}{
		{
			name:           "inputLenLessThanExpected",
			input:          []string{UserIdStr, CampaignIdStr, CountryStr, Ltv1Str, Ltv2Str, Ltv3Str, Ltv4Str, Ltv5Str, Ltv6Str},
			expectedResult: nil,
			expectedError:  true,
			errorStr:       cerror.NewCustomError("invalid csv input data len 9").Error(),
		},
		{
			name:           "inputLenMoreThanExpected",
			input:          []string{UserIdStr, CampaignIdStr, CountryStr, Ltv1Str, Ltv2Str, Ltv3Str, Ltv4Str, Ltv5Str, Ltv6Str, Ltv7Str, Ltv7Str},
			expectedResult: nil,
			expectedError:  true,
			errorStr:       cerror.NewCustomError("invalid csv input data len 11").Error(),
		},
		{
			name:           "ltv1ParseFail",
			input:          []string{UserIdStr, CampaignIdStr, CountryStr, cnst.CsvLtv1Name, Ltv2Str, Ltv3Str, Ltv4Str, Ltv5Str, Ltv6Str, Ltv7Str},
			expectedResult: nil,
			expectedError:  true,
			errorStr:       cerror.NewCustomError("failed to convert ltv data").Error(),
		},
		{
			name:           "ltv2ParseFail",
			input:          []string{UserIdStr, CampaignIdStr, CountryStr, Ltv1Str, cnst.CsvLtv2Name, Ltv3Str, Ltv4Str, Ltv5Str, Ltv6Str, Ltv7Str},
			expectedResult: nil,
			expectedError:  true,
			errorStr:       cerror.NewCustomError("failed to convert ltv data").Error(),
		},
		{
			name:           "ltv3ParseFail",
			input:          []string{UserIdStr, CampaignIdStr, CountryStr, Ltv1Str, Ltv2Str, cnst.CsvLtv3Name, Ltv4Str, Ltv5Str, Ltv6Str, Ltv7Str},
			expectedResult: nil,
			expectedError:  true,
			errorStr:       cerror.NewCustomError("failed to convert ltv data").Error(),
		},
		{
			name:           "ltv4ParseFail",
			input:          []string{UserIdStr, CampaignIdStr, CountryStr, Ltv1Str, Ltv2Str, Ltv3Str, cnst.CsvLtv4Name, Ltv5Str, Ltv6Str, Ltv7Str},
			expectedResult: nil,
			expectedError:  true,
			errorStr:       cerror.NewCustomError("failed to convert ltv data").Error(),
		},
		{
			name:           "ltv5ParseFail",
			input:          []string{UserIdStr, CampaignIdStr, CountryStr, Ltv1Str, Ltv2Str, Ltv3Str, Ltv4Str, cnst.CsvLtv5Name, Ltv6Str, Ltv7Str},
			expectedResult: nil,
			expectedError:  true,
			errorStr:       cerror.NewCustomError("failed to convert ltv data").Error(),
		},
		{
			name:           "ltv6ParseFail",
			input:          []string{UserIdStr, CampaignIdStr, CountryStr, Ltv1Str, Ltv2Str, Ltv3Str, Ltv4Str, Ltv5Str, cnst.CsvLtv6Name, Ltv7Str},
			expectedResult: nil,
			expectedError:  true,
			errorStr:       cerror.NewCustomError("failed to convert ltv data").Error(),
		},
		{
			name:           "ltv7ParseFail",
			input:          []string{UserIdStr, CampaignIdStr, CountryStr, Ltv1Str, Ltv2Str, Ltv3Str, Ltv4Str, Ltv5Str, Ltv6Str, cnst.CsvLtv7Name},
			expectedResult: nil,
			expectedError:  true,
			errorStr:       cerror.NewCustomError("failed to convert ltv data").Error(),
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			/* ARRANGE */

			/* ACT */
			result, err := NewRecordFromCsvStrings(testCase.input)

			/* ASSERT */
			// Assert expected error
			if (err != nil) != testCase.expectedError {
				t.Fatalf("NewRecordFromCsvStrings() with args %v: expected error %v, got %v", testCase.input, testCase.expectedError, err != nil)
			}

			// Assert expected error string
			if (err != nil) && (err.Error() != testCase.errorStr) {
				t.Fatalf("NewRecordFromCsvStrings() with args %v: expected error string [%s], got [%s]", testCase.input, testCase.errorStr, err.Error())
			}

			// Assert result
			if !reflect.DeepEqual(result, testCase.expectedResult) {
				t.Fatalf("NewRecordFromCsvStrings() with args %v\nexp: %+v\ngot: %+v", testCase.input, testCase.expectedResult, result)
			}
		})
	}
}

func TestNewRecordFromCsvStrings_ValidInputData(t *testing.T) {
	/* ARRANGE */
	input := []string{UserIdStr, CampaignIdStr, CountryStr, Ltv1Str, Ltv2Str, Ltv3Str, Ltv4Str, Ltv5Str, Ltv6Str, Ltv7Str}
	invalidUserIdInput := []string{cnst.CsvUserIdName, CampaignIdStr, CountryStr, Ltv1Str, Ltv2Str, Ltv3Str, Ltv4Str, Ltv5Str, Ltv6Str, Ltv7Str}

	/* ACT */
	result, normalInputErr := NewRecordFromCsvStrings(input)
	expectedResult, invalidUserIdInputErr := NewRecordFromCsvStrings(invalidUserIdInput)

	/* ASSERT */
	// Assert expected error
	if normalInputErr != nil {
		t.Fatalf("NewRecordFromCsvStrings() with args %v: expected error %v, got %v", input, nil, normalInputErr)
	}
	if invalidUserIdInputErr != nil {
		t.Fatalf("NewRecordFromCsvStrings() with args %v: expected error %v, got %v", invalidUserIdInput, nil, invalidUserIdInputErr)
	}

	// Assert result
	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("NewRecordFromCsvStrings() exp: %+v\ngot: %+v", expectedResult, result)
	}
}

func TestNewRecordFromJsonStruct(t *testing.T) {
	/* ARRANGE */
	json := types.JsonFileData{
		CampaignId: CampaignIdStr, Country: CountryStr,
		Ltv1: Ltv1Float, Ltv2: Ltv2Float, Ltv3: Ltv3Float, Ltv4: Ltv4Float,
		Ltv5: Ltv5Float, Ltv6: Ltv6Float, Ltv7: Ltv7Float, Users: Users,
	}
	expected := types.NewRecord(CampaignIdStr, CountryStr, types.LtvCollection{
		Ltv1Float, Ltv2Float, Ltv3Float, Ltv4Float, Ltv5Float, Ltv6Float, Ltv7Float})

	/* ACT */
	result := NewRecordFromJsonStruct(&json)

	/* ASSERT */
	// Assert result
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("NewRecordFromCsvStrings() exp: %+v\ngot: %+v", expected, result)
	}
}
