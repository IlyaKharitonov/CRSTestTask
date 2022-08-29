package types

import (
	"log"
	"strconv"
	"sync"
	"time"
)

const (
	keyHour  = 1
	keyDay   = 2
	keyMonth = 3
)

type GasSummary struct {
	GasInMonth            map[string]float64 `json:"gasInMonth"`
	AveragePriceGasInDay  map[string]float64 `json:"averagePriceGasInDay"`
	FrequencyDistribution map[string]int     `json:"frequencyDistribution"`
	TotalAmount           float64            `json:"totalAmount"`
	sync.Mutex
}

type SumPriceAndCount struct {
	sum   float64
	count float64
}

type sumPriceAndCount struct {
	sum   float64
	count float64
}

type sumPriceAndValue struct {
	sumPrice float64
	sumValue float64
}

func (g *GasSummary) GetGasInMonth(transactions []Gas, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, t := range transactions {
		g.addGasInMonth(t.GasValue, getKey(t, keyMonth))
	}
	return
}

func (g *GasSummary) addGasInMonth(gasValue float64, key string) {
	if _, ok := g.GasInMonth[key]; ok {
		//gasSummary.Lock()
		g.GasInMonth[key] += gasValue
		return
	}
	g.GasInMonth[key] = gasValue
	return
}

func (g *GasSummary) GetAveragePriceGasInDay(transactions []Gas, wg *sync.WaitGroup) {
	defer wg.Done()
	rawData := make(map[string]SumPriceAndCount, 0)
	for _, t := range transactions {
		g.addSumPriceAndCount(rawData, t.GasPrice, getKey(t, keyDay))
	}

	for i := range rawData {
		g.AveragePriceGasInDay[i] = rawData[i].sum / rawData[i].count
	}
	return
}

func (g *GasSummary) addSumPriceAndCount(rawData map[string]SumPriceAndCount, gasPrice float64, key string) {
	if value, ok := rawData[key]; ok {
		value.sum += gasPrice
		value.count++
		rawData[key] = value
		return
	}

	rawData[key] = SumPriceAndCount{
		sum:   gasPrice,
		count: 1,
	}
	return
}

func (g *GasSummary) GetFrequencyDistribution(transactions []Gas, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, t := range transactions {
		g.addFrequencyDistribution(getKey(t, keyHour))
	}
	return
}

func (g *GasSummary) addFrequencyDistribution(key string) {
	if _, ok := g.FrequencyDistribution[key]; ok {
		g.FrequencyDistribution[key]++
		return
	}
	g.FrequencyDistribution[key] = 1
	return
}

func (g *GasSummary) GetTotalAmount(transactions []Gas, wg *sync.WaitGroup) {
	defer wg.Done()
	rawDataForTotalAmount := &sumPriceAndValue{} // данные для поля TotalAmount

	for _, t := range transactions {
		rawDataForTotalAmount.sumPrice += t.GasPrice
		rawDataForTotalAmount.sumValue += t.GasValue
	}
	g.TotalAmount = rawDataForTotalAmount.sumPrice * rawDataForTotalAmount.sumValue
}

func getKey(t Gas, typeKey int) string {
	time, err := time.Parse("06-01-02 15:04", t.Time)
	if err != nil {
		log.Printf("types.getFrequencyDistribution #1\n Error: %s", err.Error())
	}
	year, month, day := time.Date()

	switch typeKey {
	case keyHour:
		return strconv.Itoa(time.Hour()) + " hour"
	case keyDay:
		return strconv.Itoa(day) + " " + month.String() + " " + strconv.Itoa(year)
	case keyMonth:
		return month.String() + " " + (strconv.Itoa(year))
	default:
		return ""
	}
}

//При использовании методов GetGasInMonthV2, GetAveragePriceGasInDayV2,
//котрые внутри себя обрабатывают данные в горутинах, время отклика возрастает в 2-2,5 раза
// 600-700ms против 250-280ms при использовании методов, которые обрабатывают данные последовательно

func (g *GasSummary) GetGasInMonthV2(transactions []Gas, wg *sync.WaitGroup) {
	for _, t := range transactions {
		wg.Add(1)
		go func(t Gas) {
			time, err := time.Parse("06-01-02 15:04", t.Time)
			if err != nil {
				log.Printf("types.getGasInMonth #1\n Error: %s", err.Error())
			}

			year, month, _ := time.Date()
			key := month.String() + " " + (strconv.Itoa(year))

			g.Lock()
			g.addGasInMonth(t.GasValue, key)
			g.Unlock()

			wg.Done()
		}(t)
	}
	wg.Done()
	return
}

func (g *GasSummary) GetAveragePriceGasInDayV2(transactions []Gas, wg *sync.WaitGroup) {
	rawData := make(map[string]SumPriceAndCount, 0)
	wg2 := &sync.WaitGroup{}
	for _, t := range transactions {
		wg2.Add(1)
		go func(t Gas) {
			time, err := time.Parse("06-01-02 15:04", t.Time)
			if err != nil {
				log.Printf("types.getAveragePriceGasInDay #1\n Error: %s", err.Error())
			}

			year, month, day := time.Date()
			key := strconv.Itoa(day) + " " + month.String() + " " + strconv.Itoa(year)

			g.Lock()
			g.addSumPriceAndCount(rawData, t.GasPrice, key)
			g.Unlock()
			wg2.Done()
		}(t)
	}
	wg2.Wait()

	for i := range rawData {
		g.AveragePriceGasInDay[i] = rawData[i].sum / rawData[i].count
	}
	wg.Done()
	return
}
