package cli

import (
	"flag"
	"fmt"
	cnst "playground/internal/constants"
	err "playground/internal/utils/cerror"
)

// cliParams holds the parameters parsed from the command line.
type cliParams struct {
	model     string
	source    string
	aggregate string
}

// validateParams checks the fields of the cliParams for any missing or invalid values
// and returns an error if any required parameter is not provided.
func (c *cliParams) validateParams() error {
	type paramCheck struct {
		value string
		name  string
	}

	params := []paramCheck{
		{c.Model(), cnst.CliModelParam},
		{c.Source(), cnst.CliSourceParam},
		{c.Aggregate(), cnst.CliAggregateParam},
	}

	// Add params validation logic here.
	// Nonempty params are ok for now =)
	for _, item := range params {
		if item.value == "" {
			flag.Usage()
			return err.NewCustomError(fmt.Sprintf("%q is required", item.name))
		}
	}
	return nil
}

// Model returns the model parameter.
func (c *cliParams) Model() string {
	return c.model
}

// Source returns the source parameter.
func (c *cliParams) Source() string {
	return c.source
}

// Aggregate returns the aggregate parameter.
func (c *cliParams) Aggregate() string {
	return c.aggregate
}

// NewFlags parses command line flags and returns a populated cliParams instance.
// It returns an error if any required fields are missing.
func NewFlags() (cliParams, error) {
	cmd := cliParams{}

	flag.StringVar(&cmd.model, cnst.CliModelParam, "",
		fmt.Sprintf("The prediction method to use, example: [%s, %s]",
			cnst.LinearExtrapolationPredictorModel, cnst.AveragePredictorModel))

	flag.StringVar(&cmd.source, cnst.CliSourceParam, "", "Path to the data source file")

	flag.StringVar(&cmd.aggregate, cnst.CliAggregateParam, "",
		fmt.Sprintf("Data aggregation sign, example: [%s, %s]", cnst.AggregateCountry, cnst.AggregateCampaign))

	flag.Parse()

	// Flags validation logic
	if err := cmd.validateParams(); err != nil {
		return cliParams{}, err
	}

	return cmd, nil
}
