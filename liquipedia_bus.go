package liquiwrap

import (
	"context"
	"github.com/mnlght/liquiwrap/internal"
	"sync"
	"time"
)

const LiquipediaBusTimeout = time.Second * 30
const LiquipediaBusQueueDelay = time.Millisecond * 750

type LiquipediaResponse struct {
	Body  []byte
	Error error
}

type LiquipediaRequest struct {
	Game       string
	Params     map[string]string
	ResponseCh chan LiquipediaResponse
}

type LiquipediaBus struct {
	requests []*LiquipediaRequest
	mutex    sync.Mutex
	ctx      context.Context
	stopCh   chan struct{}
	Stop     bool
}

func NewLiquipediaBus(ctx context.Context) *LiquipediaBus {
	return &LiquipediaBus{
		mutex:  sync.Mutex{},
		ctx:    ctx,
		stopCh: make(chan struct{}),
	}
}

func (lb *LiquipediaBus) MustRun() {
	seekStep := time.NewTicker(LiquipediaBusQueueDelay)
	go func() {
		for {
			select {
			case <-lb.stopCh:
				lb.Stop = true
				time.Sleep(LiquipediaBusTimeout)
				lb.Stop = false
			case <-lb.ctx.Done():
				return
			default:
				continue
			}
		}
	}()

	for {
		select {
		case <-lb.ctx.Done():
			return
		case <-seekStep.C:
			if lb.Stop == true {
				continue
			}
			if len(lb.requests) > 0 {
				if r := lb.extract(); r != nil {
					go func() {
						lb.stopCh <- struct{}{}
					}()
					client := internal.NewLiquipediaApiClient(internal.LiquipediaApiClientParams{
						Query: r.Params,
						Game:  r.Game,
					})

					resp, err := client.Do()
					r.ResponseCh <- LiquipediaResponse{
						Body:  resp,
						Error: err,
					}

					close(r.ResponseCh)
				}
			}
		}
	}
}

func (lb *LiquipediaBus) AddRequest(r *LiquipediaRequest) {
	lb.mutex.Lock()
	lb.requests = append(lb.requests, r)
	lb.mutex.Unlock()
}

func (lb *LiquipediaBus) extract() *LiquipediaRequest {
	if len(lb.requests) == 0 {
		return nil
	}

	item := lb.requests[:1]
	lb.requests = lb.requests[1:]

	return item[0]
}
