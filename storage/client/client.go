package client

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/denbeigh2000/jfsi"
	"github.com/denbeigh2000/jfsi/storage"
)

func NewClient(host string, port int) storage.Store {
	return client{
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

func (c client) Create(id jfsi.ID, r io.Reader) error {
	url := c.url(id)
	req, err := http.NewRequest(http.MethodPost, url, r)
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
	case 400:
		return storage.AlreadyExistsErr(id)
	default:
		return handleUnknownError(resp.StatusCode, resp.Body)
	}
}

func (c client) Retrieve(id jfsi.ID) (io.Reader, error) {
	url := c.url(id)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("Success received, but error reading body: %v", err)
		}

		return bytes.NewReader(body), nil
	case 404:
		return nil, storage.NotFoundErr(id)
	default:
		return nil, handleUnknownError(resp.StatusCode, resp.Body)
	}
}

func (c client) Update(id jfsi.ID, r io.Reader) error {
	url := c.url(id)
	req, err := http.NewRequest(http.MethodPut, url, r)
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
		return storage.NotFoundErr(id)
	default:
		return handleUnknownError(resp.StatusCode, resp.Body)
	}
}

func (c client) Delete(id jfsi.ID) error {
	url := c.url(id)
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
		return storage.NotFoundErr(id)
	default:
		return handleUnknownError(resp.StatusCode, resp.Body)
	}
}
