package main

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/nats-io/nats.go"
	"time"

	golang "unstable"
)

func main() {
	nc, _ := golang.NewNats()
	defer nc.Drain()
	uuid := uuid.Must(uuid.NewV4()).String()
	payload := fmt.Sprintf(`
		{
			"payload": {
				"CreateMonitoring": {
					"id": "%s",
					"url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
				}
			}
		}
	`, uuid)

	// 1. Region
	// 2. Project
	// 3. Environment
	// 4. Release
	// 5. Function Name

	// OP_ROUTER.east1.kociflow.development.979183c.envoy.observability.auth.v3.check
	// OP_ROUTER.[region].[project].[environment].[release id].[function name]
	// OP_ROUTER.[region].[project].[release id].[function name]

	c := NewClient(nc, opts{
		Region:  "useast1",
		Project: "kociflow",
		AppEnv:  "development",
	})

	c.Invoke(context.TODO(), invokeInput{
		FunctionName:     "envoy.observability.auth.v3.check",
		InvocationType:   Ptr(RequestResponse),
		Payload:          []byte(payload),
		Metadata:         nil,
		CorrelationID:    nil,
		CausationID:      nil,
		ExpectedRevision: nil,
		Release:          "latest",
	})

	rep, er := nc.Request("srv.command.monitoring.create-monitoring", []byte(payload), time.Second*3)
	if er != nil {
		fmt.Println(er)
		return
	}
	fmt.Println(string(rep.Data))
}

type opts struct {
	Region  string
	Project string
	AppEnv  string
}
type pepeg struct {
	nc      *nats.Conn
	region  string
	project string
	appEnv  string
}

func (p pepeg) Invoke(ctx context.Context, payload invokeInput) (*InvokeOutput, error) {
	rep, err := p.nc.Request("srv.command.monitoring.create-monitoring", []byte(payload), time.Second*3)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(rep.Data))
	return nil, nil
}

// TODO: read how to leverage https://docs.nats.io/nats-concepts/subject_mapping#weighted-mappings
//  as a way to micmic the AWS Lambda Alias Weighted Routing

// TODO: credit KociQQ and read https://www.jerrychang.ca/writing/the-ultimate-guide-to-aws-lambda-alias

func NewClient(nc *nats.Conn, opt opts) pepeg {
	return pepeg{
		nc:      nc,
		region:  opt.Region,
		project: opt.Project,
		appEnv:  opt.AppEnv,
	}
}

func Ptr[T any](v T) *T {
	return &v
}

type invocationType string

const (
	// Use Core Nats
	RequestResponse invocationType = "RequestResponse"
	// Use Jetstream
	Event invocationType = "Event"
	// No-Op
	DryRun invocationType = "DryRun"
)

type invokeInput struct {
	// Choose from the following options.
	//
	//    * RequestResponse (default) – Invoke the function synchronously. Keep
	//    the connection open until the function returns a response or times out.
	//    The API response includes the function response and additional data.
	//
	//    * Event – Invoke the function asynchronously. Send events that fail
	//    multiple times to the function's dead-letter queue (if one is configured).
	//    The API response only includes a status code.
	//
	//    * DryRun – Validate parameter values and verify that the user or role
	//    has permission to invoke the function.
	InvocationType   *invocationType
	FunctionName     string
	Payload          []byte
	Release          string
	Metadata         []byte
	CorrelationID    *string
	CausationID      *string
	ExpectedRevision *string
}

type InvokeOutput struct {
	ExecutedVersion *string `location:"header" locationName:"X-Amz-Executed-Version" min:"1" type:"string"`
	FunctionError   *string `location:"header" locationName:"X-Amz-Function-Error" type:"string"`
	Payload         []byte  `type:"blob" sensitive:"true"`
	StatusCode      *int64  `location:"statusCode" type:"integer"`
}
