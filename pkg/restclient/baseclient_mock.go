package restclient

import "net/http"

type MockClient struct {
	DoFuncMock func(req *http.Request) (*http.Response, error)
}

var (
	GetDoFuncMock func(req *http.Request) (*http.Response, error)
)

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return GetDoFuncMock(req)
}
