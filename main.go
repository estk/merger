package main

// Todo:
// - Use Avro to pass data.

import (
	"log"
	"os/exec"
	"time"
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

type dumpable = string

type Trace struct {
	id     string
	traces []Trace
}

type PartialEvent struct {
	traces []Trace
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
	trace Trace
	datas []dumpable
}

func (s *Server) acceptPartial(payload Payload) {
	s.eventBuffer[payload.id] = partialEventEntry{
		ts:      time.Now(),
		partial: payload.partial,
	}
	log.Println("Stored Partial:", payload)
}

func (s *Server) acceptImpression(p Payload) {
	event := s.completePartials(p)
	s.completeEvents = append(s.completeEvents, event)
	log.Println("Completed Impression: ", event)
	log.Println("Datas: ", event.datas)
}

func (s *Server) completePartials(payload Payload) completeEvent {
	traces := payload.partial.traces
	trace := Trace{payload.id, traces}
	datas := s.collectDatas(traces)
	datas = append(datas, payload.partial.data)
	return completeEvent{trace, datas}
}
func (s *Server) collectDatas(traces []Trace) []dumpable {
	datas := []dumpable{}
	for _, t := range traces {
		partialEntry := s.eventBuffer[t.id]
		cDatas := s.collectDatas(partialEntry.partial.traces)
		data := s.eventBuffer[t.id].partial.data

		datas = append(datas, cDatas...)
		datas = append(datas, data)
		delete(s.eventBuffer, t.id)
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
