package main

import (
	"testing"
	"time"

	client "github.com/estk/merger/client"
	pb "github.com/estk/merger/pb"
	server "github.com/estk/merger/server"
)

func TestToComplete(t *testing.T) {
	s := server.New(nil)

	time.Sleep(2 * time.Second)
	client := client.New()

	trace1 := client.TrackPayload([]*pb.Trace{}, []byte("data1"))
	trace2 := client.TrackPayload([]*pb.Trace{&trace1}, []byte("data2"))
	client.TrackImpression([]*pb.Trace{&trace2}, []byte("data3"))
	time.Sleep(1 * time.Second)

	if len(s.CompletedEvents) == 0 {
		t.Error("No completed events")
	}

	for _, e := range s.CompletedEvents {
		for _, d := range e.Datas {
			s := string(d)
			if !(s != "data1" || s != "data2" || s != "data3") {
				t.Errorf("Unexpected data %s", s)
			}
		}
	}
}
