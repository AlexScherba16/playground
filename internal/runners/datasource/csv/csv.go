package csv

import (
	"context"
	"encoding/csv"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	t "playground/internal/types"
	"playground/internal/utils/cerror"
	"playground/internal/utils/parser"
	"sync"
)

// csvDataSource represents a data source backed by a CSV file
type csvDataSource struct {
	ctx         context.Context
	wg          *sync.WaitGroup
	csvFilePath string
	recordCh    t.RecordChannel
	errorCh     t.ErrorChannel
}

// NewDataSource initializes and returns csvDataSource
// Returns error if some of ctx, wg, recordCh is nil, or open file failure
func NewDataSource(
	ctx context.Context,
	wg *sync.WaitGroup,
	filePath string,
	recordCh t.RecordChannel,
	errorCh t.ErrorChannel) (*csvDataSource, error) {

	// Validate parameters
	if ctx == nil {
		return nil, cerror.NewCustomError("invalid context")
	}
	if wg == nil {
		return nil, cerror.NewCustomError("invalid wait group")
	}
	if recordCh == nil {
		return nil, cerror.NewCustomError("invalid record channel")
	}
	if errorCh == nil {
		return nil, cerror.NewCustomError("invalid error channel")
	}

	return &csvDataSource{
		ctx:         ctx,
		wg:          wg,
		csvFilePath: filePath,
		recordCh:    recordCh,
		errorCh:     errorCh,
	}, nil
}

// Run interface implementation, related to CSV file specific
func (r *csvDataSource) Run() {
	go func() {
		defer r.wg.Done()
		defer close(r.recordCh)
		defer close(r.errorCh)

		// Try to open csv file
		csvFile, err := os.Open(r.csvFilePath)
		if err != nil {
			r.errorCh <- cerror.NewCustomError(fmt.Sprintf("failed to open csv file %q", r.csvFilePath))
			return
		}
		defer csvFile.Close()

		// Create a new CSV reader reading from the opened file
		reader := csv.NewReader(csvFile)

		// Scip CSV header
		_, err = reader.Read()
		if err != nil && err != io.EOF {
			r.errorCh <- cerror.NewCustomError(fmt.Sprintf("failed to read csv %q", "header"))
			return
		}

		for {
			select {
			// Handle cancel event
			case <-r.ctx.Done():
				log.Warning("csv datasource shutdown")

				// Notify next runner about cancel event
				r.recordCh <- nil
				return

			default:
				row, err := reader.Read()
				if err != nil {
					if err != io.EOF {
						r.errorCh <- cerror.NewCustomError(fmt.Sprintf("failed to read csv %q", "line"))
					}
					log.Debug("csv datasource finished work")
					return
				}

				// Convert to record csv line
				record, err := parser.NewRecordFromCsvStrings(row)
				if err != nil {
					r.errorCh <- err
					return
				}

				// Send data to next runner
				r.recordCh <- record
			}
		}
	}()
}
