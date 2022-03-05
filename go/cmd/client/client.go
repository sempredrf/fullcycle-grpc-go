package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/sempredrf/fullcycle-grpc-go/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	connection, err := grpc.Dial("localhost:5051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to gRPC Server: %v", err)
	}

	defer connection.Close()

	client := pb.NewUserServiceClient(connection)

	//-- grpc
	fmt.Printf("\n\n\n#### GRPC ####\n")
	AddUser(client)

	//-- server streaming
	fmt.Printf("\n\n\n#### Server Streaming ####\n")
	AddUserVerbose(client)

	//-- client streaming
	fmt.Printf("\n\n\n#### Client Streaming ####\n")
	AddUsers(client)

	//-- server/client streaming
	fmt.Printf("\n\n\n#### Server/Client Streaming ####\n")
	AddUsersStreamingBoth(client)
}

func AddUser(client pb.UserServiceClient) {

	req := &pb.User{
		Id:    "0",
		Name:  "diego",
		Email: "diego@email.com",
	}

	res, err := client.AddUser(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make gRPC Request: %v", err)
	}

	fmt.Println(res)
}

func AddUserVerbose(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "diego",
		Email: "diego@email.com",
	}

	responseStream, err := client.AddUserVerbose(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make gRPC Request: {%v", err)
	}

	for {
		stream, err := responseStream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Could not receive the message: %v", err)
			break
		}

		fmt.Println("Status:", stream.Status)

	}
}

func AddUsers(client pb.UserServiceClient) {
	reqs := []*pb.User{
		{
			Id:    "001",
			Name:  "diego",
			Email: "diego@email.com",
		},

		{
			Id:    "002",
			Name:  "diego 2",
			Email: "diego2@email.com",
		},

		{
			Id:    "003",
			Name:  "diego 3",
			Email: "diego3@email.com",
		},

		{
			Id:    "004",
			Name:  "diego 4",
			Email: "diego4@email.com",
		},

		{
			Id:    "005",
			Name:  "diego 5",
			Email: "diego5@email.com",
		},
	}

	stream, err := client.AddUsers(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	for _, req := range reqs {
		fmt.Println("Sending: ", req.GetName())

		stream.Send(req)
		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving respose: %v", err)
	}

	fmt.Println(res)
}

func AddUsersStreamingBoth(client pb.UserServiceClient) {

	reqs := []*pb.User{
		{
			Id:    "001",
			Name:  "diego",
			Email: "diego@email.com",
		},

		{
			Id:    "002",
			Name:  "diego 2",
			Email: "diego2@email.com",
		},

		{
			Id:    "003",
			Name:  "diego 3",
			Email: "diego3@email.com",
		},

		{
			Id:    "004",
			Name:  "diego 4",
			Email: "diego4@email.com",
		},

		{
			Id:    "005",
			Name:  "diego 5",
			Email: "diego5@email.com",
		},
	}

	stream, err := client.AddUsersStreamingBoth(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	wait := make(chan int)

	go func() {

		for _, req := range reqs {
			fmt.Println("Sending: ", req.GetName())

			stream.Send(req)
			time.Sleep(time.Second * 3)
		}

		stream.CloseSend()

	}()

	go func() {

		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("Could not receive the message: %v", err)
				break
			}

			fmt.Printf("Receiving user %v with status %v\n", res.GetUser().Name, res.GetStatus())
		}

		close(wait)

	}()

	<-wait
}
