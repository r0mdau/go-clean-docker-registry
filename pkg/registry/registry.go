package registry

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Registry struct {
	Client  *http.Client
	BaseUrl string
}

func NewRegistry(url string, insecure bool) Registry {
	client := &http.Client{}
	if insecure {
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
		client.Transport = transport
	}
	return Registry{
		Client:  client,
		BaseUrl: url,
	}
}

func (r Registry) VersionCheck() error {
	response, err := r.Client.Get(r.BaseUrl + "/v2/")
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return r.httpErr(response, "this docker registry does not implements the V2(.1) registry API and the client cannot proceed safely with other V2 operation")
	}
	return err
}

func (r Registry) ListRepositories(n int) (Response, error) {
	response, err := r.Client.Get(r.BaseUrl + "/v2/_catalog?n=" + strconv.Itoa(n))
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	rResponse := NewResponse(
		body,
		response.Header,
		response.StatusCode,
	)
	if response.StatusCode == http.StatusGatewayTimeout {
		err = r.httpErr(response, "you should retry by specifying -n parameter to limit the number of returned elements")
	}
	return rResponse, err
}

func (r Registry) ListImageTags(image string) (Response, error) {
	response, err := r.Client.Get(r.BaseUrl + "/v2/" + image + "/tags/list")
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	rResponse := NewResponse(
		body,
		response.Header,
		response.StatusCode,
	)
	return rResponse, err
}

func (r Registry) GetDigestFromManifest(image string, tag string) (string, error) {
	request, _ := http.NewRequest("HEAD", r.BaseUrl+"/v2/"+image+"/manifests/"+tag, nil)
	request.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")
	response, err := r.Client.Do(request)
	digest := response.Header.Get("Docker-Content-Digest")
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		err = r.httpErr(response, "Error while getting digest from manifest for: "+image+":"+tag)
	}
	return digest, err
}

func (r Registry) DeleteImage(image, tag, digest string) error {
	request, _ := http.NewRequest("DELETE", r.BaseUrl+"/v2/"+image+"/manifests/"+digest, nil)
	response, err := r.Client.Do(request)
	defer response.Body.Close()
	if response.StatusCode != http.StatusAccepted {
		return r.httpErr(response, "Error while deleting image:tag : "+image+":"+tag)
	}
	return err
}

func (r Registry) httpErr(response *http.Response, message string) error {
	return errors.New(fmt.Sprintf("%d %s : %s", response.StatusCode, http.StatusText(response.StatusCode), message))
}
