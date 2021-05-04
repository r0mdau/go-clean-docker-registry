package registry

import (
	"crypto/tls"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
)

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

func (r Registry) VersionCheck() error {
	response, err := r.Client.Get(r.BaseUrl + "/v2/")
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return errors.New("this docker registry does not implements the V2(.1) registry API and the client cannot proceed safely with other V2 operation")
	}
	return err
}

func (r Registry) ListRepositories() (Response, error) {
	response, err := r.Client.Get(r.BaseUrl + "/v2/_catalog?n=5000")
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	rResponse := NewResponse(
		body,
		response.Header,
		response.StatusCode,
	)
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

func (r *Registry) GetDigestFromManifest(image string, tag string) (string, error) {
	request, _ := http.NewRequest("HEAD", r.BaseUrl+"/v2/"+image+"/manifests/"+tag, nil)
	request.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")
	response, err := r.Client.Do(request)
	digest := response.Header.Get("Docker-Content-Digest")
	defer response.Body.Close()
	if response.StatusCode != 200 {
		err = errors.New("Error while getting digest from manifest for: " + image + ":" + tag + ", HTTP code " + strconv.Itoa(response.StatusCode))
	}
	return digest, err
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
