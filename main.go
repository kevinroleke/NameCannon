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

func RemoveIndex(slice []string, s int) []string {
    return append(slice[:s], slice[s+1:]...)
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
	lowBal := flag.Float64("lowBalanceLimit", 0, "If balance is below this number, don't register more names.")

	ns1 := flag.String("ns1", "", "Cloudflare nameserver 1.")
	ns2 := flag.String("ns2", "", "Cloudflare nameserver 2.")

	flag.Parse()

	if *namesiloApiKey == "" || *cloudflareApiKey == "" {
		log.Fatal("Please specify Namesilo and Cloudflare API keys with --namesiloSecret, --cloudflareSecret and -cloudflareEmail")
	}

	if *dnsRecordsFile == "" {
		log.Fatal("Please specify a DNS records file with --dnsRecordsFile.")
	}

	if *domainsFile == "" {
		log.Fatal("Please specify a domain list file with --domainsFile.")
	}

	if *ns1 == "" || *ns2 == "" {
		log.Fatal("Please specify your Cloudflare nameservers with --ns1 and --ns2.")
	}

	f, err := os.Open(*dnsRecordsFile)
	HandleErr(err)

	dnsRecords, err := ioutil.ReadAll(f)
	HandleErr(err)

	f, err = os.Open(*domainsFile)
	HandleErr(err)

	domainsList, err := ioutil.ReadAll(f)
	HandleErr(err)

	domains := strings.Split(string(domainsList), "\n")

	for i, domain := range domains {
		Balance, err = GetBalance(*namesiloApiKey)
		HandleErr(err)

		log.Printf("Balance is $%f\n", Balance)

		if Balance > *lowBal {
			log.Printf("Buying domain %s\n", domain)
			err = RegisterDomain(*namesiloApiKey, domain, *ns1, *ns2)
			if err != nil {
				log.Println(err)
				domains = RemoveIndex(domains, i)
			}
		} else {
			log.Fatal("Low balance")
		}
	}

	for _, domain := range domains {
		err = addToCloudflare(*cloudflareApiKey, domain, string(dnsRecords))
		HandleErr(err)
	}
}