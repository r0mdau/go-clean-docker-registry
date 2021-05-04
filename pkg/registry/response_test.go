package registry

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestRegistryRepository(t *testing.T) {
	t.Run("GetRepository from Response should return Repository", func(t *testing.T) {
		buffer := bytes.Buffer{}
		rResponse := Response{
			[]byte("{\"repositories\":[\"r0mdau/test\",\"r0mdau/nodejs\"]}"),
			http.Header{},
			200,
			&buffer,
		}
		actual := rResponse.GetRepository()
		expected := Repository{
			[]string{"r0mdau/test", "r0mdau/nodejs"},
		}

		require.Equal(t, expected, actual)
	})

	t.Run("GetRepository from Response should return empty Repository", func(t *testing.T) {
		buffer := bytes.Buffer{}
		rResponse := Response{
			[]byte("{}"),
			http.Header{},
			200,
			&buffer,
		}
		actual := rResponse.GetRepository()
		expected := Repository{
			[]string(nil),
		}

		require.Equal(t, expected, actual)
	})

	t.Run("GetRepository from Response should log error if response status is not 2xx", func(t *testing.T) {
		buffer := bytes.Buffer{}
		rResponse := Response{
			[]byte("{\"repositories\":[\"r0mdau/test\",\"r0mdau/nodejs\"]}"),
			http.Header{},
			404,
			&buffer,
		}
		actual := rResponse.GetRepository()
		expected := Repository{
			[]string{"r0mdau/test", "r0mdau/nodejs"},
		}

		require.Equal(t, expected, actual)
		require.Equal(t, "Non 2xx HTTP code response from registry api, got : 404, Body : {\"repositories\":[\"r0mdau/test\",\"r0mdau/nodejs\"]}", buffer.String())
	})

	t.Run("GetRepository from Response should log error if response status is not deserializable", func(t *testing.T) {
		buffer := bytes.Buffer{}
		rResponse := Response{
			[]byte("\"repositories\":[\"r0mdau/test\",\"r0mdau/nodejs\"]}"),
			http.Header{},
			200,
			&buffer,
		}
		actual := rResponse.GetRepository()
		expected := Repository{
			[]string(nil),
		}

		require.Equal(t, expected, actual)
		require.Equal(t, "Can't deserislize tree, error: invalid character ':' after top-level value\n", buffer.String())
	})
}

func TestRegistryImage(t *testing.T) {
	t.Run("GetImage from Response should return Image", func(t *testing.T) {
		buffer := bytes.Buffer{}
		rResponse := Response{
			[]byte("{\"name\":\"test\", \"tags\":[\"master-6.0.1\",\"master-6.1.0\"]}"),
			http.Header{},
			200,
			&buffer,
		}
		actual := rResponse.GetImage()
		expected := Image{
			"test",
			[]string{"master-6.0.1", "master-6.1.0"},
		}

		require.Equal(t, expected, actual)
	})

	t.Run("GetImage from Response should return empty Image", func(t *testing.T) {
		buffer := bytes.Buffer{}
		rResponse := Response{
			[]byte("{}"),
			http.Header{},
			200,
			&buffer,
		}
		actual := rResponse.GetImage()
		expected := Image{
			"",
			[]string(nil),
		}

		require.Equal(t, expected, actual)
	})

	t.Run("GetImage from Response should log error if response status is not 2xx", func(t *testing.T) {
		buffer := bytes.Buffer{}
		rResponse := Response{
			[]byte("{\"name\":\"test\", \"tags\":[\"master-6.0.1\",\"master-6.1.0\"]}"),
			http.Header{},
			404,
			&buffer,
		}
		actual := rResponse.GetImage()
		expected := Image{
			"test",
			[]string{"master-6.0.1", "master-6.1.0"},
		}

		require.Equal(t, expected, actual)
		require.Equal(t, "Non 2xx HTTP code response from registry api, got : 404, Body : {\"name\":\"test\", \"tags\":[\"master-6.0.1\",\"master-6.1.0\"]}", buffer.String())
	})

	t.Run("GetImage from Response should log error if response status is not deserializable", func(t *testing.T) {
		buffer := bytes.Buffer{}
		rResponse := Response{
			[]byte("{\"name\":\"test\", \"tags\":[\"master-6.0.1\",\"master-6.1.0\"}"),
			http.Header{},
			200,
			&buffer,
		}
		actual := rResponse.GetRepository()
		expected := Repository{
			[]string(nil),
		}

		require.Equal(t, expected, actual)
		require.Equal(t, "Can't deserislize tree, error: invalid character '}' after array element\n", buffer.String())
	})
}
