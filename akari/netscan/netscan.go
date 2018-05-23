package netscan

import (
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/1lann/rpc"
	"github.com/Sirupsen/logrus"
	"github.com/pul-s4r/vivid18/akari/geo"
)

const listenPort = "5555"

type Result struct {
	crowd  *geo.Map
	update time.Time
}

type Receiver struct {
	listener net.Listener
	ticker   *time.Ticker

	resultMutex *sync.Mutex
	results     [4]*Result

	translations []*geo.Point
}

func Receive(logger *logrus.Logger, trans []*geo.Point) (*Receiver, error) {
	listener, err := net.Listen("tcp", "0.0.0.0:5555")
	if err != nil {
		return nil, err
	}

	r := &Receiver{
		listener:     listener,
		ticker:       time.NewTicker(500 * time.Millisecond),
		resultMutex:  new(sync.Mutex),
		translations: trans,
	}

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				logger.WithError(err).Fatal("failed to listen")
			}

			client, err := rpc.NewClient(conn)
			if err != nil {
				logger.WithError(err).Warn("error receiving connection")
				continue
			}

			for i := 1; i <= 4; i++ {
				id := strconv.Itoa(i)
				idN := i
				client.On("scan-"+id, func(data interface{}) interface{} {
					defer func() {
						if r := recover(); r != nil {
							logger.WithField("recover", r).Error("recovered from scan panic")
						}
					}()

					r.onReceive(idN, data.(*geo.Map))

					return nil
				})
			}

		}
	}()

	return r, nil
}

func (r *Receiver) onReceive(id int, crowd *geo.Map) {
	r.resultMutex.Lock()
	defer r.resultMutex.Unlock()

	r.results[id] = &Result{
		crowd:  crowd,
		update: time.Now(),
	}
}

func (r *Receiver) ScanPeople(crowd *geo.Map) {
	<-r.ticker.C

	crowd.Lock()
	defer crowd.Unlock()

	crowd.Clear()

	r.resultMutex.Lock()
	for id, result := range r.results {
		if result == nil {
			continue
		}

		if time.Since(result.update) < 3*time.Second {
			crowd.Merge(result.crowd, r.translations[id])
		}
	}
}
