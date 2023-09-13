package parser

import (
	"fmt"
	cnst "playground/internal/constants"
	"playground/internal/types"
	"playground/internal/utils/cerror"
	"strconv"
)

// NewRecordFromCsvStrings creates a new Record from a slice of CSV strings.
// Returns error in cases of invalid slice length or data conversion failures
func NewRecordFromCsvStrings(row []string) (*types.Record, error) {
	if len(row) != cnst.CsvDataLen {
		return nil, cerror.NewCustomError(fmt.Sprintf("invalid csv input data len %d", len(row)))
	}

	ltvs := types.LtvCollection{}
	for i, val := range row[cnst.CsvLtv1Position:] {
		ltv, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return nil, cerror.NewCustomError(fmt.Sprintf("failed to convert ltv data"))
		}
		ltvs[i] = ltv
	}
	return types.NewRecord(row[cnst.CsvCampaignIdPosition], row[cnst.CsvCountryPosition], ltvs), nil
}

// NewRecordFromJsonStruct creates a new Record from a JSON struct.
func NewRecordFromJsonStruct(jsonData *types.JsonFileData) *types.Record {
	return types.NewRecord(jsonData.CampaignId, jsonData.Country, types.LtvCollection{
		jsonData.Ltv1,
		jsonData.Ltv2,
		jsonData.Ltv3,
		jsonData.Ltv4,
		jsonData.Ltv5,
		jsonData.Ltv6,
		jsonData.Ltv7,
	})
}
