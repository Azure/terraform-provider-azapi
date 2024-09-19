package clients_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
)

// check that MockDataPlaneClient implements the clients.DataPlaneRequester interface
var _ clients.DataPlaneRequester = &MockDataPlaneClient{}

// MockDataPlaneClient is a mock implementation of the clients.DataPlaneRequester interface
type MockDataPlaneClient struct {
	t            *testing.T
	response     interface{}
	err          error
	retries      int
	requestCount int
	retryErr     error
	mu           sync.Mutex
	requestTimes []time.Time
}

func NewMockDataPlaneClient(t *testing.T, resp interface{}, err error, retries int, retryErr error) *MockDataPlaneClient {
	return &MockDataPlaneClient{
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

func (m *MockDataPlaneClient) CreateOrUpdateThenPoll(ctx context.Context, id parse.DataPlaneResourceId, body interface{}, options clients.RequestOptions) (interface{}, error) {
	return m.respond(ctx)
}

func (m *MockDataPlaneClient) Get(ctx context.Context, id parse.DataPlaneResourceId, options clients.RequestOptions) (interface{}, error) {
	return m.respond(ctx)
}

func (m *MockDataPlaneClient) DeleteThenPoll(ctx context.Context, id parse.DataPlaneResourceId, options clients.RequestOptions) (interface{}, error) {
	return m.respond(ctx)
}

func (m *MockDataPlaneClient) Action(ctx context.Context, resourceID string, action string, apiVersion string, method string, body interface{}, options clients.RequestOptions) (interface{}, error) {
	return m.respond(ctx)
}

func (m *MockDataPlaneClient) respond(ctx context.Context) (interface{}, error) {
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

func (m *MockDataPlaneClient) RequestCount() int {
	return m.requestCount
}

func (m *MockDataPlaneClient) RequestTimes() []time.Time {
	return m.requestTimes
}
