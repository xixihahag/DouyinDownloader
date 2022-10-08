package httputil

import (
	"douyin/pkg/util/logutil"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	DefaultUserAgent = `Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1`
)

func Get(url string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		msg := fmt.Sprintf("NewRequest fail, url: %s, err: %v", url, err)
		logutil.Error(msg)
		return "", errors.New(msg)
	}
	req.Header.Add("User-Agent", DefaultUserAgent)
	req.Header.Add("Upgrade-Insecure-Requests", "1")

	rsp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		logutil.Errorf("http Do fail, err: %v", err)
		return "", err
	}
	defer rsp.Body.Close()
	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		logutil.Errorf("ReadAll fail, err: %v", err)
		return "", err
	}
	return string(body), nil
}
