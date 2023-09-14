package campaign

import (
	tp "playground/internal/types"
	"reflect"
	"testing"
)

func TestNewCampaignAggregatorStrategy(t *testing.T) {
	/* ARRANGE */
	// Prepare records and expected aggregated data
	records := []*tp.Record{
		tp.NewRecord("9566c74d-1003-4c4d-bbbb-0407d1e2c649", "JP", tp.LtvCollection{1.73305638789404, 1.7684248856061633, 2.781764692566589, 0, 0, 0, 0}),
		tp.NewRecord("6325253f-ec73-4dd7-a9e2-8bf921119c16", "US", tp.LtvCollection{1.9466884664338124, 3.166483202629052, 4.892883942338033, 0, 0, 0, 0}),
		tp.NewRecord("680b4e7c-8b76-4a1b-9d49-d4955c848621", "DE", tp.LtvCollection{1.281468676817884, 1.5047392622480078, 1.7456670792496436, 0, 0, 0, 0}),
	}
	expectedAggregatedData := []*tp.AggregatedData{}
	for _, record := range records {
		agg := tp.NewAggregatedData(record.CampaignId(), record.Ltv())
		expectedAggregatedData = append(expectedAggregatedData, agg)
	}
	strategy := NewCampaignAggregatorStrategy()

	/* ACT */
	expectedStrategy := reflect.ValueOf(campaignAggregatorStrategy).Pointer()

	/* ASSERT */
	if reflect.ValueOf(strategy).Pointer() != expectedStrategy {
		t.Fatalf("NewCampaignAggregatorStrategy() exp: %+v\ngot: %+v", expectedStrategy, strategy)
	}
	for i, data := range records {
		aggregatedResult := strategy(data)
		if !reflect.DeepEqual(expectedAggregatedData[i], aggregatedResult) {
			t.Fatalf("NewCampaignAggregatorStrategy() exp: %+v\ngot: %+v",
				expectedAggregatedData[i], aggregatedResult)
		}
	}
}
