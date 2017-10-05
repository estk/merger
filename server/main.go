package main

import (
	"context"
	"fmt"
	"log"
	pb "merger/pb"
	"time"
)

func main() {
	fmt.Println("vim-go")
}

type ServerConfig struct {
	TTL time.Duration
}

type Server struct {
	ServerConfig
	eventBuffer    partialMap
	completeEvents []CompleteEvent
}

type partialMap = map[string]partialEventEntry

type partialEventEntry struct {
	ts     time.Time
	traces []*pb.Trace
	data   dumpable
}

type dumpable = []byte

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
		time.Now(),
		trace.Traces,
		data,
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
		cDatas := s.collectDatas(partialEntry.traces)
		data := s.eventBuffer[t.Id].data

		datas = append(datas, cDatas...)
		datas = append(datas, data)
		delete(s.eventBuffer, t.Id)
	}
	return datas
}
