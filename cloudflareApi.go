package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"strings"
	"errors"
)

var (
	CFBaseUrl string = "https://api.cloudflare.com/client/v4/"
)

type Error struct {
	Code int
	Message string
}

type Accounts struct {
	Success bool
	Errors []Error
	Messages []interface{}
	Result []struct {
		Id string
	}
}

type DNS struct {
	Success bool
	Errors []Error
}

type Zone struct {
	Success bool
	Errors []Error
	Messages []interface{}
	Result struct {
		Id string
		Name_Servers []string
	}
}

func AuthenticatedCloudflareReq(url string, data string, apiKey string) (*http.Response, error) {
	var req *http.Request
	var err error

	if data == "" {
		req, err = http.NewRequest("GET", url, nil)
	} else {
		req, err = http.NewRequest("POST", url, strings.NewReader(data))
	}

	if err != nil {
		return &http.Response{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer " + apiKey)
	// req.Header.Set("X-Auth-Key", apiKey)
	// req.Header.Set("X-Auth-Email", email)

	client := &http.Client{}
	return client.Do(req)
}

func GetAccountId(apiKey string) (string, error) {
	reqUrl := CFBaseUrl + "accounts"
	resp, err := AuthenticatedCloudflareReq(reqUrl, "", apiKey)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var accounts Accounts
	err = json.Unmarshal(body, &accounts)
	if err != nil {
		return "", err
	}

	if len(accounts.Errors) > 0 {
		return "", fmt.Errorf(accounts.Errors[0].Message)
	}

	return accounts.Result[0].Id, err
}

func AddZone(apiKey string, domain string, id string) ([]string, string, error) {
	reqUrl := CFBaseUrl + "zones"
	data := fmt.Sprintf("{\"name\": \"%s\", \"account\": {\"id\": \"%s\"}}", domain, id)

	resp, err := AuthenticatedCloudflareReq(reqUrl, data, apiKey)
	if err != nil {
		return []string{}, "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []string{}, "", err
	}

	var createZone Zone
	err = json.Unmarshal(body, &createZone)
	if err != nil {
		return []string{}, "", err
	}

	if len(createZone.Errors) > 0 {
		return []string{}, "", fmt.Errorf("%d: %s", createZone.Errors[0].Code, createZone.Errors[0].Message)
	}

	return createZone.Result.Name_Servers, createZone.Result.Id, err
}

func AddDnsRecords(apiKey string, domain string, zoneId string, text string) error {
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		spl := strings.Fields(line)

		var priority string 
		if len(spl) == 4 {
			priority = ""
		} else if len(spl) == 5 {
			priority = spl[4]
		} else {
			return errors.New("Invalid DNS data.")
		}

		typ := spl[0]
		pname := spl[1]
		content := spl[2]
		ttl := spl[3]
		var name string
		var proxied string

		if pname == "@" {
			name = domain
		} else {
			name = pname + "." + domain
		}

		if typ == "A" || typ == "CNAME" || typ == "AAAA" {
			proxied = "true"
		} else {
			proxied = "false"
		}

		err := AddDnsRecord(apiKey, domain, zoneId, typ, name, content, ttl, priority, proxied)
		if err != nil {
			return err
		}
	}

	return nil
}

func AddDnsRecord(apiKey string, domain string, zoneId string, typ string, name string, content string, ttl string, priority string, proxied string) error {
	reqUrl := CFBaseUrl + "zones/" + zoneId + "/dns_records"
	var data string

	if priority != "" {
		data = fmt.Sprintf("{\"type\":\"%s\",\"name\":\"%s\",\"content\":\"%s\",\"ttl\":%s,\"priority\":%s,\"proxied\":%s}", typ, name, content, ttl, priority, proxied)
	} else {
		data = fmt.Sprintf("{\"type\":\"%s\",\"name\":\"%s\",\"content\":\"%s\",\"ttl\":%s,\"proxied\":%s}", typ, name, content, ttl, proxied)
	}

	resp, err := AuthenticatedCloudflareReq(reqUrl, data, apiKey)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var dnsResp DNS
	err = json.Unmarshal(body, &dnsResp)
	if err != nil {
		return err
	}

	if len(dnsResp.Errors) > 0 {
		return fmt.Errorf("%d: %s", dnsResp.Errors[0].Code, dnsResp.Errors[0].Message)
	}

	return err
}