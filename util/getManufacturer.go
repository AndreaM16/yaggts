package util

import (
	"github.com/AndreaM16/yaggts/model"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
)

type Q struct {
	Page uint `json:"page"`
	Size uint `json:"size"`
}

func GetManufacturers() (model.Manufacturers, error){
	baseUrl := "localhost:8080/manufacturer?page=1&size=100"
	fmt.Println(fmt.Sprintf("Using %s url for request", baseUrl))
	response, responseError := http.Get("http://" + baseUrl); if responseError != nil {
		return model.Manufacturers{}, responseError
	}
	defer response.Body.Close()
	return unmarshalManufacturer(response), nil
}

func unmarshalManufacturer(r *http.Response) model.Manufacturers {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	var manufacturers model.Manufacturers
	err = json.Unmarshal(body, &manufacturers)
	if err != nil {
		panic(err)
	}
	return manufacturers
}