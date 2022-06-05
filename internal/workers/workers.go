package workers

//go:generate mockery --case underscore --name WorkerPool
//go:generate mockery --case underscore --name Distributor

import (
	"context"
	b64 "encoding/base64"
	"io"
	"net/http"
	"time"

	"github.com/hyonosake/HTTP-Multiplexer/internal/types"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type WorkerPool interface {
	Distributor
}

type Distributor interface {
	RunWorkers(ctx context.Context)
	Close() error
	AddTask(job WorkerJob)
}

type WorkerJob struct {
	Req      *types.MultiplyRequest
	RespChan chan *types.MultiplyResponse
}

type WorkerResult struct {
	Err  error
	Resp *types.MultiplyResponse
	out  chan *types.MultiplyResponse
}

type RequestHandler struct {
	logger *zap.Logger
	cfg    *types.Env
	jobs   chan WorkerJob
	client types.Client
}

func NewRequestHandler(cfg *types.Env, lg *zap.Logger) (*RequestHandler, error) {
	return &RequestHandler{
		logger: lg,
		cfg:    cfg,
		jobs:   make(chan WorkerJob, cfg.PoolWorkersSize),
		client: http.DefaultClient,
	}, nil
}

// Close used to close job queue. After call no more jobs can be added with AddTask()
func (r *RequestHandler) Close() error {
	r.logger.Info("Stopped accepting new connections...")
	close(r.jobs)
	return nil
}

// AddTask adds a Job task to job queue
func (r *RequestHandler) AddTask(job WorkerJob) {
	r.jobs <- job
}

// RunWorkers creates a limited pool of workers that are handling provided jobs
func (r *RequestHandler) RunWorkers(ctx context.Context) {
	streamRestraint := make(chan struct{}, r.cfg.PoolWorkersSize) // limit concurrent goroutines
	for job := range r.jobs {
		streamRestraint <- struct{}{} // this blocks if {r.cfg.PoolWorkersSize} are running
		currJob := job
		go func() {
			r.runJob(ctx, currJob)
			<-streamRestraint
		}()
	}
}

// runJob runs provided job and sends resulting response back to handler
func (r *RequestHandler) runJob(ctx context.Context, job WorkerJob) {
	aggregatedResp := new(types.MultiplyResponse)
	jobResp, err := r.run(ctx, job)
	aggregatedResp.Data = jobResp
	if err != nil {
		aggregatedResp.Error = err.Error()
	}
	job.RespChan <- aggregatedResp
}

// run is actually the main one in charge. It creates limted amount of goroutines that are synced with errgroup, so that
// if one goroutine returns err, all of the work will be stopped.
func (r *RequestHandler) run(ctx context.Context, job WorkerJob) ([]types.JsonData, error) {
	var (
		respCh          = make(chan types.JsonData, len(job.Req.URLs)) // collect all data from sent requests here
		respData        = make([]types.JsonData, 0, len(job.Req.URLs))
		streamRestraint = make(chan struct{}, r.cfg.MaxParallelQueries) // use it to handle amount of concurrently running goroutines
	)
	defer close(streamRestraint)
	eg, egCtx := errgroup.WithContext(ctx) // stops spawning goroutines if err occurs
	for _, url := range job.Req.URLs {
		streamRestraint <- struct{}{}
		innerUrl := url
		eg.Go(func() error {
			resp, err := r.makeSingleRequest(egCtx, innerUrl)
			<-streamRestraint
			if err != nil {
				return err
			}
			respCh <- types.JsonData{Data: resp, URL: innerUrl}
			return nil
		},
		)
	}
	err := eg.Wait()
	if err != nil {
		return nil, err
	}
	close(respCh)
	for data := range respCh {
		respData = append(respData, data)
	}
	return respData, nil
}

// makeSingleRequest makes a single request to provided URL. Returns ErrTimeoutRequest if request takes too long
func (r *RequestHandler) makeSingleRequest(ctx context.Context, url string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		r.logger.Error("Unable to create request", zap.Error(err))
		return "", err
	}
	req = req.WithContext(ctx)
	resp, err := r.client.Do(req)
	select {
	case <-ctx.Done():
		return "", errors.Wrap(types.ErrTimeoutRequest, "url "+url)
	default:
	}
	if err != nil {
		r.logger.Error("Unable to make request", zap.Error(err))
		return "", err
	}
	defer resp.Body.Close()
	var data []byte
	data, err = io.ReadAll(resp.Body)
	if err != nil {
		r.logger.Error("Unable to make request", zap.Error(err))
		return "", err
	}
	strData := b64.StdEncoding.EncodeToString(data)
	return strData, nil
}
