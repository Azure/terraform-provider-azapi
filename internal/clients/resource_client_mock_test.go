package clients_test

import (
	"context"
	"sync"
	"testing"
	"time"

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
	requestTimes []time.Time
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
		requestTimes: make([]time.Time, 0, 10),
	}
}

func (m *MockResourceClient) Get(ctx context.Context, resourceID string, apiVersion string, options clients.RequestOptions) (interface{}, error) {
	return m.respond(ctx)
}

func (m *MockResourceClient) CreateOrUpdate(ctx context.Context, resourceID string, apiVersion string, body interface{}, options clients.RequestOptions) (interface{}, error) {
	return m.respond(ctx)
}

func (m *MockResourceClient) Delete(ctx context.Context, resourceID string, apiVersion string, options clients.RequestOptions) (interface{}, error) {
	return m.respond(ctx)
}

func (m *MockResourceClient) List(ctx context.Context, resourceID string, apiVersion string, options clients.RequestOptions) (interface{}, error) {
	return m.respond(ctx)
}

func (m *MockResourceClient) Action(ctx context.Context, resourceID string, action string, apiVersion string, method string, body interface{}, options clients.RequestOptions) (interface{}, error) {
	return m.respond(ctx)
}

func (m *MockResourceClient) respond(ctx context.Context) (interface{}, error) {
	select {
	case <-ctx.Done():
		m.t.Logf("context cancelled")
		return nil, ctx.Err()
	default:
		m.mu.Lock()
		defer m.mu.Unlock()
		timeSinceLastRequest := time.Duration(0)
		if len(m.requestTimes) != 0 {
			timeSinceLastRequest = time.Since(m.requestTimes[len(m.requestTimes)-1])
		}
		m.t.Logf("request: %d, time since last: %s", m.requestCount, timeSinceLastRequest)
		m.requestTimes = append(m.requestTimes, time.Now())
		if m.requestCount < m.retries && m.retryErr != nil {
			m.requestCount++
			return nil, m.retryErr
		}
		return m.response, m.err
	}
}

func (m *MockResourceClient) RequestCount() int {
	return m.requestCount
}

func (m *MockResourceClient) RequestTimes() []time.Time {
	return m.requestTimes
}
