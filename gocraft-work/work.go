package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

// Make a redis pool
var redisPool = &redis.Pool{
	MaxActive: 5,
	MaxIdle:   5,
	Wait:      true,
	Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp", ":6379")
	},
}

type Context struct {
	// customerID int64
}

// Make an enqueuer with a particular namespace
var enqueuer = work.NewEnqueuer("my_app_namespace", redisPool)

func main() {

	// Enqueue a job named "send_email" with the specified parameters.
	_, err := enqueuer.Enqueue("send_email", work.Q{"address": "test@example.com", "subject": "hello world"})
	if err != nil {
		log.Fatal(err)
	}

	// _, err = enqueuer.EnqueueIn("send_welcome_email", 60, work.Q{"address": "new@example.com", "subject": "hello world"})
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Make a new pool. Arguments:
	// Context{} is a struct that will be the context for the request.
	// 10 is the max concurrency
	// "my_app_namespace" is the Redis namespace
	// redisPool is a Redis pool
	pool := work.NewWorkerPool(Context{}, 10, "my_app_namespace", redisPool)

	// Add middleware that will be executed for each job
	pool.Middleware((*Context).Log)

	// Map the name of jobs to handler functions
	pool.PeriodicallyEnqueue("30 * * * * *", "send_email")
	pool.Job("send_email", (*Context).SendEmail)

	// pool.Job("send_welcome_email", (*Context).SendEmail)

	// Customize options:
	// pool.JobWithOptions("export", work.JobOptions{Priority: 10, MaxFails: 1}, (*Context).Export)

	// Start processing jobs
	pool.Start()

	// Wait for a signal to quit:
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan

	// Stop the pool
	pool.Stop()
}

func (c *Context) Log(job *work.Job, next work.NextMiddlewareFunc) error {
	fmt.Println("Starting job: ", job.Name)
	return next()
}

func (c *Context) SendEmail(job *work.Job) error {
	// Extract arguments:
	addr := job.ArgString("address")
	subject := job.ArgString("subject")
	fmt.Println(addr)
	fmt.Println(subject)
	if err := job.ArgError(); err != nil {
		return err
	}

	fmt.Println("Pretending to send email...")
	// Go ahead and send the email...
	// sendEmailTo(addr, subject)

	return nil
}

// func (c *Context) Export(job *work.Job) error {
// 	return nil
// }
