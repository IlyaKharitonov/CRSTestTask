package types

import (
	//"fmt"
	"strconv"
	"sync"
	"time"
	"log"
)

type GasSummary struct {
	GasInMonth            map[string]float64 `json:"gasInMonth"`
	AveragePriceGasInDay  map[string]float64 `json:"averagePriceGasInDay"`
	FrequencyDistribution map[int]int   `json:"frequencyDistribution"`
	TotalAmount           float64            `json:"totalAmount"`
	sync.Mutex
}

func (g *GasSummary) GetGasInMonth(transactions []Gas, wg *sync.WaitGroup){
	for _,t := range transactions{
		wg.Add(1)
		go func(t Gas){
			time, err := time.Parse("06-01-02 15:04",t.Time)
			if err != nil {
				log.Printf("types.getGasInMonth #1\n Error: %s", err.Error())
			}

			year,month,_ := time.Date()
			key := month.String()+" "+(strconv.Itoa(year))

			g.Lock()
			g.addGasInMonth(t.GasValue, key)
			g.Unlock()
			//g.Lock()
			//if _,ok := g.GasInMonth[key]; ok{
			//	//gasSummary.Lock()
			//	g.Lock()
			//	g.GasInMonth[key] += t.GasValue
			//	g.Unlock()
			//	return
			//}
			//g.Unlock()

			//g.Lock()
			//g.GasInMonth[key] = t.GasValue
			//g.Unlock()

			wg.Done()
			return
		}(t)
	}

}

func(g *GasSummary) addGasInMonth(gasValue float64, key string){
	if _,ok := g.GasInMonth[key]; ok{
		//gasSummary.Lock()
		g.GasInMonth[key] += gasValue
		return
	}
	g.GasInMonth[key] = gasValue
	return
}



func (g *GasSummary)GetAveragePriceGasInDay(transactions []Gas){
	rawData := make(map[string]SumPriceAndCount,0)
	wg := &sync.WaitGroup{}

	for _,t := range transactions{
		wg.Add(1)
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
			//if value, ok := rawData[key]; ok {
			//	value.sum += t.GasPrice
			//	value.count++
			//	rawData[key] = value
			//	return
			//}
			//
			//rawData[key] = SumPriceAndCount{
			//	sum:   t.GasPrice,
			//	count: 1,
			//}
		}(t)
		wg.Done()
	}
	wg.Wait()

	g.Lock()
	for i := range rawData{
		g.AveragePriceGasInDay[i] = rawData[i].sum/rawData[i].count
	}
	g.Unlock()
}

func(g *GasSummary) addSumPriceAndCount(rawData map[string]SumPriceAndCount, gasPrice float64, key string) {
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



//func (g *GasSummary) GetGasInMonth2(transactions []Gas, wg *sync.WaitGroup){
//	wg.Add(1)
//	for _,t := range transactions{
//			//fmt.Printf("запущена горутина %v\n", i)
//			time, err := time.Parse("06-01-02 15:04",t.Time)
//			if err != nil {
//				log.Printf("api.getGasInMonth #1\n Error: %s", err.Error())
//			}
//
//			year,month,_ := time.Date()
//			key := month.String()+" "+(strconv.Itoa(year))
//
//			g.Lock()
//			if _,ok := g.GasInMonth[key]; ok{
//				//gasSummary.Lock()
//				g.GasInMonth[key] += t.GasValue
//				//continue
//			}
//			g.Unlock()
//
//			g.Lock()
//			g.GasInMonth[key] = t.GasValue
//			g.Unlock()
//	}
//	wg.Done()
//}

type SumPriceAndCount struct {
	sum float64
	count float64
}

//func (g *GasSummary)GetAveragePriceGasInDay(transactions []Gas, wg * sync.WaitGroup){
//	rawData := make(map[string]SumPriceAndCount,0)
//	wg.Add(1)
//	for _,t := range transactions{
//		time, err := time.Parse("06-01-02 15:04",t.Time)
//		if err != nil {
//			 log.Printf("api.getAveragePriceGasInDay #1\n Error: %s", err.Error())
//		}
//
//		year,month,day := time.Date()
//		key := strconv.Itoa(day)+" "+month.String()+" "+strconv.Itoa(year)
//
//		g.Lock()
//		if value, ok := rawData[key]; ok{
//			value.sum += t.GasPrice
//			value.count ++
//			rawData[key] = value
//			continue
//		}
//		g.Unlock()
//		g.Lock()
//		rawData[key] = SumPriceAndCount{
//			sum: t.GasPrice,
//			count: 1,
//		}
//		g.Unlock()
//	}
//
//	for i := range rawData{
//		g.AveragePriceGasInDay[i] = rawData[i].sum/rawData[i].count
//	}
//	wg.Done()
//}