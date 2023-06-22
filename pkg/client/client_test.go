package client_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/cksidharthan/go-jaeger/pkg/client"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Parallel()

	testLogger := logrus.StandardLogger()

	testClient, err := client.New(&client.Opts{
		CollectorURL: "http://jaeger:8080/api/traces",
		ServiceName:  "test_service",
		Logger:       testLogger,
	})
	assert.Nil(t, err)
	assert.NotNil(t, testClient)
	assert.Equal(t, "http://jaeger:8080/api/traces", testClient.CollectorURL)
	assert.Equal(t, "test_service", testClient.ServiceName)
}

func TestJaegerClient_Disconnect(t *testing.T) {
	t.Parallel()

	testLogger := logrus.StandardLogger()

	testClient, err := client.New(&client.Opts{
		CollectorURL: "http://jaeger:8080/api/traces",
		ServiceName:  "test_service",
		Logger:       testLogger,
	})
	assert.Nil(t, err)
	assert.NotNil(t, testClient)

	err = testClient.Disconnect(context.Background())
	assert.Nil(t, err)
}

func TestNewTest(t *testing.T) {
	t.Parallel()

	jaegerCollectorURL := fmt.Sprintf("http://localhost:14268/api/traces")
	testLogger := logrus.StandardLogger()

	testClient, err := client.New(&client.Opts{
		CollectorURL: jaegerCollectorURL,
		ServiceName:  "test_service",
		Logger:       testLogger,
	})
	assert.Nil(t, err)
	assert.NotNil(t, testClient)
	assert.Equal(t, jaegerCollectorURL, testClient.CollectorURL)
	assert.Equal(t, "test_service", testClient.ServiceName)

	err = testClient.Disconnect(context.Background())
	assert.Nil(t, err)
}
