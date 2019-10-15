package main

import (
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"

	"github.com/nats-io/nats.go"
)

func main() {
	server := flag.String("c", nats.DefaultURL, "Specify a NATS connection string")

	file := flag.String("f", "", "Specify a JSON file to publish to a subject")
	subject := flag.String("s", "", "Specify a subject to publish to")
	flag.Parse()

	conn, _ := nats.Connect(*server)

	if len(*file) > 0 && len(*subject) > 0 {
		data, err := ioutil.ReadFile(*file)
		if err != nil {
			panic(err)
		}
		if !json.Valid(data) {
			panic(errors.New("Invalid JSON supplied"))
		}
		err = conn.Publish(*subject, data)
		if err != nil {
			panic(err)
		}
	}

	conn.Drain()
	conn.Close()
}
