package main

import (
	"log"
	"testing"

	client "github.com/estk/merger/client"
	pb "github.com/estk/merger/pb"
	server "github.com/estk/merger/server"
)

func TestSum(t *testing.T) {
	t.Errorf("Sum was incorrect, got: %d, want: %d.", 9, 10)
}

func TestToComplete(t *testing.T) {
	t.Errorf("run")
	server.New(nil)
	log.Println("server")
	client := client.New()
	log.Println("client")

	trace1 := client.TrackPayload([]*pb.Trace{}, []byte("hello!"))
	log.Println("1")
	trace2 := client.TrackPayload([]*pb.Trace{&trace1}, []byte("bon jour"))
	log.Println("2")
	client.TrackImpression([]*pb.Trace{&trace2}, []byte("ni hau"))
}
