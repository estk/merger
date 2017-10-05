package main

// Todo:
// - Use Reggy to load avro schemas.
// - Use avro to pass data.

import (
	"log"
	"os/exec"
	"time"

	pb "merger/mergerpb"

	context "golang.org/x/net/context"
)

// Test main
func main() {
	s := NewServer(nil)
	c := NewClient(s)

	// Test 1
	// t1 := c.trackPayload([]Trace{}, "partial1")
	// t2 := c.trackPayload([]Trace{t1}, "partial2")
	// c.trackImpression([]Trace{t2}, "impression")
	// time.Sleep(time.Second)

	// Test 2
	t1a := c.trackPayload([]Trace{}, "partial1a")
	t1b := c.trackPayload([]Trace{}, "partial1b")
	t2 := c.trackPayload([]Trace{t1a, t1b}, "partial2")
	c.trackImpression([]Trace{t2}, "impression")
	time.Sleep(time.Second)
}

// ======
// Shared
// ======

type dumpable = []byte

type Trace struct {
	id     string
	traces []Trace
}

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
	completeEvents []completeEvent
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
		[]completeEvent{},
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

type completeEvent struct {
	trace pb.Trace
	datas []dumpable
}

func (s *Server) PartialEvents(ctx context.Context, er *pb.EventRequest) (*pb.Empty, error) {
	for _, dw := range er.Payload {
		s.acceptPartial(*dw.Trace, dw.Data)
	}
	return &pb.Empty{}, nil
}

func (s *Server) CompleteEvents(ctx context.Context, er *pb.EventRequest) (*pb.Empty, error) {
	for _, dw := range er.Payload {
		s.acceptPartial(*dw.Trace, dw.Data)
	}
	return &pb.Empty{}, nil
}

func (s *Server) acceptPartial(trace pb.Trace, data []byte) {
	s.eventBuffer[trace.Id] = partialEventEntry{
		ts:      time.Now(),
		partial: PartialEvent{trace.Traces, data},
	}
	log.Println("Stored Partial:", trace)
}

func (s *Server) acceptImpression(p Payload) {
	event := s.completePartials(p)
	s.completeEvents = append(s.completeEvents, event)
	log.Println("Completed Impression: ", event)
	log.Println("Datas: ", event.datas)
}

func (s *Server) completePartials(payload Payload) completeEvent {
	traces := payload.partial.traces
	trace := pb.Trace{payload.id, traces}
	datas := s.collectDatas(traces)
	datas = append(datas, payload.partial.data)
	return completeEvent{trace, datas}
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
}

type Client struct {
	ClientConfig
	server *Server
}

func NewClient(server *Server) *Client {
	return &Client{ClientConfig{}, server}
}

// TODO: This should send avro data over the wire
func (c Client) sendPartial(p Payload) {
	log.Println("Sending Partial:", p)
	c.server.acceptPartial(p)
}

func (c Client) sendImpression(p Payload) {
	log.Println("Sending Impression:", p)
	c.server.acceptImpression(p)
}

func (c Client) trackPayload(traces []Trace, data dumpable) Trace {
	id := mkID()
	// TODO: maybe change to queue
	go c.sendPartial(Payload{id: id, partial: PartialEvent{traces: traces, data: data}})
	return Trace{id: id, traces: traces}
}

func (c Client) trackImpression(traces []Trace, data dumpable) Trace {
	id := mkID()
	// TODO: maybe change to queue
	go c.sendImpression(Payload{id: id, partial: PartialEvent{traces: traces, data: data}})
	return Trace{id: id, traces: traces}
}

// TODO: Need to prove things about collisions etc
func mkID() string {
	out, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(out)
}
