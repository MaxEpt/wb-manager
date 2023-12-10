package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"wb-manager/config"
	"wb-manager/internal/dto"
	"wb-manager/internal/wb"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	newPrices, err := readNewPrices()
	if err != nil {
		log.Fatal(err)
	}

	wbApi := wb.New(&cfg.WbApi)
	err = wbApi.UpdatePrices(newPrices)
	if err != nil {
		log.Fatal(err)
	}
}

func readNewPrices() (dto.WbPriceUpdateRequest, error) {
	f, err := os.Open("new_prices.csv")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	csvReader.Comma = ';'
	newPrices := make(dto.WbPriceUpdateRequest, 0)
	nmIdProcessed := make(map[int]bool)
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		nmId, err := strconv.Atoi(rec[0])
		if err != nil {
			return nil, err
		}
		price, err := strconv.Atoi(rec[1])
		if err != nil {
			return nil, err
		}

		if _, ok := nmIdProcessed[nmId]; ok {
			return nil, fmt.Errorf("duplicate nmId %d", nmId)
		}
		nmIdProcessed[nmId] = true
		newPrices = append(newPrices, dto.WbPrice{
			NmId:  nmId,
			Price: price,
		})
	}

	if len(newPrices) == 0 {
		return nil, fmt.Errorf("no prices in csv file")
	}

	return newPrices, nil
}
