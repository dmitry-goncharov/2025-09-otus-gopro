package internalgrpc

//go:generate protoc -I ../../../api EventService.proto --go_out=. --go-grpc_out=.

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/app"
	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	addr string
	log  app.Logger
	app  app.Application
	srv  *grpc.Server
}

func NewServer(conf *config.ServerConf, log app.Logger, app app.Application) *Server {
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(logInterceptor(log)))
	grpc.ChainStreamInterceptor()
	server := &Server{
		addr: net.JoinHostPort(conf.Host, conf.Port),
		log:  log,
		app:  app,
		srv:  grpcServer,
	}

	RegisterEventServiceServer(grpcServer, server)

	return server
}

func logInterceptor(log app.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ any, err error) {
		log.Debug(fmt.Sprintf("before method: %s, request: %+v", info.FullMethod, req))
		resp, err := handler(ctx, req)
		log.Debug(fmt.Sprintf("after method: %s, request: %+v", info.FullMethod, req))
		return resp, err
	}
}

func (s *Server) Start(_ context.Context) error {
	lsn, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("failed to listen an address: %w", err)
	}
	s.log.Info("grpc server is starting on addr: " + s.addr)
	err = s.srv.Serve(lsn)
	if err != nil {
		s.log.Error("failed to start grpc server: " + err.Error())
	}
	return nil
}

func (s *Server) Stop(_ context.Context) {
	s.log.Info("grpc server is stopping")
	s.srv.GracefulStop()
}

func (s *Server) AddEvent(ctx context.Context, req *CreateEventReq) (*CreateEventResp, error) {
	ctx, cancel := context.WithTimeout(context.WithValue(ctx, app.UserIDKey, req.UserID), time.Second*5)
	defer cancel()

	appEvent, err := s.app.CreateEvent(ctx, req.Title)
	if err != nil {
		return nil, fmt.Errorf("error adding event: %w", err)
	}

	event := &Event{
		ID:     appEvent.ID,
		Title:  appEvent.Title,
		Date:   timestamppb.New(appEvent.Date),
		UserID: appEvent.UserID,
	}

	resp := &CreateEventResp{
		Event: event,
	}

	return resp, nil
}

func (s *Server) UpdateEvent(ctx context.Context, req *UpdateEventReq) (*UpdateEventResp, error) {
	appEvent := mapGrpcEventToAppEvent(req.Event)

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	err := s.app.UpdateEvent(ctx, appEvent)
	if err != nil {
		return nil, fmt.Errorf("error updating event: %w", err)
	}

	event := &Event{
		ID:     appEvent.ID,
		Title:  appEvent.Title,
		Date:   timestamppb.New(appEvent.Date),
		UserID: appEvent.UserID,
	}

	resp := &UpdateEventResp{
		Event: event,
	}

	return resp, nil
}

func (s *Server) DeleteEvent(ctx context.Context, req *DeleteEventReq) (*DeleteEventResp, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	err := s.app.DeleteEvent(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("error deleting event: %w", err)
	}

	resp := &DeleteEventResp{}

	return resp, nil
}

func (s *Server) GetEvents(ctx context.Context, req *GetEventsReq) (*GetEventsResp, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	date := time.Date(int(req.Year), time.Month(req.Month), int(req.Day), 0, 0, 0, 0, time.Local)

	var events []app.Event
	var err error
	switch req.Period {
	case GetEventsReq_PERIOD_DAY:
		events, err = s.app.GetDayEvents(ctx, date)
		if err != nil {
			return nil, fmt.Errorf("error getting day events: %w", err)
		}
	case GetEventsReq_PERIOD_WEEK:
		events, err = s.app.GetWeekEvents(ctx, date)
		if err != nil {
			return nil, fmt.Errorf("error getting week events: %w", err)
		}
	case GetEventsReq_PERIOD_MONTH:
		events, err = s.app.GetMonthEvents(ctx, date)
		if err != nil {
			return nil, fmt.Errorf("error getting month events: %w", err)
		}
	case GetEventsReq_PERIOD_UNSPECIFIED:
		return nil, fmt.Errorf("bad request: period unspecified")
	}

	resp := &GetEventsResp{
		Events: mapAppEventsToGrpsEvents(events),
	}

	return resp, nil
}

func (s *Server) mustEmbedUnimplementedEventServiceServer() {}

func mapGrpcEventToAppEvent(grpcEvent *Event) *app.Event {
	return &app.Event{
		ID:     grpcEvent.ID,
		Title:  grpcEvent.Title,
		Date:   grpcEvent.Date.AsTime(),
		UserID: grpcEvent.UserID,
	}
}

func mapAppEventToGrpsEvent(appEvent *app.Event) *Event {
	return &Event{
		ID:     appEvent.ID,
		Title:  appEvent.Title,
		Date:   timestamppb.New(appEvent.Date),
		UserID: appEvent.UserID,
	}
}

func mapAppEventsToGrpsEvents(appEvents []app.Event) []*Event {
	grpcEvents := make([]*Event, len(appEvents))
	for i, appEvent := range appEvents {
		grpcEvents[i] = mapAppEventToGrpsEvent(&appEvent)
	}
	return grpcEvents
}
