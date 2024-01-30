package golang

import (
	"fmt"
	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	services "github.com/nats-io/nats.go/micro"
	"os"
	"strings"
)

func MustNewEventStore() *esdb.Client {
	settings, err := esdb.ParseConnectionString("esdb://127.0.0.1:2113?tls=false&keepAliveTimeout=10000&keepAliveInterval=10000")
	Must(err)

	client, err := esdb.NewClient(settings)
	Must(err)

	return client
}

func NewNats() (*nats.Conn, jetstream.JetStream) {
	natsUrl := os.Getenv("NATS_URL")
	if len(strings.TrimSpace(natsUrl)) == 0 {
		natsUrl = nats.DefaultURL
	}
	nc, err := nats.Connect(natsUrl)
	Must(err)

	js, err := jetstream.New(nc)
	Must(err)

	return nc, js
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

type CommandHandlerResponse struct {
	NextExpectedVersion uint64 `json:"nextExpectedVersion"`
}
type ServiceCommandHandler[Command any] func(command Command) (
	*CommandHandlerResponse,
	error,
)

type UnmarshalCommand[Command any] func(data []byte) (Command, error)

func NewService[Command any](
	nc *nats.Conn,
	serviceName string,
	subjectName string,
	unmarshalCommand UnmarshalCommand[Command],
	appHandler ServiceCommandHandler[Command],
) (services.Service, error) {
	fullSubjectName := fmt.Sprintf("svc.onepiece.%s", subjectName)
	fullServiceName := fmt.Sprintf("onepiece-%s", serviceName)

	fmt.Printf("service name: %s\n", fullServiceName)
	fmt.Printf("service subject: %s\n", fullSubjectName)
	return services.AddService(nc, services.Config{
		Name:    fullServiceName,
		Version: "1.0.0",
		Endpoint: &services.EndpointConfig{
			Subject: fullSubjectName,
			Handler: services.HandlerFunc(func(req services.Request) {
				command, err := unmarshalCommand(req.Data())
				if err != nil {
					req.Error("error", err.Error(), nil)
					return
				}

				resp, err := appHandler(command)
				if err != nil {
					req.Error("error", err.Error(), nil)
					return
				}

				req.RespondJSON(resp)
			}),
		},
	})
}
