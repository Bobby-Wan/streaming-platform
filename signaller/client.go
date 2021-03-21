package signaller

import(
	"google.golang.org/grpc"
	"flag"
	"log"
)

func SignalChange(signal StreamStateChangeSignal) *grpc.SignalServiceClient{
	serverAddr := flag.String("server_addr", "localhost:8080", "The server address in the format of host:port")

	conn, err:= grpc.Dial(*serverAddr, grpc.WithInsecure())
	if err!=nil{
		log.Println("could not connect to %s: %v",*serverAddr, err)
		return nil
	}
	defer conn.Close()

	client := NewSignalServiceClient(conn)
	client.SignalStreamStateChange(context.Background(), signal, )
	
	return client
}