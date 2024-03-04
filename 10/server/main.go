package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	pb "grpc_workshop/10/service"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedFileServiceServer
}

func (s *server) DownloadFile(req *pb.FileRequest, stream pb.FileService_DownloadFileServer) error {
	file, err := os.Open(req.Filename)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)

		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	totalSize := fileInfo.Size()
	fmt.Println(totalSize, "totalsize")
	var sentSize int64 = 0
	buf := make([]byte, 65536)
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		sentSize += int64(n)
		progress := int32((float64(sentSize) / float64(totalSize)) * 100)
		if err := stream.Send(&pb.FileChunk{Content: buf[:n], Progress: progress}); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	//file, err := os.Open("10/server/file.txt")
	//if err != nil {
	//	log.Fatalf("Failed to open file: %v", err)
	//}
	//defer file.Close() // It's good practice to defer Close() immediately after checking for an error
	//
	//buffer := make([]byte, 1024)
	//
	//for {
	//	n, err := file.Read(buffer)
	//	if err != nil && err != io.EOF {
	//		log.Fatal(err)
	//	}
	//
	//	if err == io.EOF {
	//		break
	//	}
	//
	//	fmt.Println(string(buffer[:n]))
	//}

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterFileServiceServer(s, &server{})
	log.Printf("server listening at %v", listener.Addr())
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
