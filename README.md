# NameCannon

Automatically register a list of domain names, add them as zones on Cloudflare, then add DNS records.

## Usage

```sh
$ ./NameCannon --namesiloSecret <namesilo api key> --cloudflareSecret <cloudflare token> --dnsRecordsFile <path/to/dnsFile.txt> --domainsFile <path/to/listOfDomains.txt>
```

### api keys

- Create Namesilo account, fund it in Account Funds Manager, create an API key. 
- Create a Cloudflare account, create an API Token, give it permissions of edit on Zone\*, read account settings, read user details.  

### dns file format

```
type subdomain answer ttl (priority)?

A @ 1.1.1.1 3600
MX @ google.com 3600 15
CNAME cats 1.0.0.1 3600
```
