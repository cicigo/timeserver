package utils

import (
	"fmt"
	log "github.com/cihub/seelog"
	"io/ioutil"
	"net/http"
)

type AuthClient struct {
	serverEndpoint string
	client         *http.Client
}

func (c *AuthClient) Get(uuid string) (string, error) {
	url := fmt.Sprintf("%s/get?cookie=%s", c.serverEndpoint, uuid)
	r, err := c.client.Get(url)
	defer r.Body.Close()
	if err != nil {
		log.Errorf("Get request failed: %s", err)
		return "", err

	}
	name, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Read response failed: %s", err)
		return "", err
	}
	log.Info(name)
	return string(name[:]), nil
}

func (c *AuthClient) Set(uuid string, name string) error {
	url := fmt.Sprintf("%s/set?cookie=%s&name=%s", c.serverEndpoint, uuid, name)
	log.Infof("Set cookie request: %s", url)
	r, err := c.client.PostForm(url, nil)
	defer r.Body.Close()

	if err != nil {
		log.Errorf("Set request failed: %s", err)
		return err
	}
	return nil
}

func (c *AuthClient) Delete(uuid string) error {
	url := fmt.Sprintf("%s/set?cookie=%s", c.serverEndpoint, uuid)

	r, err := c.client.PostForm(url, nil)
	defer r.Body.Close()

	if err != nil {
		log.Errorf("Set request failed: %s", err)
		return err
	}
	return nil
}

func NewAuthClient(serverEndpoint string) *AuthClient {
	return &AuthClient{
		serverEndpoint: serverEndpoint,
		client:         &http.Client{},
	}
}
