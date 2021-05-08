package registry

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Repository struct {
	List []string `json:"repositories"`
}

type Image struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

type Response struct {
	Body       []byte
	Header     http.Header
	StatusCode int
	Writer     io.Writer
}

func NewResponse(body []byte, header http.Header, status int) Response {
	return Response{
		Body:       body,
		Header:     header,
		StatusCode: status,
		Writer:     os.Stderr,
	}
}

func (r Response) GetRepository() Repository {
	var repository Repository
	err := json.Unmarshal(r.Body, &repository)
	r.logError(err)
	return repository
}

func (r Response) GetImage() Image {
	var registryImage Image
	err := json.Unmarshal(r.Body, &registryImage)
	r.logError(err)
	return registryImage
}

func (r Response) logError(err error) {
	if r.StatusCode < 200 || r.StatusCode > 299 {
		fmt.Fprintf(r.Writer, "Non 2xx HTTP code response from registry api, got : %d, Body : %s", r.StatusCode, string(r.Body))
	} else if err != nil {
		fmt.Fprintf(r.Writer, "Can't deserislize tree, error: %v\n", err)
	}
}
