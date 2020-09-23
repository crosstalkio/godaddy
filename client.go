package godaddy

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Client struct {
	client *http.Client
	url    string
	key    string
	secret string
}

func NewClient(url, key, secret string, timeout time.Duration) *Client {
	return &Client{
		client: &http.Client{Timeout: time.Second * time.Duration(timeout)},
		url:    url,
		key:    key,
		secret: secret,
	}
}

func (c *Client) PutRecord(domain, kind, name, addr string, ttl int) error {
	if domain == "" {
		return errors.New("Empty domain")
	}
	if kind == "" {
		return errors.New("Empty type")
	}
	if name == "" {
		return errors.New("Empty name")
	}
	if addr == "" {
		return errors.New("Empty data (addr)")
	}
	if ttl <= 0 {
		return fmt.Errorf("Invalid TTL: %d", ttl)
	}
	rec := &struct {
		Data string `json:"data"`
		TTL  int    `json:"ttl"`
	}{
		Data: addr,
		TTL:  ttl,
	}
	data, err := json.Marshal(&[]interface{}{rec})
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/v1/domains/%s/records/%s/%s", c.url, domain, kind, name)
	log.Println(url)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("sso-key %s:%s", c.key, c.secret))
	req.Header.Set("Content-Type", "application/json")
	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	data, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("%s: %s", res.Status, data)
	}
	return nil
}
