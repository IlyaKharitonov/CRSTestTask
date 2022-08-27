package types

type CryptoData struct {
	Ethereum Ethereum `json:"ethereum"`
}

type Ethereum struct {
	Transactions []Gas `json:"transactions"`
}

type Gas struct {
	Time           string  `json:"time"`
	GasPrice       float64 `json:"gasPrice"`
	GasValue       float64 `json:"gasValue"`
	Average        float64 `json:"average"`
	MaxGasPrice    float64 `json:"maxGasPrice"`
	MedianGasPrice float64 `json:"medianGasPrice"`
}