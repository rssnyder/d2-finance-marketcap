package main

import (
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	D2Circulating = "https://d2.finance/api/v1/supply?q=circulating"
)

func GetCirculating() (result float64, err error) {

	req, err := http.NewRequest("GET", D2Circulating, nil)
	if err != nil {
		return result, err
	}

	req.Header.Add("User-Agent", "Mozilla/5.0")
	req.Header.Add("accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}

	results, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	return strconv.ParseFloat(string(results), 0.0)
}
