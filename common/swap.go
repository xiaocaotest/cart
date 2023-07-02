package common

import "encoding/json"

func SwapTo(request, product interface{}) (err error) {
	dataByte, err := json.Marshal(request)
	if err != nil {
		return
	}
	return json.Unmarshal(dataByte, product)
}
