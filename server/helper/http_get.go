package helper

import (
	"fmt"
	"net/http"
	"time"
)

const timeout = time.Second * 20

var client = http.Client{
	Timeout: timeout,
}

func GetWithTimeout(url string) (*http.Response, error) {
	fmt.Printf("getting %s with timeout %s\n", url, timeout)
	return client.Get(url)
}
