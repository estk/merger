package main

// Todo:
// - Refactor into separate packages
// - Use Reggy to load avro schemas.

import (
	"log"
	"os/exec"
	"time"

	pb "merger/mergerpb"

	context "golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Test main
func main() {
	// s := NewServer(nil)
	// c := NewClient(s)

	// Test 1
	// t1 := c.trackPayload([]Trace{}, "partial1")
	// t2 := c.trackPayload([]Trace{t1}, "partial2")
	// c.trackImpression([]Trace{t2}, "impression")
	// time.Sleep(time.Second)

	// Test 2
	// t1a := c.trackPayload([]Trace{}, "partial1a")
	// t1b := c.trackPayload([]Trace{}, "partial1b")
	// t2 := c.trackPayload([]Trace{t1a, t1b}, "partial2")
	// c.trackImpression([]Trace{t2}, "impression")
	// time.Sleep(time.Second)
}

// ======
// Shared
// ======

type dumpable = []byte

type PartialEvent struct {
	traces []*pb.Trace
	data   dumpable
}

type Payload struct {
	id      string
	partial PartialEvent
}

// ======
// SERVER
// ======

type ServerConfig struct {
	TTL time.Duration
}

type Server struct {
	ServerConfig
	eventBuffer    partialMap
	completeEvents []CompleteEvent
}

type partialEventEntry struct {
	ts      time.Time
	partial PartialEvent
}

type partialMap = map[string]partialEventEntry

func NewServer(sc *ServerConfig) *Server {
	if sc == nil {
		sc = &ServerConfig{
			TTL: 30 * time.Second,
		}
	}
	s := &Server{
		*sc,
		make(partialMap),
		[]CompleteEvent{},
	}

	go func() {
		for _ = range time.Tick(s.TTL) {
			cleanBuffer(s.eventBuffer)
		}
	}()
	return s
}

func cleanBuffer(b partialMap) {
	tcutoff := time.Now().Add(time.Duration(60) * time.Second)
	for k, v := range b {
		if v.ts.Before(tcutoff) {
			delete(b, k)
		}
	}
}

type CompleteEvent struct {
	trace pb.Trace
	datas []dumpable
}

func (s *Server) PartialEvents(ctx context.Context, er *pb.EventRequest) (*pb.Empty, error) {
	for _, dw := range er.Payload {
		s.processPartial(*dw.Trace, dw.Data)
	}
	return &pb.Empty{}, nil
}

func (s *Server) CompleteEvents(ctx context.Context, er *pb.EventRequest) (*pb.Empty, error) {
	for _, dw := range er.Payload {
		s.processComplete(*dw.Trace, dw.Data)
	}
	return &pb.Empty{}, nil
}

func (s *Server) processPartial(trace pb.Trace, data []byte) {
	s.eventBuffer[trace.Id] = partialEventEntry{
		ts:      time.Now(),
		partial: PartialEvent{trace.Traces, data},
	}
	log.Println("Stored Partial:", trace)
}

func (s *Server) processComplete(trace pb.Trace, data []byte) {
	event := s.completePartials(trace, data)
	s.completeEvents = append(s.completeEvents, event)
	log.Println("Completed Impression: ", trace)
}

func (s *Server) completePartials(trace pb.Trace, data []byte) CompleteEvent {
	datas := s.collectDatas(trace.Traces)
	datas = append(datas, data)
	return CompleteEvent{trace, datas}
}

func (s *Server) collectDatas(traces []*pb.Trace) []dumpable {
	datas := []dumpable{}
	for _, t := range traces {
		partialEntry := s.eventBuffer[t.Id]
		cDatas := s.collectDatas(partialEntry.partial.traces)
		data := s.eventBuffer[t.Id].partial.data

		datas = append(datas, cDatas...)
		datas = append(datas, data)
		delete(s.eventBuffer, t.Id)
	}
	return datas
}

// ======
// CLIENT
// ======

type ClientConfig struct {
	ServerAddr string
}

type Client struct {
	ClientConfig
	msc pb.MergeServiceClient
}

func NewClient(server *Server) *Client {
	cc := ClientConfig{ServerAddr: "localhost:3000"}
	conn, err := grpc.Dial(cc.ServerAddr)
	if err != nil {
		log.Fatal(err)
	}
	msc := pb.NewMergeServiceClient(conn)
	return &Client{cc, msc}
}

func (c Client) TrackPayload(traces []*pb.Trace, data dumpable) pb.Trace {
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

func (c Client) TrackImpression(traces []*pb.Trace, data dumpable) pb.Trace {
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
