package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/tlandenberger/grpc-todo/pb"
)

func main() {
	serverAddr := flag.String(
		"server", "localhost:8080",
		"The server address in the format of host:port",
	)
	flag.Parse()
  
  var opts []grpc.DialOption

  if *serverAddr == "localhost:8080" {
    opts = append(opts, grpc.WithInsecure())
  } else {
    creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: false})
    opts = append(opts, grpc.WithTransportCredentials(creds))
  }
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, *serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial:", err)
	}
	defer conn.Close()

	client := pb.NewToDoServiceClient(conn)

  todo := &pb.ToDo{
    Title:       "Learn gRPC",
    Description: "Experiment with gRPC and Go",
  }
  res, err := client.CreateToDo(ctx, &pb.CreateToDoRequest{Todo: todo})
	if err != nil {
		log.Fatalf("error sending request:", err)
	}

	fmt.Println("result:", res);
}
