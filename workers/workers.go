package main // workers.go
import (
	"log"

	"github.com/hibiken/asynq"
	"notify.is-go/tasks"
)

func main() {
	r := asynq.RedisClientOpt{Addr: "localhost:6379"}
	srv := asynq.NewServer(r, asynq.Config{
		Concurrency: 10,
	})

	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.EmailDelivery, tasks.HandleEmailDeliveryTask)

	if err := srv.Run(mux); err != nil {
		log.Fatal(err)
	}
}
