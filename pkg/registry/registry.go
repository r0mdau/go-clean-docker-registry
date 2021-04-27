package registry

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Image struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

type Response struct {
	Body   []byte
	Header http.Header
}

func (r Response) GetImage() Image {
	var registryImage Image
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

func (r *Registry) Configure(url string, insecure bool) {
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

func (r Registry) GetTagsList(path string) Response {
	resp, err := r.Client.Get(r.BaseUrl + path)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	registryResponse := Response{
		body,
		resp.Header,
	}
	return registryResponse
}

func (r Registry) GetCatalog(path string) []byte {
	resp, err := r.Client.Get(r.BaseUrl + path)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body
}
