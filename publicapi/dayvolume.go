package publicapi

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type DayVolumes struct {
	DayVolumes      map[string]DayVolume
	PrimaryCurrency map[string]float64
}

type DayVolume map[string]float64

func convertToDayVolume(value map[string]interface{}) (DayVolume, error) {

	dv := make(DayVolume)
	for k, v := range value {

		if v, ok := v.(string); ok {

			if val, err := strconv.ParseFloat(v, 64); err != nil {
				return nil, fmt.Errorf("parsefloat : %v", v)
			} else {
				dv[k] = val
			}

		} else {
			return nil, fmt.Errorf("type error: %v", v)
		}
	}
	return dv, nil
}

func (a *DayVolumes) UnmarshalJSON(data []byte) error {

	adv := make(map[string]interface{})

	if err := json.Unmarshal(data, &adv); err != nil {
		return fmt.Errorf("unmarshal adv: %v", err)
	}

	a.DayVolumes = make(map[string]DayVolume)
	a.PrimaryCurrency = make(map[string]float64)

	for key, value := range adv {

		switch value := value.(type) {

		case map[string]interface{}:
			if res, err := convertToDayVolume(value); err != nil {
				return fmt.Errorf("convert to dayvolume: %v", err)
			} else {
				a.DayVolumes[key] = res
			}
		case string:
			res, _ := strconv.ParseFloat(value, 64)
			a.PrimaryCurrency[key] = res
		default:
			return fmt.Errorf("unmarshal adv error type")
		}
	}

	return nil
}

func (client *PublicClient) GetDayVolumes() (*DayVolumes, error) {

	params := map[string]string{
		"command": "return24hVolume",
	}

	url := buildUrl(params)

	resp, err := client.do("GET", url, "", false)
	if err != nil {
		return nil, fmt.Errorf("get: %v", err)
	}

	res := DayVolumes{}

	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, fmt.Errorf("json unmarshal: %v", err)
	}

	return &res, nil
}
