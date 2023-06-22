package main

import (
	"context"
	"fmt"
	"github.com/cksidharthan/go-jaeger/pkg/client"
	"github.com/cksidharthan/go-jaeger/pkg/trace"
	"github.com/sirupsen/logrus"
	"time"
)

func main() {

	exLogger := logrus.StandardLogger()

	jaegerClient, err := client.New(&client.Opts{
		CollectorURL: "http://localhost:14268/api/traces",
		ServiceName:  "test_service",
		Environment:  "dev",
		Logger:       exLogger,
	})
	if err != nil {
		fmt.Println(err)
	}

	defer jaegerClient.Disconnect(context.Background())

	methodTrace := trace.NewTraceWithContext(context.Background())
	defer methodTrace.Close()

	firstFunc(methodTrace, "John", "Doe")
	time.Sleep(1 * time.Second)
	firstFunc(methodTrace, "Jane", "Doe")
}

func firstFunc(funcTrace *trace.Trace, firstName string, lastName string) {
	subTrace := funcTrace.StartNewSpanWithName(fmt.Sprintf("changeName(%s, %s)", firstName, lastName))
	subTrace.SetTag("process", "firstFunc")
	defer subTrace.Close()
	time.Sleep(1 * time.Second)

	changeName(subTrace, firstName, lastName)
}

func changeName(funcTrace *trace.Trace, firstName string, lastName string) {
	subTrace := funcTrace.StartNewSpanWithName(fmt.Sprintf("changeName(%s, %s)", firstName, lastName))
	defer subTrace.Close()

	time.Sleep(1 * time.Second)

	subTrace.SetTag("first_name", firstName)
	subTrace.SetTag("last_name", lastName)
	subTrace.SetTag("process", "firstFunc")
	time.Sleep(1 * time.Second)
}
