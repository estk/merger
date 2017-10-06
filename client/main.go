package client

import (
	"context"
	"fmt"
	"log"
	"os/exec"

	pb "github.com/estk/merger/pb"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("vim-go")
}

type ClientConfig struct {
	ServerAddr string
}

type Client struct {
	ClientConfig
	msc pb.MergeServiceClient
}

func New() *Client {
	cc := ClientConfig{ServerAddr: "127.0.0.1:3000"}
	conn, err := grpc.Dial(cc.ServerAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	msc := pb.NewMergeServiceClient(conn)
	return &Client{cc, msc}
}

func (c Client) TrackPayload(traces []*pb.Trace, data []byte) pb.Trace {
	id := mkID()
	trace := pb.Trace{id, traces}
	go c.msc.PartialEvents(context.Background(), &pb.EventRequest{
		[]*pb.DataWrapper{
			&pb.DataWrapper{
				&trace,
				&pb.DataMeta{"myschema", "myversion"},
				data,
			},
		},
	})
	return trace
}

func (c Client) TrackImpression(traces []*pb.Trace, data []byte) pb.Trace {
	id := mkID()
	trace := pb.Trace{id, traces}
	go c.msc.CompleteEvents(context.Background(), &pb.EventRequest{
		[]*pb.DataWrapper{
			&pb.DataWrapper{
				&trace,
				&pb.DataMeta{"myschema", "myversion"},
				data,
			},
		},
	})
	return trace
}

// TODO: Need to prove things about collisions etc
func mkID() string {
	out, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(out)
}
