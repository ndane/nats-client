package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	server := flag.String("c", nats.DefaultURL, "Specify a NATS connection string")

	file := flag.String("f", "", "Specify a JSON file to publish to a subject")
	subject := flag.String("s", "", "Specify a subject to publish to")
	request := flag.Bool("request", false, "Specify that this is a request and wait for a reply")
	timeout := flag.Int("t", int(10*time.Second), "Specify a request timeout")
	flag.Parse()

	conn, _ := nats.Connect(*server)

	if len(*subject) > 0 {
		var data []byte
		var err error

		if len(*file) > 0 {
			fmt.Println("Reading ", *file)
			data, err = ioutil.ReadFile(*file)
			if err != nil {
				panic(err)
			}
			if !json.Valid(data) {
				panic(errors.New("Invalid JSON supplied"))
			}
		} else {
			data = *new([]byte)
		}

		fmt.Println("Sending ", string(data))

		if *request == true {
			if response, err := conn.Request(*subject, data, time.Duration(*timeout)); err != nil {
				panic(err)
			} else {
				fmt.Println("Reply: ", string(response.Data))
			}
		} else {
			if err := conn.Publish(*subject, data); err != nil {
				panic(err)
			}
		}
	}

	conn.Drain()
	conn.Close()
}