package main

import (
	"crypto/tls"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	MaxIdleConnections int = 20
	RequestTimeout     int = 15
)

func createHTTPClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: MaxIdleConnections,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: time.Duration(RequestTimeout) * time.Second,
	}

	return client
}

func main() {
	httpClient := createHTTPClient()

	var endPoint string = "https://localhost/proxy/"
	var done sync.WaitGroup

	for i := 'a'; i <= 'z'; i++ {
		done.Add(1)
		go func(url string) {
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				log.Fatalf("couldn't connect: %+v", err)
			}

			response, err := httpClient.Do(req)
			if err != nil && response == nil {
				log.Fatalf("couldn't send: %+v", err)
			}
			defer response.Body.Close()
			defer done.Done()

			io.Copy(ioutil.Discard, response.Body)
		}(endPoint + "nettest" + string(rune(i)) + ".fau.xxx")
	}
	done.Wait()
}
