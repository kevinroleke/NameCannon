package main

import (
	"log"
	"flag"
	"os"
	"time"
	"strings"
	"io/ioutil"
)

func HandleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func addToCloudflare(apiKey string, domain string, dnsRecords string) error {
	id, err := GetAccountId(apiKey)
	if err != nil {
		return err
	}

	log.Println(id)
	if err != nil {
		return err
	}

	_, zoneId, err := AddZone(apiKey, domain, id)
	if err != nil {
		if strings.HasPrefix(err.Error(), "1049") {
			log.Println("Domain not registered. Waiting 120 seconds.")
			time.Sleep(120 * time.Second)
			return addToCloudflare(apiKey, domain, dnsRecords)
		} else {
			return err
		}
	}

	err = AddDnsRecords(apiKey, domain, zoneId, dnsRecords)
	return err
}

func main() {
	namesiloApiKey := flag.String("namesiloSecret", "", "NameSilo API key.")
	cloudflareApiKey := flag.String("cloudflareSecret", "", "Cloudflare API key.")
	dnsRecordsFile := flag.String("dnsRecordsFile", "", "Filename containing DNS records.")
	domainsFile := flag.String("domainsFile", "", "Filename containing a newline deliminated list of domains to register.")

	ns1 := flag.String("ns1", "", "Cloudflare nameserver 1.")
	ns2 := flag.String("ns2", "", "Cloudflare nameserver 2.")

	flag.Parse()

	if *namesiloApiKey == "" || *cloudflareApiKey == "" {
		log.Fatal("Please specify Namesilo and Cloudflare API keys with --namesiloSecret, --cloudflareSecret and -cloudflareEmail")
	}

	if *dnsRecordsFile == "" {
		log.Fatal("Please specify a DNS records file with --dnsRecordsFile.")
	}

	if *ns1 == "" || *ns2 == "" {
		log.Fatal("Please specify your Cloudflare nameservers with --ns1 and --ns2.")
	}

	domain := "jeffersonaeroplane.xyz"

	var err error
	Balance, err = GetBalance(*namesiloApiKey)
	HandleErr(err)

	log.Printf("Balance is $%f\n", Balance)
	log.Printf("Buying domain %s\n", domain)

	if Balance > 10 {
		err = RegisterDomain(*namesiloApiKey, domain, *ns1, *ns2)
		HandleErr(err)
	} else {
		log.Fatal("Low balance")
	}

	Balance, err = GetBalance(*namesiloApiKey)
	HandleErr(err)

	log.Printf("Balance is $%f\n", Balance)

	f, err := os.Open(*dnsRecordsFile)
	HandleErr(err)

	dnsRecords, err := ioutil.ReadAll(f)
	HandleErr(err)

	err = addToCloudflare(*cloudflareApiKey, domain, string(dnsRecords))
	HandleErr(err)
}