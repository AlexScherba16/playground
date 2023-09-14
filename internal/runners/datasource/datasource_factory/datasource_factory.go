package datasource_factory

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	cnst "playground/internal/constants"
	"playground/internal/runners/common"
	"playground/internal/runners/datasource/runner/csv"
	"playground/internal/runners/datasource/runner/json"
	t "playground/internal/types"
	"playground/internal/utils/cerror"
	"sync"
)

// NewRunner creates a new data source runner to stream data
// In data processing pipeline
func NewRunner(
	ctx context.Context,
	wg *sync.WaitGroup,
	filePath string,
	recordCh t.RecordChannel,
	errorCh t.ErrorChannel) (common.IRunner, error) {

	// Check for file exists
	_, err := os.Stat(filePath)
	if err != nil {
		return nil, cerror.NewCustomError(fmt.Sprintf("%q no such file", filePath))
	}

	// Get file extension
	ext := filepath.Ext(filePath)

	// General Factory logic, create data source depends on file extension
	switch ext {
	case cnst.CsvDataSource:
		return csv.NewDataSourceRunner(ctx, wg, filePath, recordCh, errorCh)
	case cnst.JsonDataSource:
		return json.NewDataSourceRunner(ctx, wg, filePath, recordCh, errorCh)
	default:
		return nil, cerror.NewCustomError(fmt.Sprintf("%q invalid data source type extension", ext))
	}
}
