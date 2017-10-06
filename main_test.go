package main

// TODO:
// - Back with datastore
// - Deserialize and repackage avro.

import (
	"log"
	"testing"
	"time"

	client "github.com/estk/merger/client"
	pb "github.com/estk/merger/pb"
	server "github.com/estk/merger/server"
)

func BenchmarkComplete(b *testing.B) {
	s := server.New(nil)
	client := client.New()
	defer client.Close()
	defer s.Stop()

	for i := 1; i <= b.N; i++ {
		trace1, _ := client.TrackPayload(mkTraces(), []byte("data1"))
		trace2, _ := client.TrackPayload(mkTraces(&trace1), []byte("data2"))
		client.TrackImpression(mkTraces(&trace2), []byte("data3"))
	}
	time.Sleep(1 * time.Second)
	if len(s.CompletedEvents) != b.N {
		log.Fatal("Not all events completed")
	}

}
func mkTraces(ts ...*pb.Trace) []*pb.Trace {
	return ts
}

func TestToComplete(t *testing.T) {
	s := server.New(nil)
	client := client.New()
	defer client.Close()
	defer s.Stop()

	trace1, _ := client.TrackPayload(mkTraces(), []byte("data1"))
	trace2, _ := client.TrackPayload(mkTraces(&trace1), []byte("data2"))
	client.TrackImpression(mkTraces(&trace2), []byte("data3"))

	if len(s.CompletedEvents) != 1 {
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
	time.Sleep(1 * time.Second)
}
