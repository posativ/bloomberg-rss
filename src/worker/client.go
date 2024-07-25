package worker

import (
	"fmt"
	"golang.org/x/net/proxy"
	"net/http"
	"strings"
)

type Client struct {
	httpClient *http.Client
	headers    map[string]string
}

func newClient(socksProxy string, cookie string) *Client {
	var transport = &http.Transport{}
	if socksProxy != "" {
		if strings.HasPrefix(socksProxy, "socks5://") {
			socksProxy = strings.TrimPrefix(socksProxy, "socks5://")
		}

		var auth = proxy.Auth{}
		if strings.Contains(socksProxy, "@") {
			// socks5://user:password@host:port
			parts := strings.Split(socksProxy, "@")
			authParts := strings.Split(parts[0], ":")
			auth = proxy.Auth{
				User:     authParts[0],
				Password: authParts[1],
			}
			socksProxy = parts[1]
		}

		dialSocksProxy, err := proxy.SOCKS5("tcp", socksProxy, &auth, proxy.Direct)
		if err != nil {
			fmt.Println("Error connecting to proxy:", err)
		}
		transport = &http.Transport{Dial: dialSocksProxy.Dial}
	}

	return &Client{
		httpClient: &http.Client{
			Transport: transport,
		},
		headers: map[string]string{
			"Cookie":          cookie,
			"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:125.0) Gecko/20100101 Firefox/125.0",
			"Accept-Language": "en",
		},
	}
}

func (c *Client) get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
