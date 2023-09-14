package predictor

func LinearExtrapolation(data []float64, day float64) float64 {
	// m = [counted_values * sum(x * y) - sum(x) * sum(y)] / [counted_values * sum(x * x) - sum(x) * sum(x)]
	// b = sum(y) - m * sum(y)
	// y = m * x + b

	var sumX, sumY, sumXY, sumXX float64
	for i, y := range data {
		x := float64(i + 1)
		sumX += x
		sumY += y
		sumXY += x * y
		sumXX += x * x
	}
	count := float64(len(data))
	m := (count*sumXY - sumX*sumY) / (count*sumXX - sumX*sumX)
	b := (sumY - m*sumX) / count

	return m*day + b
}

// IMPORTANT: Expected len(data) != 0
func Average(data []float64, day float64) float64 {
	dataLen := len(data)
	delta := (data[dataLen-1] - data[0]) / float64(dataLen)
	return data[dataLen-1] + delta*float64(day-float64(dataLen-1))
}
