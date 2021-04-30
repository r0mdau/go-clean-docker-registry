package registry

import (
	"bytes"
	"crypto/tls"
	"github.com/stretchr/testify/require"
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
	t.Run("GetImage from Response should return Image", func(t *testing.T) {
		rResponse := Response{
			[]byte("{\"name\":\"test\", \"tags\":[\"master-6.0.1\",\"master-6.1.0\"]}"),
			http.Header{},
		}
		actualImage := rResponse.GetImage()
		expectedImage := Image{
			"test",
			[]string{"master-6.0.1", "master-6.1.0"},
		}

		require.Equal(t, expectedImage, actualImage)
	})

	t.Run("GetImage from Response should return empty Image", func(t *testing.T) {
		rResponse := Response{
			[]byte("{}"),
			http.Header{},
		}
		actualImage := rResponse.GetImage()
		expectedImage := Image{
			"",
			[]string(nil),
		}

		require.Equal(t, expectedImage, actualImage)
	})
}

func TestRegistry(t *testing.T) {
	t.Run("Configure Registry secure (default) configuration", func(t *testing.T) {
		client := &http.Client{}
		expectedRegistry := Registry{
			client,
			url,
		}
		actualRegistry := Registry{}
		actualRegistry.Configure(url, false)
		require.Equal(t, expectedRegistry, actualRegistry)
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
		actualRegistry := Registry{}
		actualRegistry.Configure(url, true)
		require.Equal(t, expectedRegistry, actualRegistry)
	})

	t.Run("GetCatalog API with roundtripper should return OK", func(t *testing.T) {
		client := NewTestClient(func(req *http.Request) *http.Response {
			require.Equal(t, req.URL.String(), url+"/v2/_catalog?n=5000")

			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`OK`)),
				Header:     make(http.Header),
			}
		})

		api := Registry{client, url}
		body, err := api.GetCatalog()
		require.Equal(t, []byte("OK"), body)
		require.Nil(t, err)
	})

	t.Run("GetTagsList API with roundtripper should return OK", func(t *testing.T) {
		client := NewTestClient(func(req *http.Request) *http.Response {
			require.Equal(t, req.URL.String(), url+"/v2/image/tags/list")

			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`OK`)),
				Header:     make(http.Header),
			}
		})

		api := Registry{client, url}
		body, err := api.GetTagsList("image")
		require.Equal(t, []byte("OK"), body.Body)
		require.Nil(t, err)
	})

	/*	t.Run("DeleteImageTag API with roundtripper should return OK", func(t *testing.T) {
		client := NewTestClient(func(req *http.Request) *http.Response {
			require.Equal(t, req.URL.String(), url+"/v2/image/manifests/1.0.0")

			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`OK`)),
				Header:     make(http.Header),
			}
		})

		api := Registry{client, url}
		err := api.DeleteImageTag("image", "1.0.0")
		require.Nil(t, err)
	})*/
}
