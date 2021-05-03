package registry

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
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

func (r Registry) ListRepositories() ([]byte, error) {
	response, err := r.Client.Get(r.BaseUrl + "/v2/_catalog?n=5000")
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	return body, err
}

func (r Registry) ListImageTags(image string) (Response, error) {
	response, err := r.Client.Get(r.BaseUrl + "/v2/" + image + "/tags/list")
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	registryResponse := Response{
		body,
		response.Header,
	}
	return registryResponse, err
}

func (r *Registry) GetExistingManifest(image string, tag string) (*http.Response, string, error) {
	request, _ := http.NewRequest("HEAD", r.BaseUrl+"/v2/"+image+"/manifests/"+tag, nil)
	request.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")
	response, err := r.Client.Do(request)
	digest := response.Header.Get("Docker-Content-Digest")
	defer response.Body.Close()
	return response, digest, err
}

func (r *Registry) DeleteImage(image, tag, digest string) error {
	request, _ := http.NewRequest("DELETE", r.BaseUrl+"/v2/"+image+"/manifests/"+digest, nil)
	res, err := r.Client.Do(request)
	defer res.Body.Close()
	if res.StatusCode != 202 {
		return errors.New("Error while deleting image:tag : " + image + ":" + tag + " HTTP code " + strconv.Itoa(res.StatusCode))
	}
	return err
}
