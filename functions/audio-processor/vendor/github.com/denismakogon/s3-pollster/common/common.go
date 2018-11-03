package common

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

type PreSignedURLs struct {
	GetURL string `json:"get_url"`
	PutURL string `json:"put_url"`
}

type RequestPayload struct {
	S3Endpoint    string        `json:"s3_endpoint"`
	Bucket        string        `json:"bucket"`
	Object        string        `json:"object"`
	PreSignedURLs PreSignedURLs `json:"presigned_urls"`
}

func WithDefault(key, defaultValue string) string {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}
	return v
}

func SetupHTTPClient() *http.Client {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 120 * time.Second,
		}).DialContext,
		MaxIdleConnsPerHost: 512,
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig: &tls.Config{
			ClientSessionCache: tls.NewLRUClientSessionCache(4096),
		},
	}
	return &http.Client{Transport: transport}
}

func DoRequest(req *http.Request, httpClient *http.Client) error {
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode > http.StatusAccepted {
		return fmt.Errorf("unable to submit webhoot successfully, "+
			"status code: %d, response body: '%s'", resp.StatusCode, string(b))
	}
	log.Printf("request accepted, response: %s", string(b))

	return nil
}
