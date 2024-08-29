package testutils

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

type call struct {
	Method string
	Path   string
	Body   map[string]interface{}
}

type response struct {
	body   string
	status int
}

type requestHandlerMock struct {
	Calls     []call
	Responses []response
	t         *testing.T
}

func (mock *requestHandlerMock) AddResponse(status int, body string) *requestHandlerMock {
	mock.Responses = append(mock.Responses, response{
		body:   body,
		status: status,
	})

	return mock
}

func (mock *requestHandlerMock) AssertCalledTimes(times int) {
	calls := len(mock.Calls)

	assert.Equal(mock.t, times, calls, "Request handler mock not called the expected number of times")
}

func (mock *requestHandlerMock) AssertCalledWithMethod(method string) {
	methods := ""

	for _, call := range mock.Calls {
		if call.Method == method {
			return
		}

		methods += call.Method
	}

	assert.Fail(mock.t, fmt.Sprintf("Handler not called with method:\nexpected: %s, actual: %s", method, methods))
}

func (mock *requestHandlerMock) AssertCalledWithBody(body string) {
	unmarshalled, err := UmarshalJson([]byte(body))
	assert.Nil(mock.t, err)

	bodies := ""

	for _, call := range mock.Calls {
		if cmp.Equal(call.Body, unmarshalled) {
			return
		}
		data, err := MarshalJson(call.Body)
		assert.Nil(mock.t, err)

		bodies += data + "\n"
	}

	assert.Fail(mock.t, fmt.Sprintf("Handler not called with body:\nexpected: %s\nactual: %s", body, bodies))
}
