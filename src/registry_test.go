package main

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"testing"
)

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

const url = "https://example.com"

func TestRegistryImage(t *testing.T) {
	t.Run("GetRegistryImage from RegistryResponse should return RegistryImage", func(t *testing.T) {
		rResponse := RegistryResponse{
			[]byte("{\"name\":\"test\", \"tags\":[\"master-6.0.1\",\"master-6.1.0\"]}"),
			http.Header{},
		}
		gotRegistryImage := rResponse.getRegistryImage()
		expectedRegistryImage := RegistryImage{
			"test",
			[]string{"master-6.0.1", "master-6.1.0"},
		}

		equals(t, expectedRegistryImage, gotRegistryImage)
	})

	t.Run("GetRegistryImage from RegistryResponse should return empty RegistryImage", func(t *testing.T) {
		rResponse := RegistryResponse{
			[]byte("{}"),
			http.Header{},
		}
		gotRegistryImage := rResponse.getRegistryImage()
		expectedRegistryImage := RegistryImage{
			"",
			[]string(nil),
		}

		equals(t, expectedRegistryImage, gotRegistryImage)
	})
}

func TestRegistry(t *testing.T) {
	t.Run("Configure Registry secure (default) configuration", func(t *testing.T) {
		client := &http.Client{}
		expectedRegistry := Registry{
			client,
			url,
		}
		gotRegistry := Registry{}
		gotRegistry.configure(url, false)
		equals(t, expectedRegistry, gotRegistry)
	})

	t.Run("Configure Registry insecure configuration", func(t *testing.T) {
		client := &http.Client{}
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
		client.Transport = transport
		expectedRegistry := Registry{
			client,
			url,
		}
		gotRegistry := Registry{}
		gotRegistry.configure(url, true)
		equals(t, expectedRegistry, gotRegistry)
	})

	t.Run("Get API with roundtripper should return OK", func(t *testing.T) {
		client := NewTestClient(func(req *http.Request) *http.Response {
			equals(t, req.URL.String(), url+"/some/path")

			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`OK`)),
				Header:     make(http.Header),
			}
		})

		api := Registry{client, url}
		body := api.getApi("/some/path")
		equals(t, []byte("OK"), body.Body)
	})
}
