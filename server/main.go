package server

import (
	"fmt"
	stdlog "log"
	"net"
	"os"
	"time"

	"golang.org/x/net/context"

	pb "github.com/estk/merger/pb"
	"github.com/rs/zerolog"
	log "github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type ServerConfig struct {
	ServerAddr string
	TTL        time.Duration
}

type Server struct {
	ServerConfig
	eventBuffer     partialMap
	CompletedEvents []CompleteEvent
	grpcServer      *grpc.Server
}

type partialMap = map[string]partialEventEntry

type partialEventEntry struct {
	ts     time.Time
	traces []*pb.Trace
	data   dumpable
}

type dumpable = []byte

func init() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log := zerolog.New(os.Stdout).With().
		Str("module", "merger-server").
		Logger()

	stdlog.SetFlags(0)
	stdlog.SetOutput(log)
}

func New(sc *ServerConfig) *Server {
	if sc == nil {
		sc = &ServerConfig{
			ServerAddr: ":3000",
			TTL:        30 * time.Second,
		}
	}

	lis, err := net.Listen("tcp", sc.ServerAddr)
	if err != nil {
		fmt.Errorf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	s := &Server{
		*sc,
		make(partialMap),
		[]CompleteEvent{},
		grpcServer,
	}
	pb.RegisterMergeServiceServer(grpcServer, s)

	go grpcServer.Serve(lis)

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
	Datas []dumpable
}

func (s *Server) Stop() {
	s.grpcServer.Stop()
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
	log.Debug().Msgf("Stored Partial: %s", trace)
}

func (s *Server) processComplete(trace pb.Trace, data []byte) {
	event := s.completePartials(trace, data)
	s.CompletedEvents = append(s.CompletedEvents, event)
	log.Debug().Msgf("Completed Impression: %s", event)
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
