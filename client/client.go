package main // client.go
import (
	"fmt"
	"log"

	"github.com/hibiken/asynq"
	"notify.is-go/tasks"
)

func main() {
	r := asynq.RedisClientOpt{Addr: "localhost:6379"}
	client := asynq.NewClient(r)

	t1 := tasks.NewEmailDeliveryTask(42, "some:template:id")

	// Process the task immediately.
	res, err := client.EnqueueAt("* * * * * *", t1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("result: %+v\n", res)

}
