package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	pb "github.com/KathurimaKimathi/shippy/shippy-service-consignment/proto/consignment"
	micro "github.com/micro/go-micro/v2"
)

const (
	address         = "localhost:50051"
	defaultFilename = "consignment.json"
)

func parseFile(file string) (*pb.Consignment, error) {
	var consignment *pb.Consignment
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(data, &consignment)
	return consignment, err
}

func main() {
	service := micro.NewService(micro.Name("shippy.consignment.cli"))
	service.Init()

	client := pb.NewShippingService("shippy.service.consignment", service.Client())

	//Contact the server and print out its response
	file := defaultFilename
	if len(os.Args) > 1 {
		file = os.Args[1]
	}

	consignment, err := parseFile(file)
	if err != nil {
		log.Fatalf("Could not parse file %s: %v", file, err)
	}

	r, err := client.CreateConsignment(context.Background(), consignment)
	if err != nil {
		log.Fatalf("Could not greet: %v", err)
	}
	log.Printf("Created: %t", r.Created)

	getAll, err := client.GetConsignment(context.Background(), &pb.GetRequest{})
	if err != nil {
		log.Fatalf("Could not list consignment: %v", err)
	}
	for _, v := range getAll.Consignments {
		log.Println(v)
	}
}
