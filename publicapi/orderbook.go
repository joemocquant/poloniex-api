package publicapi

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type OrderBooks map[string]OrderBook

type OrderBook struct {
	Asks     []Order `json:"asks"`
	Bids     []Order `json:"bids"`
	IsFrozen bool
	Seq      uint64
}

type Order struct {
	Rate     float64
	Quantity float64
}

func (o *OrderBook) UnmarshalJSON(data []byte) error {

	type alias OrderBook
	aux := struct {
		IsFrozen string
		*alias
	}{
		alias: (*alias)(o),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("unmarshal aux: %v", err)
	}

	if aux.IsFrozen != "0" {
		o.IsFrozen = false
	} else {
		o.IsFrozen = true
	}

	return nil
}

func (o *Order) UnmarshalJSON(data []byte) error {

	var rateStr string
	tmp := []interface{}{&rateStr, &o.Quantity}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return fmt.Errorf("unmarshal order: %v", err)
	}

	if got, want := len(tmp), 2; got != want {
		return fmt.Errorf("wrong number of fields in Order: %d != %d",
			got, want)
	}

	if val, err := strconv.ParseFloat(rateStr, 64); err != nil {
		return fmt.Errorf("parsefloat: %v", err)
	} else {
		o.Rate = val
	}

	return nil
}

func (client *PublicClient) GetOrderBooks(depth int) (OrderBooks, error) {

	params := map[string]string{
		"command":      "returnOrderBook",
		"currencyPair": "all",
		"depth":        strconv.Itoa(depth),
	}

	url := buildUrl(params)

	resp, err := client.do("GET", url, "", false)
	if err != nil {
		return nil, fmt.Errorf("get: %v", err)
	}

	res := make(OrderBooks)

	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, fmt.Errorf("json unmarshal: %v", err)
	}

	return res, nil
}

func (client *PublicClient) GetOrderBook(currencyPair string, depth int) (*OrderBook, error) {

	params := map[string]string{
		"command":      "returnOrderBook",
		"currencyPair": strings.ToUpper(currencyPair),
		"depth":        strconv.Itoa(depth),
	}

	url := buildUrl(params)

	resp, err := client.do("GET", url, "", false)
	if err != nil {
		return nil, fmt.Errorf("get: %v", err)
	}

	res := OrderBook{}

	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, fmt.Errorf("json unmarshal: %v", err)
	}

	return &res, nil
}
