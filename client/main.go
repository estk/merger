package client

import (
	"context"
	"fmt"
	stdlog "log"
	"os"
	"os/exec"

	pb "github.com/estk/merger/pb"
	"github.com/rs/zerolog"

	"google.golang.org/grpc"
)

func init() {
	log := zerolog.New(os.Stdout).With().
		Str("module", "merger-client").
		Logger()

	stdlog.SetFlags(0)
	stdlog.SetOutput(log)
}

type ClientConfig struct {
	ServerAddr string
}

type Client struct {
	ClientConfig
	msc  pb.MergeServiceClient
	conn *grpc.ClientConn
}

func New() *Client {
	cc := ClientConfig{ServerAddr: "127.0.0.1:3000"}
	conn, err := grpc.Dial(cc.ServerAddr, grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
	}
	msc := pb.NewMergeServiceClient(conn)
	return &Client{cc, msc, conn}
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c Client) TrackPayload(traces []*pb.Trace, data []byte) (pb.Trace, error) {
	id := c.mkID()
	trace := pb.Trace{id, traces}
	_, err := c.msc.PartialEvents(context.Background(), &pb.EventRequest{
		[]*pb.DataWrapper{
			&pb.DataWrapper{
				&trace,
				&pb.DataMeta{"myschema", "myversion"},
				data,
			},
		},
	})
	return trace, err
}

func (c Client) TrackImpression(traces []*pb.Trace, data []byte) (pb.Trace, error) {
	id := c.mkID()
	trace := pb.Trace{id, traces}
	_, err := c.msc.CompleteEvents(context.Background(), &pb.EventRequest{
		[]*pb.DataWrapper{
			&pb.DataWrapper{
				&trace,
				&pb.DataMeta{"myschema", "myversion"},
				data,
			},
		},
	})
	return trace, err
}

// TODO: Need to prove things about collisions etc
func (c Client) mkID() string {
	out, err := exec.Command("uuidgen").Output()
	if err != nil {
		fmt.Println(err)
	}
	return string(out)
}
