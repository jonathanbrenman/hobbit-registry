package clients

import (
	"log"
	"net/http"
	"time"
)

type HttpClient interface {
	CheckConnectivity()
}

type httpClient struct{
	BaseUrl string
	Client *http.Client
}

func NewHttpClient(baseUrl string) HttpClient {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    2 * time.Second,
		DisableCompression: true,
	}
	return &httpClient{
		Client: &http.Client{Transport: tr},
		BaseUrl: baseUrl,
	}
}

func (hc *httpClient) CheckConnectivity() {
	resp, err := hc.Client.Get(hc.BaseUrl)
	if err != nil {
		log.Fatal("No connectivity with the private registry, please check it and try again!")
	}
	defer resp.Body.Close()
	if resp == nil || resp.StatusCode != 200 && resp.StatusCode != 301 {
		log.Fatal("No connectivity with the private registry, please check it and try again!")
	}
}