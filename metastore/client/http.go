package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/denbeigh2000/jfsi"
	"github.com/denbeigh2000/jfsi/metastore"
	"github.com/denbeigh2000/jfsi/utils"
)

func NewClient(host string, port int) metastore.MetaStore {
	return &client{
		Host:   fmt.Sprintf("http://%v:%v", host, port),
		client: http.DefaultClient,
	}
}

type client struct {
	Host string

	client *http.Client
}

func (c client) url(id jfsi.ID) string {
	return fmt.Sprintf("%v/%v", c.Host, string(id))
}

func (c client) Create(key jfsi.ID, n int) (r metastore.Record, err error) {
	url := c.url(key)
	metaReq := metastore.CreateRequest{NChunks: n}
	body, err := json.Marshal(metaReq)
	if err != nil {
		return
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	switch resp.StatusCode {
	case 200:
		decoder.Decode(&r)
		if err != nil {
			err = fmt.Errorf("Error decoding JSON response: %v", err.Error())
		}
	case 400:
		var errResp utils.StringResponse
		err = decoder.Decode(&errResp)
		if err != nil {
			break
		}
		if errResp.Error == metastore.ZeroLenStr {
			err = metastore.ZeroLengthCapacityRecordErr{}
			break
		}

		err = metastore.KeyAlreadyExistsErr(key)
	default:
		var errResp utils.StringResponse
		err = decoder.Decode(&errResp)
		if err != nil {
			break
		}

		err = fmt.Errorf(errResp.Error)
	}

	return
}

func (c client) Retrieve(key jfsi.ID) (r metastore.Record, err error) {
	url := c.url(key)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	switch resp.StatusCode {
	case 200:
		decoder.Decode(&r)
		if err != nil {
			err = fmt.Errorf("Error decoding JSON response: %v", err.Error())
		}
	case 404:
		err = metastore.KeyNotFoundErr(key)
	default:
		var errResp utils.StringResponse
		err = decoder.Decode(&errResp)
		if err != nil {
			break
		}

		err = fmt.Errorf(errResp.Error)
	}

	return
}

func (c client) Update(key jfsi.ID, r metastore.Record) error {
	url := c.url(key)
	body, err := json.Marshal(r)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	switch resp.StatusCode {
	case 200:
		return nil
	case 404:
		return metastore.KeyNotFoundErr(key)
	default:
		var errResp utils.StringResponse
		err = json.NewDecoder(resp.Body).Decode(&errResp)
		if err != nil {
			return err
		}

		return fmt.Errorf(errResp.Error)
	}
}

func (c client) Delete(key jfsi.ID) error {
	url := c.url(key)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	switch resp.StatusCode {
	case 200:
		return nil
	case 404:
		return metastore.KeyNotFoundErr(key)
	default:
		var errResp utils.StringResponse
		err = json.NewDecoder(resp.Body).Decode(&errResp)
		if err != nil {
			return err
		}

		return fmt.Errorf(errResp.Error)
	}
}
