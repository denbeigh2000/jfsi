package client

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/denbeigh2000/jfsi"
	"github.com/denbeigh2000/jfsi/storage"
)

func NewClient(host string) storage.Store {
	return client{
		Host:   host,
		client: http.DefaultClient,
	}
}

type client struct {
	Host string

	client *http.Client
}

func (c client) url(id jfsi.ID) string {
	return fmt.Sprintf("%v/id/%v", c.Host, string(id))
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
	case 500:
		errBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		return fmt.Errorf(string(errBytes))
	}

	return nil
}

func (c client) Retrieve(jfsi.ID) (io.Reader, error) {

}

func (c client) Update(jfsi.ID, io.Reader) error {

}

func (c client) Delete(jfsi.ID) error {

}
