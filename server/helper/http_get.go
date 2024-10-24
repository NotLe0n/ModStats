package helper

import (
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

var timeout = time.Second * time.Duration(viper.GetInt32("timeout"))

var client = http.Client{
	Timeout: timeout,
}

func GetWithTimeout(url string) (*http.Response, error) {
	fmt.Printf("getting %s with timeout %s\n", url, timeout)
	return client.Get(url)
}
