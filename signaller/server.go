package signaller

import(
	"net"
	"log"
	grpc "google.golang.org/grpc"
	"fmt"
	"context"
	"github.com/bobby-wan/streaming-platform/webserver/db"
)

type StreamSignalServer struct{
	UnimplementedSignalServiceServer
	streamController *StreamControllerInterface
}

func Serve(port int) error{
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d",port))
	if err!=nil{
		log.Fatalf("could not listen to port %d: %v", port, err)
	}

	s := grpc.NewServer()
	RegisterSignalServiceServer(s, StreamSignalServer{})

	ptrDb, err:= db.Initialize()
	if err!=nil{
		log.Fatalf("could not connect to db: %v", err)
	}

	streamController, err := db.NewStreamControllerGORM(ptrDb)
	if err!=nil{
		log.Fatal("failed to create stream controller: %s", err)
	}

	err = s.Serve(lis)
	if err!=nil{
		log.Fatalf("could not serve: %v", err)
	}
}

func (StreamSignalServer) SignalStreamStateChange(ctx context.Context, signal *StreamStateChangeSignal) (*SignalResponse, error){
	if signal==nil{
		log.Println("invalid data at grpc signal")
		return nil, fmt.Errorf("invalid signal")
	}

	streamController.Update()
	
	
	resp := SignalResponse{Ok:true}
	return &resp, nil
}