package json

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	t "playground/internal/types"
	"playground/internal/utils/cerror"
	"playground/internal/utils/parser"
	"sync"
)

// jsonDataSourceRunner represents a data source runner backed by a JSON file
type jsonDataSourceRunner struct {
	ctx          context.Context
	wg           *sync.WaitGroup
	jsonFilePath string
	recordCh     t.RecordChannel
	errorCh      t.ErrorChannel
}

// NewDataSourceRunner initializes and returns jsonDataSourceRunner
// Returns error if some of ctx, wg, recordCh, errorCh is nil
func NewDataSourceRunner(
	ctx context.Context,
	wg *sync.WaitGroup,
	filePath string,
	recordCh t.RecordChannel,
	errorCh t.ErrorChannel) (*jsonDataSourceRunner, error) {

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

	return &jsonDataSourceRunner{
		ctx:          ctx,
		wg:           wg,
		jsonFilePath: filePath,
		recordCh:     recordCh,
		errorCh:      errorCh,
	}, nil
}

// Run interface implementation, related to JSON file specific
func (r *jsonDataSourceRunner) Run() {
	go func() {
		defer r.wg.Done()
		defer close(r.recordCh)
		defer close(r.errorCh)

		// Try to open json file
		jsonDump, err := os.ReadFile(r.jsonFilePath)
		if err != nil {
			r.errorCh <- cerror.NewCustomError(fmt.Sprintf("failed to read json file %q", r.jsonFilePath))
			return
		}

		// Unmarshal the json content into a slice of jsonFileData
		var jsonData []t.JsonFileData
		err = json.Unmarshal(jsonDump, &jsonData)
		if err != nil {
			r.errorCh <- cerror.NewCustomError(fmt.Sprintf("failed to unmarchall json data %q", r.jsonFilePath))
			return
		}

		for {
			select {
			// Handle cancel event
			case <-r.ctx.Done():
				log.Warning("json datasource shutdown")

				// Notify next runner about cancel event
				r.recordCh <- nil
				return

			default:
				for _, data := range jsonData {
					// Well, as far as I understand
					// The json data contains a set of Ltv associated with the number of users, right?
					// So I divide the sample by the number of users to get ltv per user
					data.Ltv1 = data.Ltv1 / float64(data.Users)
					data.Ltv2 = data.Ltv2 / float64(data.Users)
					data.Ltv3 = data.Ltv3 / float64(data.Users)
					data.Ltv4 = data.Ltv4 / float64(data.Users)
					data.Ltv5 = data.Ltv5 / float64(data.Users)
					data.Ltv6 = data.Ltv6 / float64(data.Users)
					data.Ltv7 = data.Ltv7 / float64(data.Users)

					// Send data to next runner
					r.recordCh <- parser.NewRecordFromJsonStruct(&data)
				}
				log.Debug("json datasource finished work")
				return
			}
		}
	}()
}
