package country

import (
	"fmt"
	tp "playground/internal/types"
	"reflect"
	"testing"
)

func TestNewPostProcessorStrategy(t *testing.T) {
	/* ARRANGE */
	//Prepare predicted data
	predictedData := []*tp.PredictedData{
		tp.NewPredictedData("JP", 123.123),
		tp.NewPredictedData("US", 9999.99999),
	}
	expectedPostProcData := []string{}
	for _, data := range predictedData {
		expectedPostProcData = append(expectedPostProcData,
			fmt.Sprintf("%s: %.2f", data.Key(), data.Predicted()))
	}
	strategy := NewPostProcessorStrategy()

	/* ACT */
	expectedStrategy := reflect.ValueOf(countryPostProcessor).Pointer()

	/* ASSERT */
	if reflect.ValueOf(strategy).Pointer() != expectedStrategy {
		t.Fatalf("NewPostProcessorStrategy() exp: %+v\ngot: %+v", expectedStrategy, strategy)
	}
	for i, data := range predictedData {
		resultPostProc := strategy(data)
		if !reflect.DeepEqual(expectedPostProcData[i], resultPostProc) {
			t.Fatalf("NewPostProcessorStrategy() exp: %+v\ngot: %+v",
				expectedPostProcData[i], resultPostProc)
		}
	}
}
