package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"testTask/types"
)

func gasSummary(w http.ResponseWriter, req *http.Request){

		//mu := &sync.Mutex{}
		//mu.Lock()
		cryptoData := &types.CryptoData{}
		//mu.Unlock()
		//mu.Lock()
		gasSummary := &types.GasSummary{
			GasInMonth: make(map[string]float64, 0),
			AveragePriceGasInDay: make(map[string]float64, 0),
			FrequencyDistribution: make(map[string]int64, 0),
		}
		//mu.Unlock()


	body, err := ioutil.ReadAll(req.Body)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("api.gasSummary#1 \n Error: %s\n", err.Error())
		return
	}

	err = json.Unmarshal(body, cryptoData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//fmt.Println(crypto.Ethereum)

	go getGasInMonth(cryptoData.Ethereum.Transactions, gasSummary)

	go getAveragePriceGasInDay(cryptoData.Ethereum.Transactions, gasSummary)
	//
	//getFrequencyDistribution(cryptoData.Ethereum.Transactions, gasSummary)
	//
	//getTotalAmount(cryptoData.Ethereum.Transactions, gasSummary)


	err = json.NewEncoder(w).Encode(gasSummary)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("api.gasSummary#2 \n Error: %s\n", err.Error())
		return
	}
}

func getGasInMonth(transactions []types.Gas, gasSummary *types.GasSummary)error{
	for i := range transactions{
		time, err := time.Parse("06-01-02 15:04",transactions[i].Time)
		if err != nil {
			return fmt.Errorf("api.getGasInMonth #1\n Error: %s", err.Error())
		}

		year,month,_ := time.Date()
		key := month.String()+" "+(strconv.Itoa(year))

		if value, ok := gasSummary.GasInMonth[key]; ok{
			gasSummary.GasInMonth[key] = value + transactions[i].GasValue
			continue
		}
		gasSummary.GasInMonth[key] = transactions[i].GasValue
	}
	return nil
}

type sumPriceAndCount struct {
	sum float64
	count float64
}

type sumPriceAndValue struct {
	sumPrice float64
	sumValue float64
}

func getAveragePriceGasInDay(transactions []types.Gas, gasSummary *types.GasSummary)error{
	rawData := make(map[string]sumPriceAndCount,0)
	rawDataForTotalAmount := &sumPriceAndValue{} // данные для поля TotalAmount

	for i := range transactions{
		time, err := time.Parse("06-01-02 15:04",transactions[i].Time)
		if err != nil {
			return fmt.Errorf("api.getAveragePriceGasInDay #1\n Error: %s", err.Error())
		}

		year,month,day := time.Date()
		key := strconv.Itoa(day)+" "+month.String()+" "+strconv.Itoa(year)

		if value, ok := rawData[key]; ok{
 			value.sum += transactions[i].GasPrice
			value.count ++
			continue
		}
		rawData[key] = sumPriceAndCount{
			sum: transactions[i].GasPrice,
			count: 1,
		}

		rawDataForTotalAmount.sumPrice += transactions[i].GasPrice
		rawDataForTotalAmount.sumValue += transactions[i].GasValue
	}

	for i := range rawData{
		gasSummary.AveragePriceGasInDay[i] = rawData[i].sum/rawData[i].count
	}

	gasSummary.TotalAmount = rawDataForTotalAmount.sumPrice * rawDataForTotalAmount.sumValue

	return nil

}

func getFrequencyDistribution(transactions []types.Gas ,gasSummary *types.GasSummary){

}

func getTotalAmount(transactions []types.Gas ,gasSummary *types.GasSummary){

}
