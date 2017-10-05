package main

import (
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

func NewClient() *Client {
	cc := ClientConfig{ServerAddr: "localhost:3000"}
	conn, err := grpc.Dial(cc.ServerAddr)
	if err != nil {
		log.Fatal(err)
	}
	msc := pb.NewMergeServiceClient(conn)
	return &Client{cc, msc}
}

func (c Client) TrackPayload(traces []*pb.Trace, data []byte) pb.Trace {
	id := mkID()
	trace := pb.Trace{id, traces}
	go c.msc.PartialEvent(nil, &pb.EventRequest{
		[]*pb.DataWrapper{
			&pb.DataWrapper{
				&trace,
				&pb.DataMeta{"myschema", "myversion"},
				[]byte("mydata"),
			},
		},
	})
	return trace
}

func (c Client) TrackImpression(traces []*pb.Trace, data []byte) pb.Trace {
	id := mkID()
	trace := pb.Trace{id, traces}
	go c.msc.CompleteEvent(nil, &pb.EventRequest{
		[]*pb.DataWrapper{
			&pb.DataWrapper{
				&trace,
				&pb.DataMeta{"myschema", "myversion"},
				[]byte("mydata"),
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
