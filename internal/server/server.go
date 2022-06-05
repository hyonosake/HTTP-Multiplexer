package server

//go:generate mockery --case underscore --name Handler

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/hyonosake/HTTP-Multiplexer/internal/types"
	"github.com/hyonosake/HTTP-Multiplexer/internal/workers"
	"go.uber.org/zap"
)

type Handler interface {
	Handle(ctx context.Context)
}

type Server struct {
	Cfg    *types.Env
	ctx    context.Context
	mux    *http.ServeMux
	logger *zap.Logger
	d      workers.Distributor
}

func New(_ context.Context, cfg *types.Env, distributor workers.Distributor) (s *Server, err error) {
	lg, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	sMux := http.NewServeMux()
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return &Server{
		logger: lg,
		mux:    sMux,
		Cfg:    cfg,
		d:      distributor,
	}, nil
}

// Handle is used to define all endpoints of the server
func (s *Server) Handle(ctx context.Context) {
	go func() {
		s.mux.HandleFunc("/multiply", s.handleMultiply)
		if err := http.ListenAndServe(":"+strconv.Itoa(s.Cfg.Port), s.mux); err != nil {
			s.logger.Fatal("Server down: ", zap.Error(err))
		}
	}()
	s.logger.Info("Server started", zap.Ints("ports", []int{s.Cfg.Port}))
}

// handleMultiply handles all requests on "/multiply" endpoint
func (s *Server) handleMultiply(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("handleMultiply operating fact", zap.Any("method", r.Method))

	req := new(types.MultiplyRequest)

	ctx := r.Context()

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Info("Error unmarshalling data ", zap.String("endpoint", r.URL.Path), zap.Error(err))
		return
	}

	err = s.checkRequest(ctx, req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	respCh := make(chan *types.MultiplyResponse)
	s.d.AddTask(workers.WorkerJob{Req: req, RespChan: respCh})
	response := <-respCh
	sendResponse(w, response, http.StatusOK)
}

func (s *Server) checkRequest(_ context.Context, req *types.MultiplyRequest) error {
	if len(req.URLs) > 20 {
		return types.ErrTooManyURLs
	}
	return nil
}

func sendResponse(w http.ResponseWriter, resp *types.MultiplyResponse, status int) error {
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		return err
	}
	return err
}
