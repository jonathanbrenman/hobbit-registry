package clients

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type HttpClient interface {
	CheckConnectivity()
	CheckImage(image string) bool
}

type httpClient struct{
	BaseUrl string
	Client *http.Client
	RegistryCatalog struct{
		Repositories []string `json:"repositories"`
	}
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
	resp, err := hc.Client.Get(fmt.Sprintf("%s/v2/_catalog",hc.BaseUrl))
	if err != nil {
		log.Fatal("No connectivity with the private registry, please check it and try again!")
	}
	defer resp.Body.Close()
	if resp == nil || resp.StatusCode != 200 && resp.StatusCode != 301 {
		log.Fatal("No connectivity with the private registry, please check it and try again!")
	}
	json.NewDecoder(resp.Body).Decode(&hc.RegistryCatalog)
}

func (hc *httpClient) CheckImage(image string) bool {
	for _, img := range hc.RegistryCatalog.Repositories {
		if img == image {
			return true
		}
	}
	return false
}