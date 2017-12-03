package util

import (
	"github.com/andream16/yaggts/model"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"strings"
	"strconv"
	"errors"
	"bytes"
)

type Q struct {
	Page uint `json:"page"`
	Size uint `json:"size"`
}

var baseUrl string

func GetManufacturers(flags model.Flag) (model.Manufacturers, error){
	baseUrl = strings.Join([]string{*flags.Host, strconv.Itoa(*flags.Port)}, ":")
	fullUrl := strings.Join([]string{baseUrl, *flags.Route}, "/")
	paramsUrl := strings.Join([]string{fullUrl, "page=" + strconv.Itoa(*flags.Page) + "&size=" + strconv.Itoa(*flags.Size)}, "?")
	fmt.Println(fmt.Sprintf("Using %s url for request", paramsUrl))
	response, responseError := http.Get(paramsUrl); if responseError != nil {
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

func PostTrend(trendEntry model.Trend) error {
	fmt.Println(fmt.Sprintf("Posting new trend entry for manufacturer %s . . .", trendEntry.Manufacturer))
	body, bodyError := json.Marshal(trendEntry); if bodyError != nil {
		fmt.Println(fmt.Sprintf("Unable to marshal new trend entry for manufacturer %s, got error: %s", trendEntry.Manufacturer, bodyError.Error()))
		return bodyError
	}
	response, requestErr := http.Post(strings.Join([]string{baseUrl, "trend"}, "/"), "application/json", bytes.NewBuffer(body)); if requestErr != nil {
		fmt.Println(fmt.Sprintf("Unable to post new trend entry for manufacturer %s, got error: %s", trendEntry.Manufacturer, requestErr.Error()))
		return requestErr
	}
	if response.StatusCode != http.StatusOK {
		fmt.Println(fmt.Sprintf("Unable to post new trend entry for manufacturer %s, got status code: %d", trendEntry.Manufacturer, response.StatusCode))
		return errors.New(fmt.Sprintf("Unable to post trend entry for manufacturer %s", trendEntry.Manufacturer))
	}
	fmt.Println(fmt.Sprintf("Successfully posted new trend entry for manufacturer %s. Returning.", trendEntry.Manufacturer))
	return nil
}