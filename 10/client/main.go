package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	pb "grpc_workshop/10/service"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewFileServiceClient(conn)

	stream, err := client.DownloadFile(context.Background(), &pb.FileRequest{Filename: "10/server/file.mov"})
	if err != nil {
		log.Fatalf("could not download file: %v", err)
	}

	outFile, err := os.Create("10/client/file.mov")
	if err != nil {
		log.Fatalf("failed to create file: %v", err)
	}
	defer outFile.Close()

	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("failed to receive a chunk: %v", err)
		}
		if _, err := outFile.Write(chunk.Content); err != nil {
			log.Fatalf("failed to write to file: %v", err)
		}
		fmt.Printf("\rProgress: %d%%", chunk.Progress)
		os.Stdout.Sync()
	}
}
