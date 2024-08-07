package clients_test

import (
	"context"
	"sync"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/clients"
)

// check that MockResourceClient implements the clients.Requester interface
var _ clients.Requester = &MockResourceClient{}

// MockResourceClient is a mock implementation of the clients.Requester interface
type MockResourceClient struct {
	t            *testing.T
	response     interface{}
	err          error
	retries      int
	requestCount int
	retryErr     error
	mu           sync.Mutex
}

func NewMockResourceClient(t *testing.T, resp interface{}, err error, retries int, retryErr error) *MockResourceClient {
	return &MockResourceClient{
		t:            t,
		response:     resp,
		err:          err,
		retryErr:     retryErr,
		retries:      retries,
		requestCount: 0,
		mu:           sync.Mutex{},
	}
}

func (m *MockResourceClient) Get(ctx context.Context, resourceID string, apiVersion string) (interface{}, error) {
	return m.respond()
}

func (m *MockResourceClient) CreateOrUpdate(ctx context.Context, resourceID string, apiVersion string, body interface{}) (interface{}, error) {
	return m.respond()
}

func (m *MockResourceClient) Delete(ctx context.Context, resourceID string, apiVersion string) (interface{}, error) {
	return m.respond()
}

func (m *MockResourceClient) List(ctx context.Context, resourceID string, apiVersion string) (interface{}, error) {
	return m.respond()
}

func (m *MockResourceClient) Action(ctx context.Context, resourceID string, action string, apiVersion string, method string, body interface{}) (interface{}, error) {
	return m.respond()
}

func (m *MockResourceClient) respond() (interface{}, error) {
	m.t.Logf("request: %d", m.requestCount)
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.requestCount < m.retries && m.retryErr != nil {
		m.requestCount++
		return nil, m.retryErr
	}
	return m.response, m.err
}

func (m *MockResourceClient) RequestCount() int {
	return m.requestCount
}
