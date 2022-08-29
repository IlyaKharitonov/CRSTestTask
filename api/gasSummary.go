package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"testTask/types"
)

func getGasSummary(w http.ResponseWriter, req *http.Request) {

	var (
		wg         = &sync.WaitGroup{}
		cryptoData = &types.CryptoData{}
		gasSummary = &types.GasSummary{
			GasInMonth:            make(map[string]float64, 0),
			AveragePriceGasInDay:  make(map[string]float64, 0),
			FrequencyDistribution: make(map[string]int, 0),
		}
	)

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("api.gasSummary#4 \n Error: %s\n", err.Error())
		return
	}

	err = json.Unmarshal(body, cryptoData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	wg.Add(4)
	go gasSummary.GetGasInMonth(cryptoData.Ethereum.Transactions, wg)

	go gasSummary.GetAveragePriceGasInDay(cryptoData.Ethereum.Transactions, wg)

	go gasSummary.GetFrequencyDistribution(cryptoData.Ethereum.Transactions, wg)

	go gasSummary.GetTotalAmount(cryptoData.Ethereum.Transactions, wg)

	wg.Wait()

	err = json.NewEncoder(w).Encode(gasSummary)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("api.gasSummary#5 \n Error: %s\n", err.Error())
		return
	}
}
