package server

//go:generate mockery --case underscore --name Handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	b64 "encoding/base64"

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
	m      workers.Worker
}

func New(_ context.Context) (s *Server, err error) {
	lg, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	sMux := http.NewServeMux()
	cfg, err := s.parseConfig()
	if err != nil {
		return nil, err
	}
	w, err := workers.New(cfg)
	if err != nil {
		return nil, err
	}
	return &Server{
		logger: lg,
		mux:    sMux,
		Cfg:    cfg,
		m:      w,
	}, nil
}

func (s *Server) sendResponse(w http.ResponseWriter, resp *types.MultiplyResponse, status int) {
	w.WriteHeader(status)
	bytez, err := json.Marshal(resp)
	if err != nil {
		s.logger.Error("Unable to marshal response", zap.Any("response", resp), zap.Error(err))
	}
	w.Write(bytez)
}

func (s *Server) Handle(ctx context.Context) {
	go func() {
		s.mux.HandleFunc("/multiply", s.handleMultiply)
		if err := http.ListenAndServe(":"+strconv.Itoa(s.Cfg.Port), s.mux); err != nil {
			s.logger.Fatal("Server down: ", zap.Error(err))
		}
	}()
	s.logger.Info("Server started", zap.Ints("ports", []int{s.Cfg.Port}))
}

func (s *Server) parseConfig() (*types.Env, error) {
	return &types.Env{
		PoolWorkersSize: 100,
		URLAmountLimit:  4,
		URLQueryTimeout: time.Second,
		ShutdownTimeout: time.Second * 5,
		Port:            1234,
	}, nil
}

// HandleFact handles GET and POST for URL/fact
func (s *Server) handleMultiply(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("handleMultiply operating fact", zap.Any("method", r.Method))

	var (
		req  = types.MultiplyRequest{}
		resp *types.MultiplyResponse
	)

	ctx := r.Context()
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.logger.Info("Unable to parse data body", zap.String("endpoint", r.URL.Path), zap.Error(err))
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &req)
	if err != nil {
		s.logger.Info("Error unmarshalling data ", zap.String("endpoint", r.URL.Path), zap.Error(err))
		return
	}

	err = s.checkRequest(ctx, &req)
	if err != nil {
		resp.Error = err.Error()
		s.sendResponse(w, resp, http.StatusBadRequest)
		return
	}
	resp = s.multiply(ctx, &req)
	s.sendResponse(w, resp, http.StatusOK)
	//resp, err := s.postNewFacts(r)
}

// prob worker's job
func (s *Server) multiply(ctx context.Context, req *types.MultiplyRequest) *types.MultiplyResponse {
	//eg, egCtx := errgroup.WithContext(ctx)
	var datas []types.JsonData
	for _, url := range req.URLs {
		d := request(ctx, url)
		datas = append(datas, d)
	}
	return &types.MultiplyResponse{Data: datas}
}

func request(ctx context.Context, url string) types.JsonData {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("ERR 1: ", err.Error())
		return types.JsonData{}
		// log me
	}
	req = req.WithContext(ctx)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("ERR 2: ", err.Error())
		return types.JsonData{}
		// log me
	}
	defer resp.Body.Close()
	var data []byte
	data, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ERR 3")
		return types.JsonData{}
		// log
	}
	strData := b64.StdEncoding.EncodeToString(data)
	return types.JsonData{URL: url, Data: strData}
}

// prob worker's job
func (s *Server) checkRequest(_ context.Context, req *types.MultiplyRequest) error {
	if len(req.URLs) > 20 {
		return types.ErrTooManyURLs
	}
	return nil
}
