package requester

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestInspector(t *testing.T) {

	var dumpedReqBody []byte

	var doer DoerFunc = func(req *http.Request) (*http.Response, error) {
		dumpedReqBody, _ = ioutil.ReadAll(req.Body)
		resp := &http.Response{
			StatusCode: 201,
			Body:       ioutil.NopCloser(strings.NewReader("pong")),
		}
		return resp, nil
	}

	i := Inspector{}

	resp, body, err := Receive(&i, doer, Body("ping"))
	require.NoError(t, err)

	assert.Equal(t, 201, resp.StatusCode)
	assert.Equal(t, "pong", string(body))

	require.NotNil(t, i.Request)

	assert.Equal(t, "ping", i.RequestBody.String())
	assert.Equal(t, "ping", string(dumpedReqBody))

	require.NotNil(t, i.Response)
	assert.Equal(t, 201, i.Response.StatusCode)

	assert.Equal(t, "pong", i.ResponseBody.String())
}

func TestInspector_Clear(t *testing.T) {

	i := Inspector{
		Request:      &http.Request{},
		Response:     &http.Response{},
		RequestBody:  bytes.NewBuffer(nil),
		ResponseBody: bytes.NewBuffer(nil),
	}

	i.Clear()

	assert.Nil(t, i.Request)
	assert.Nil(t, i.Response)
	assert.Nil(t, i.RequestBody)
	assert.Nil(t, i.ResponseBody)
}

func TestInspect(t *testing.T) {

	r := MustNew()

	i := Inspect(r)

	r.Receive(MockDoer(201))

	assert.NotNil(t, i.Request)
	assert.Equal(t, 201, i.Response.StatusCode)

}

// Inspect returns an Inspector, which captures the traffic to and from a Requester.  It's
// a tool for writing tests.
func ExampleInspect() {

	r := MustNew(
		MockDoer(201, Body("pong")),
		Header(HeaderAccept, MediaTypeTextPlain),
		Body("ping"),
	)

	i := Inspect(r)

	r.Receive(nil)

	fmt.Println(i.Request.Header.Get(HeaderAccept))
	fmt.Println(i.RequestBody.String())
	fmt.Println(i.Response.StatusCode)
	fmt.Println(i.ResponseBody.String())

	// Output:
	// text/plain
	// ping
	// 201
	// pong
}
