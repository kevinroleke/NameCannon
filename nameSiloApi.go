package main

import (
	"fmt"
	"errors"
	"strconv"
	"net/http"
	"io/ioutil"
)

var (
	NSBaseUrl string = "https://www.namesilo.com/api/"
	Balance float64
)

func GetBalance(apiKey string) (float64, error) {
	reqUrl := fmt.Sprintf("%sgetAccountBalance?version=1&type=xml&key=%s", NSBaseUrl, apiKey)
	
	resp, err := http.Get(reqUrl)
	if err != nil {
		return 0.0, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0.0, err
	}

	xmap, err := GetXmlMap(string(body), ">")
	if err != nil {
		return 0.0, err
	}

	if xmap["namesilo>reply>detail"] != "success" {
		return 0.0, errors.New(xmap["namesilo>reply>detail"])
	}

	balance, err := strconv.ParseFloat(xmap["namesilo>reply>balance"], 64)
	if err != nil {
		return 0.0, err
	}

	return balance, err
}

func RegisterDomain(apiKey string, domain string, ns1 string, ns2 string) error {
	reqUrl := fmt.Sprintf("%sregisterDomain?version=1&type=xml&years=1&private=1&auto_renew=1&domain=%s&key=%s&ns1=%s&ns2=%s", NSBaseUrl, domain, apiKey, ns1, ns2)

	resp, err := http.Get(reqUrl)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	xmap, err := GetXmlMap(string(body), ">")
	if err != nil {
		return err
	}

	if xmap["namesilo>reply>detail"] != "success" {
		return errors.New(xmap["namesilo>reply>detail"])
	}

	orderAmount, err := strconv.ParseFloat(xmap["namesilo>reply>order_amount"], 64)
	if err != nil {
		return err
	}

	// alternative to getting balance on each buy.
	Balance -= orderAmount

	return err
}

// not currently used in NameCannon, but can be used to buy the most amount of domains possible?
func GetPrice(apiKey string, tld string) (float64, error) {
	reqUrl := fmt.Sprintf("%sgetPrices?version=1&type=xml&key=%s", NSBaseUrl, apiKey)

	resp, err := http.Get(reqUrl)
	if err != nil {
		return 0.0, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0.0, err
	}

	xmap, err := GetXmlMap(string(body), ">")
	if err != nil {
		return 0.0, err
	}
	
	if xmap["namesilo>reply>detail"] != "success" {
		return 0.0, errors.New(xmap["namesilo>reply>detail"])
	}

	// this is why we need the XmlMap trick and not just unmarshal
	// the namesilo api dynamically generates these <com>, <net>, <xyz> items...
	price, err := strconv.ParseFloat(xmap["namesilo>reply>" + tld + ">registration"], 64)
	if err != nil {
		return 0.0, err
	}

	return price, nil
}