package types


type GasSummary struct {
	GasInMonth            map[string]float64 `json:"gasInMonth"`
	AveragePriceGasInDay  map[string]float64 `json:"averagePriceGasInDay"`
	FrequencyDistribution map[string]int64   `json:"frequencyDistribution"`
	TotalAmount           float64            `json:"totalAmount"`

}
