package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/plugin"
)

func main() {
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	request := &plugin_go.CodeGeneratorRequest{}
	if err := proto.Unmarshal(data, request); err != nil {
		log.Fatal(err)
	}

	response := &plugin_go.CodeGeneratorResponse{}

	generator := New(request, response)
	if err := generator.Generate(); err != nil {
		log.Fatal(err)
	}

	out, err := proto.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stdout.Write(out); err != nil {
		log.Fatal(err)
	}
}
