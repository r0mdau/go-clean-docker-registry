package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type RegistryImage struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

type RegistryResponse struct {
	Body   []byte
	Header http.Header
}

func (r RegistryResponse) getRegistryImage() RegistryImage {
	var registryImage RegistryImage
	err := json.Unmarshal(r.Body, &registryImage)
	if err != nil {
		fmt.Printf("- Can't deserislize tree, error: %v\n", err)
	}
	return registryImage
}

type Registry struct {
	Client  *http.Client
	BaseUrl string
}

func (r *Registry) configure(url string, insecure bool) {
	client := &http.Client{}
	if insecure {
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
		client.Transport = transport
	}
	r.Client = client
	r.BaseUrl = url
}

func (r Registry) getTagsList(path string) RegistryResponse {
	resp, err := r.Client.Get(r.BaseUrl + path)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	registryResponse := RegistryResponse{
		body,
		resp.Header,
	}
	return registryResponse
}

func (r Registry) getCatalog(path string) []byte {
	resp, err := r.Client.Get(r.BaseUrl + path)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body
}
