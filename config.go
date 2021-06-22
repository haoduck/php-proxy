package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/url"
)

//CaCert/CA.crt should be trusted by local OS
//this is goagent CA,you can gen your root ca by openssl
var CaCert = []byte(`-----BEGIN CERTIFICATE-----
MIIDUjCCAjoCAQAwDQYJKoZIhvcNAQELBQAwbzELMAkGA1UEBhMCQ04xETAPBgNV
BAgMCEludGVybmV0MQ8wDQYDVQQHDAZDZXJuZXQxEDAOBgNVBAoMB0dvQWdlbnQx
FTATBgNVBAsMDEdvQWdlbnQgUm9vdDETMBEGA1UEAwwKR29BZ2VudCBDQTAeFw0y
MTA2MTcxMTMxNTNaFw0zMTA2MTcxMTMxNTNaMG8xCzAJBgNVBAYTAkNOMREwDwYD
VQQIDAhJbnRlcm5ldDEPMA0GA1UEBwwGQ2VybmV0MRAwDgYDVQQKDAdHb0FnZW50
MRUwEwYDVQQLDAxHb0FnZW50IFJvb3QxEzARBgNVBAMMCkdvQWdlbnQgQ0EwggEi
MA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDQk/hAE0yDjfXWZbX3VOEiroqM
S97RVmVs2GLcWP9GVFJzdIvJVlejsjTMq1XG1boybww3upHCCqaA88N9CLSUQDuL
j32l1/zbYZpf7SlMdcQ/ekOK8+tFLUj4VIAX08Ud6v+xj8v8PcwCM21wf+WCIgaD
TT30SXuPSWfebpJu6fAHlhA+CLTSR35oBrwHu6e5DElw6GgK92X5S6wlFeIWywzc
Mo3bKn6AX31wI9spfnnzMHeOZEScyEv8AG3O2LTrl14EdvSAv6S3Mp7qW8G/nKJY
3rUcFOyKz82yZtkQ6vVlKBrNmSPmVZxncMpjuMFPzHSU1jeuEZ8EPOSTPol7AgMB
AAEwDQYJKoZIhvcNAQELBQADggEBAK4dWpx62OvYEJc6eiurUDcj6xZMew7uCPTp
s8dcAYQcfmMBy7RrDCBJ0DhdBa6+BHe2yqDJpIiMDM1rhRKmVGSm0BbH7/rkZAFA
Ie34g2IylKMd86H2rAyi3dDBhMZHJrTOzhIJpnmAJ+msMZ3NSOomik0VsXYZQuOD
oNSbr4+NMxEeJfKNNFqtJ6Jnf+5Z5NLpd3MZlEkkZalAvQBZqSb2FezCIeaOwZp2
uUEF1MRCft2inuPTeSCH1ssdRrSM4F79+RdI13rejkx/YNym6JfIBZ8PbmyFWBJl
hzlF2tPatVqnOvUFVBWh+EVUVkgTyfsn0dH6CquH4buj8RpF9qc=
-----END CERTIFICATE-----`)

var CaKey = []byte(`-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDQk/hAE0yDjfXW
ZbX3VOEiroqMS97RVmVs2GLcWP9GVFJzdIvJVlejsjTMq1XG1boybww3upHCCqaA
88N9CLSUQDuLj32l1/zbYZpf7SlMdcQ/ekOK8+tFLUj4VIAX08Ud6v+xj8v8PcwC
M21wf+WCIgaDTT30SXuPSWfebpJu6fAHlhA+CLTSR35oBrwHu6e5DElw6GgK92X5
S6wlFeIWywzcMo3bKn6AX31wI9spfnnzMHeOZEScyEv8AG3O2LTrl14EdvSAv6S3
Mp7qW8G/nKJY3rUcFOyKz82yZtkQ6vVlKBrNmSPmVZxncMpjuMFPzHSU1jeuEZ8E
POSTPol7AgMBAAECggEBAInZEVu/pXTYcJ4mkHGK7lQOiNCaIAO4FsYt6IB6bRPd
DLTzVKNW5grw7wZJQiJsBGfqjmeSbVyRz2MwN4W+KCJGpVPiHIdrzNhslCtLwVyg
BHhzZIpEFLyeZjiBGDsnIYJZWm3OfGETsm3N8UlFrbgopqdGeGin2/ph6DzQVQTp
c7DohIH9fmPJZVW7Zfq/AgqF5KvQNBBlm+c9sHukSdyY+HKTC03812+9uJKw0qEL
IJ3N/6TEgoM3ZKp0XDelAXt7Yaooit4C9OzfI7EhXOrhxVfE7cihj6t6jI9kAjqm
D8FSHdGolz6JQqJq9Jxw+fGAgfYaf/cOF0EJnKaOfzkCgYEA73M2DpTcln4ZVj4v
VpXtlfY2XHVUP2mScGPVu+G3VJsiZjvN1d+rCw4CEntCo4XGxs7fYgcCLHTb2eIg
6Q9s27vbcjqwxOCyulA3YD7AFFd51m57l2wzPTBNsNUvuQz39dcYvyS7thCF5Qio
E553phDXw858ENOIWrm7V0uGDF8CgYEA3v6DkTqcdNm4/jXPyQj4mjN3QzZV3nlk
7xPJmM3Vzh7dj6S3jHUiZmYeFFMOBc8mLJ68kDIu00htAsdBqGlVgxFukeYMcY8f
uR2sbOgExuizxGHoMVeJKwYeCVkedP+FQx19cg/Gi1bpJGD9toIUEoNqAS7evytC
0W6kuJLFWGUCgYBahevBx3U9T560g/3RdgzDzLjwa0rWTksWQifjR4nPanauv50p
Zc17+GfAJOkkeMaVElBQ9uVTeTpEPMDEWxiEWZi0rot1Yp0u4nSM5iwnhIqDDnGa
5UTZtREp8O6Bvu1e+1pXqMNuKQD1fThNcnM5TNTFKaKtmcrKwbyZW+vpcQKBgQCm
ZxS05iDkjagfgvZoVVp2b1ta+4v+dWYdhg2VCly28I9zZn5VwP8HnMJrdkLrkNYy
y814aQpKPiyiuyBC1T+ri/GPzDSS9TO+BuepaUZPTE0BifIkB+djBLCbVzaEJj1C
hRocaKtHRXa63+nULKNf4VLUSS6NR3IYKNGgrl23hQKBgEyEcyRvim0UBjCLvhH8
brl1Ltf3/MiRXxPQTOT7A5O7F/fgOvGHPoZBRY/qdFK8q/FvSvPi1rR2K3PIigd4
OR5JBs8MGf8pFqojIW+10cpgZm8PApJGZ2Tf+TdiLFsW2IUFVpJAF+oAPiwQ2Zks
vn/ZWUFQVg9bQ5WCJ3JHlPAn
-----END PRIVATE KEY-----`)

type config struct {
	//php fetchserver path
	fetchserver string
	//password
	password string
	//when connect https php server,TLS sni extension
	sni string
	//local listen address
	listen string
	//fetchserver info parsed by url.Parse
	server_url url.URL
	//ca sign ssl cert for middle intercept
	signer *CaSigner
	//root ca info
	Ca tls.Certificate
}

func (c *config) init_config() {
	//
	flag.StringVar(&c.listen, "l", "127.0.0.1:8081", "Local listen address")
	flag.StringVar(&c.password, "p", "123456", "php server password")
	flag.StringVar(&c.sni, "sni", "", "HTTPS sni extension ServerName(default fetchserver hostname)")
	flag.StringVar(&c.fetchserver, "s", "https://a.bc.com/php-proxy/index.php", "php fetchserver path(http/https)")
	flag.Parse()
	//
	server_url, err := url.Parse(c.fetchserver)
	if err != nil {
		log.Fatal(err)
	}
	c.server_url = *server_url
	//
	c.signer = NewCaSignerCache(1024)
	ca, err := tls.X509KeyPair(CaCert, CaKey)
	c.Ca = ca
	if err != nil {
		log.Println(err)
	} else {
		c.signer.Ca = &ca
	}
	log.Printf("php Fetch server:%s\n", c.fetchserver)
}