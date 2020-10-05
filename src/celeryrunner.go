package main

import (
	"fmt"
	. "github.com/gocelery/gocelery"
	"github.com/gomodule/redigo/redis"
	"os"
	"time"
)

func getREDISPool(redis_url string) *redis.Pool {

	// create redis connection pool
	return &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(redis_url)
			if err != nil {
				return nil, err
			}
			return c, err
		},
	}

}

func createWorker(workers int) {

	redisPool := getREDISPool(os.Getenv("REDIS_URL"))

	// initialize celery client
	cli, _ := NewCeleryClient(
		NewRedisBroker(redisPool),
		&RedisCeleryBackend{Pool: redisPool},
		workers, // number of workers
	)

	// register tasks
	doreport_a := func(a, b string) string {
		return sqltasker_a(a, b)
	}
	cli.Register("doreport_a", doreport_a)

	doreport_b := func(a, b string) string {
		return sqltasker_b(a, b)
	}
	cli.Register("doreport_b", doreport_b)

	doreport_c := func(a, b string) string {
		return sqltasker_c(a, b)
	}
	cli.Register("doreport_c", doreport_c)

	// start workers (non-blocking call) - they`ll stop at program exit anyway
	cli.StartWorker()

}

func celeryClient(taskName, a, b string) {

	redisPool := getREDISPool(os.Getenv("REDIS_URL"))

	// initialize celery client
	cli, _ := NewCeleryClient(
		NewRedisBroker(redisPool),
		&RedisCeleryBackend{Pool: redisPool},
		1,
	)

	// run task
	asyncResult, err := cli.Delay(taskName, a, b)

	checkErr(err)

	// get results from backend with timeout
	res, err := asyncResult.Get(10 * time.Second)
	checkErr(err)

	fmt.Printf("task finished, result: %+v \r\n", res)

}
