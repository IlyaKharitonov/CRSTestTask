package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"testTask/types"
)

type sumPriceAndCount struct {
	sum float64
	count float64
}

type sumPriceAndValue struct {
	sumPrice float64
	sumValue float64
}

func getGasSummary(w http.ResponseWriter, req *http.Request){

	var (
		//mu = &sync.Mutex{}
		wg = &sync.WaitGroup{}
		cryptoData = &types.CryptoData{}
	)

	//mu.Lock()
	gasSummary := &types.GasSummary{
		GasInMonth: make(map[string]float64, 0),
		AveragePriceGasInDay: make(map[string]float64, 0),
		FrequencyDistribution: make(map[int]int, 0),
	}
	//mu.Unlock()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("api.gasSummary#4 \n Error: %s\n", err.Error())
		return
	}

	err = json.Unmarshal(body, cryptoData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	 go gasSummary.GetGasInMonth(cryptoData.Ethereum.Transactions, wg)
	 //go gasSummary.GetGasInMonth2(cryptoData.Ethereum.Transactions, wg)
	 go gasSummary.GetAveragePriceGasInDay(cryptoData.Ethereum.Transactions)

		//err = getGasInMonth(cryptoData.Ethereum.Transactions, gasSummary)
		//if err != nil {
		//	w.WriteHeader(http.StatusInternalServerError)
		//	log.Printf("api.gasSummary#1 \n Error: %s\n", err.Error())
		//	return
		//}
		//
		//
		//
		//err = getAveragePriceGasInDay(cryptoData.Ethereum.Transactions, gasSummary)
		//if err != nil {
		//	w.WriteHeader(http.StatusInternalServerError)
		//	log.Printf("api.gasSummary#2 \n Error: %s\n", err.Error())
		//	return
		//}
		//
		//
		//
		//err = getFrequencyDistribution(cryptoData.Ethereum.Transactions, gasSummary)
		//if err != nil {
		//	w.WriteHeader(http.StatusInternalServerError)
		//	log.Printf("api.gasSummary#3 \n Error: %s\n", err.Error())
		//	return
		//}


	wg.Wait()

	err = json.NewEncoder(w).Encode(gasSummary)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("api.gasSummary#5 \n Error: %s\n", err.Error())
		return
	}
}



func getGasInMonth(transactions []types.Gas, gasSummary *types.GasSummary)error{
	rawDataForTotalAmount := &sumPriceAndValue{} // данные для поля TotalAmount

	for _,t := range transactions{
			time, err := time.Parse("06-01-02 15:04",t.Time)
			if err != nil {
				return fmt.Errorf("api.getGasInMonth #1\n Error: %s", err.Error())
			}

			year,month,_ := time.Date()
			key := month.String()+" "+(strconv.Itoa(year))

			rawDataForTotalAmount.sumPrice += t.GasPrice
			rawDataForTotalAmount.sumValue += t.GasValue

			if _,ok := gasSummary.GasInMonth[key]; ok{
				gasSummary.GasInMonth[key] += t.GasValue
				continue
			}
			gasSummary.GasInMonth[key] = t.GasValue
	}

	gasSummary.TotalAmount = rawDataForTotalAmount.sumPrice * rawDataForTotalAmount.sumValue

	return nil
}

func getAveragePriceGasInDay(transactions []types.Gas, gasSummary *types.GasSummary)error{
	rawData := make(map[string]sumPriceAndCount,0)

	for _,t := range transactions{
		time, err := time.Parse("06-01-02 15:04",t.Time)
		if err != nil {
			return fmt.Errorf("api.getAveragePriceGasInDay #1\n Error: %s", err.Error())
		}

		year,month,day := time.Date()
		key := strconv.Itoa(day)+" "+month.String()+" "+strconv.Itoa(year)

		if value, ok := rawData[key]; ok{
 			value.sum += t.GasPrice
			value.count ++
			rawData[key] = value
			continue
		}
		rawData[key] = sumPriceAndCount{
			sum: t.GasPrice,
			count: 1,
		}
	}

	for i := range rawData{
		gasSummary.AveragePriceGasInDay[i] = rawData[i].sum/rawData[i].count
	}

	return nil
}

func getFrequencyDistribution(transactions []types.Gas ,gasSummary *types.GasSummary)error{

	for _,t := range transactions{
		time, err := time.Parse("06-01-02 15:04",t.Time)
		if err != nil {
			return fmt.Errorf("api.getFrequencyDistribution #1\n Error: %s", err.Error())
		}

		key := time.Hour()

		if _,ok := gasSummary.FrequencyDistribution[key]; ok{
			gasSummary.FrequencyDistribution[key] ++
			continue
		}
		gasSummary.FrequencyDistribution[key] = 1
	}
	return nil
}

func getTotalAmount(transactions []types.Gas ,gasSummary *types.GasSummary){

}
