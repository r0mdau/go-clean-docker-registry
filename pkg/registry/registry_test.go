package registry

import (
	"bytes"
	"crypto/tls"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
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

func getHttpResponse() *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewBufferString("OK")),
		Header:     make(http.Header),
	}
}

func TestRegistry(t *testing.T) {
	t.Run("NewRegistry secure (default) configuration", func(t *testing.T) {
		client := &http.Client{
			Timeout: 30 * time.Second,
		}
		expectedRegistry := Registry{
			client,
			url,
		}
		actualRegistry := NewRegistry(url, false)
		require.Equal(t, expectedRegistry, actualRegistry)
	})

	t.Run("NewRegistry insecure configuration", func(t *testing.T) {
		client := &http.Client{
			Timeout: 30 * time.Second,
		}
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
		actualRegistry := NewRegistry(url, true)
		require.Equal(t, expectedRegistry, actualRegistry)
	})

	expectedResponse := NewResponse(
		[]byte("OK"),
		make(http.Header),
		http.StatusOK,
	)

	t.Run("API Version Check with roundtripper should return no error if v2", func(t *testing.T) {
		client := NewTestClient(func(req *http.Request) *http.Response {
			require.Equal(t, url+"/v2/", req.URL.String())
			return getHttpResponse()
		})

		api := Registry{client, url}
		err := api.VersionCheck()
		require.NoError(t, err)
	})

	t.Run("API Version Check with roundtripper should return KO if not v2", func(t *testing.T) {
		client := NewTestClient(func(req *http.Request) *http.Response {
			require.Equal(t, url+"/v2/", req.URL.String())

			return &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`OK`)),
				Header:     make(http.Header),
			}
		})

		api := Registry{client, url}
		err := api.VersionCheck()
		require.Error(t, err)
	})

	t.Run("ListRepositories API with roundtripper should return OK", func(t *testing.T) {
		client := NewTestClient(func(req *http.Request) *http.Response {
			require.Equal(t, url+"/v2/_catalog?n=5000", req.URL.String())
			return getHttpResponse()
		})

		api := Registry{client, url}
		response, err := api.ListRepositories(5000)

		require.Equal(t, expectedResponse, response)
		require.NoError(t, err)
	})

	t.Run("ListImageTags API with roundtripper should return OK", func(t *testing.T) {
		client := NewTestClient(func(req *http.Request) *http.Response {
			require.Equal(t, url+"/v2/image/tags/list", req.URL.String())
			return getHttpResponse()
		})

		api := Registry{client, url}
		actual, err := api.ListImageTags("image")

		require.Equal(t, expectedResponse, actual)
		require.NoError(t, err)
	})

	t.Run("GetDigestFromManifest API with roundtripper should return image hash", func(t *testing.T) {
		client := NewTestClient(func(req *http.Request) *http.Response {
			require.Equal(t, url+"/v2/image/manifests/tag", req.URL.String())

			header := make(http.Header)
			header.Add("Docker-Content-Digest", "sha256sum")
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`OK`)),
				Header:     header,
			}
		})

		api := Registry{client, url}
		digest, err := api.GetDigestFromManifest("image", "tag")
		require.Equal(t, "sha256sum", digest)
		require.NoError(t, err)
	})

	t.Run("DeleteImage API with roundtripper should return no error", func(t *testing.T) {
		client := NewTestClient(func(req *http.Request) *http.Response {
			require.Equal(t, url+"/v2/image/manifests/sha256sum", req.URL.String())

			return &http.Response{
				StatusCode: http.StatusAccepted,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`OK`)),
				Header:     make(http.Header),
			}
		})

		api := Registry{client, url}
		err := api.DeleteImage("image", "tag", "sha256sum")
		require.NoError(t, err)
	})

	t.Run("DeleteImage API with roundtripper should return error if StatusCode != 202", func(t *testing.T) {
		client := NewTestClient(func(req *http.Request) *http.Response {
			require.Equal(t, url+"/v2/image/manifests/sha256sum", req.URL.String())
			return getHttpResponse()
		})

		api := Registry{client, url}
		err := api.DeleteImage("image", "tag", "sha256sum")
		require.Error(t, err)
	})
}
