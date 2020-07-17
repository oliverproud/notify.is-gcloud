package main // client.go
import (
	"fmt"
	"os"
	"os/signal"

	"github.com/hibiken/asynq"
	"notify.is-go/tasks"
)

func main() {
	r := asynq.RedisClientOpt{Addr: "localhost:6379"}
	client := asynq.NewClient(r)

	t := tasks.NewEmailDeliveryTask(42, "some:template:id")

	// Process the task immediately.
	res, err := client.Enqueue(t)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("result: %+v\n", res)

	// Wait for a signal to quit:
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan

}
