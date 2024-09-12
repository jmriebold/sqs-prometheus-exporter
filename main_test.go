package main

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
			os.Setenv("SQS_MONITOR_INTERVAL_SECONDS", tt.envValue)
			defer os.Unsetenv("SQS_MONITOR_INTERVAL_SECONDS")

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
