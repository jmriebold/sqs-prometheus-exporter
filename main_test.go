package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockSQSClient is a mock of the SQS client
type MockSQSClient struct {
	mock.Mock
}

// GetQueueAttributes mocks the GetQueueAttributes method
func (m *MockSQSClient) GetQueueAttributes(ctx context.Context, params *sqs.GetQueueAttributesInput, optFns ...func(*sqs.Options)) (*sqs.GetQueueAttributesOutput, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*sqs.GetQueueAttributesOutput), args.Error(1)
}

func TestGetMonitorInterval(t *testing.T) {
	tests := []struct {
		name           string
		envValue       string
		expectedResult time.Duration
	}{
		{"Default value", "", 30 * time.Second},
		{"Custom value", "60", 60 * time.Second},
		{"Invalid value", "invalid", 30 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NoError(t, os.Setenv("SQS_MONITOR_INTERVAL_SECONDS", tt.envValue))
			defer func() {
				require.NoError(t, os.Unsetenv("SQS_MONITOR_INTERVAL_SECONDS"))
			}()

			result := getMonitorInterval()
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestMonitorQueue(t *testing.T) {
	mockClient := new(MockSQSClient)
	originalSvc := svc
	svc = mockClient
	defer func() { svc = originalSvc }()

	queueURL := "https://sqs.us-west-2.amazonaws.com/123456789012/MyQueue"
	expectedOutput := &sqs.GetQueueAttributesOutput{
		Attributes: map[string]string{
			"ApproximateNumberOfMessages":           "10",
			"ApproximateNumberOfMessagesDelayed":    "5",
			"ApproximateNumberOfMessagesNotVisible": "2",
		},
	}

	mockClient.On("GetQueueAttributes", mock.Anything, mock.Anything, mock.Anything).Return(expectedOutput, nil)

	c := make(chan queueResult)
	go monitorQueue(queueURL, c)

	result := <-c

	assert.Equal(t, queueURL, result.QueueURL)
	assert.Equal(t, "MyQueue", result.QueueName)
	assert.Equal(t, expectedOutput, result.QueueResults)

	mockClient.AssertExpectations(t)
}

func TestHealthcheck(t *testing.T) {
	assert.NotPanics(t, func() { healthcheck(nil, nil) })
}

func TestSQSMetrics(t *testing.T) {
	// Set up mock SQS client
	mockClient := new(MockSQSClient)
	originalSvc := svc
	svc = mockClient
	defer func() { svc = originalSvc }()

	// Set up test queue URL and expected metrics
	queueURL := "https://sqs.us-west-2.amazonaws.com/123456789012/TestQueue"
	expectedOutput := &sqs.GetQueueAttributesOutput{
		Attributes: map[string]string{
			"ApproximateNumberOfMessages":           "10",
			"ApproximateNumberOfMessagesDelayed":    "5",
			"ApproximateNumberOfMessagesNotVisible": "2",
		},
	}

	// Configure mock client to return the expected output
	mockClient.On("GetQueueAttributes", mock.Anything, mock.Anything, mock.Anything).Return(expectedOutput, nil)

	// Set environment variables
	require.NoError(t, os.Setenv("SQS_QUEUE_URLS", queueURL))
	require.NoError(t, os.Setenv("SQS_MONITOR_INTERVAL_SECONDS", "1"))
	defer func() {
		require.NoError(t, os.Unsetenv("SQS_QUEUE_URLS"))
		require.NoError(t, os.Unsetenv("SQS_MONITOR_INTERVAL_SECONDS"))
	}()

	// Set the monitor interval
	originalInterval := monitorInterval
	monitorInterval = getMonitorInterval()
	defer func() { monitorInterval = originalInterval }()

	// Create a cancellable context for the monitoring goroutine
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the monitor in a goroutine
	go monitorQueues(ctx, []string{queueURL})

	// Wait for metrics to be collected
	time.Sleep(2 * time.Second)

	// Set up test HTTP server
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	// Send request to /metrics endpoint
	resp, err := http.Get(testServer.URL + "/metrics")
	assert.NoError(t, err)
	defer func() {
		require.NoError(t, resp.Body.Close())
	}()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	// Check for expected metrics
	bodyString := string(body)
	expectedMetrics := []string{
		`sqs_approximatenumberofmessages{queue="TestQueue"} 10`,
		`sqs_approximatenumberofmessagesdelayed{queue="TestQueue"} 5`,
		`sqs_approximatenumberofmessagesnotvisible{queue="TestQueue"} 2`,
	}

	for _, metric := range expectedMetrics {
		assert.True(t, strings.Contains(bodyString, metric), fmt.Sprintf("Metric not found: %s", metric))
	}

	// Clear registered metrics to avoid conflicts in other tests
	prometheus.DefaultRegisterer = prometheus.NewRegistry()
}
