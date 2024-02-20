package main

import (
	"fmt"
	"github.com/gofrs/uuid"
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

	rep, er := nc.Request("srv.command.monitoring.create-monitoring", []byte(payload), time.Second*3)
	if er != nil {
		fmt.Println(er)
		return
	}
	fmt.Println(string(rep.Data))
}
