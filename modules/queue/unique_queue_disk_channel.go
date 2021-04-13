// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package queue

import (
	"context"
	"sync"
	"time"

	"github.com/sitename/sitename/modules/log"
)

// PersistableChannelUniqueQueueType is the type for persistable queue
const PersistableChannelUniqueQueueType Type = "unique-persistable-channel"

// PersistableChannelUniqueQueueConfiguration is the configuration for a PersistableChannelUniqueQueue
type PersistableChannelUniqueQueueConfiguration struct {
	Name         string
	DataDir      string
	BatchLength  int
	QueueLength  int
	Timeout      time.Duration
	MaxAttempts  int
	Workers      int
	MaxWorkers   int
	BlockTimeout time.Duration
	BoostTimeout time.Duration
	BoostWorkers int
}

// PersistableChannelUniqueQueue wraps a channel queue and level queue together
//
// Please note that this Queue does not guarantee that a particular
// task cannot be processed twice or more at the same time. Uniqueness is
// only guaranteed whilst the task is waiting in the queue.
type PersistableChannelUniqueQueue struct {
	*ChannelUniqueQueue
	delayedStarter
	lock   sync.Mutex
	closed chan struct{}
}

// NewPersistableChannelUniqueQueue creates a wrapped batched channel queue with persistable level queue backend when shutting down
// This differs from a wrapped queue in that the persistent queue is only used to persist at shutdown/terminate
func NewPersistableChannelUniqueQueue(handle HandlerFunc, cfg, exemplar interface{}) (Queue, error) {
	configInterface, err := toConfig(PersistableChannelUniqueQueueConfiguration{}, cfg)
	if err != nil {
		return nil, err
	}
	config := configInterface.(PersistableChannelUniqueQueueConfiguration)

	channelUniqueQueue, err := NewChannelUniqueQueue(handle, ChannelUniqueQueueConfiguration{
		WorkerPoolConfiguration: WorkerPoolConfiguration{
			QueueLength:  config.QueueLength,
			BatchLength:  config.BatchLength,
			BlockTimeout: config.BlockTimeout,
			BoostTimeout: config.BoostTimeout,
			BoostWorkers: config.BoostWorkers,
			MaxWorkers:   config.MaxWorkers,
		},
		Workers: config.Workers,
		Name:    config.Name + "-channel",
	}, exemplar)
	if err != nil {
		return nil, err
	}

	// the level backend only needs temporary workers to catch up with the previously dropped work
	levelCfg := LevelUniqueQueueConfiguration{
		ByteFIFOQueueConfiguration: ByteFIFOQueueConfiguration{
			WorkerPoolConfiguration: WorkerPoolConfiguration{
				QueueLength:  config.QueueLength,
				BatchLength:  config.BatchLength,
				BlockTimeout: 0,
				BoostTimeout: 0,
				BoostWorkers: 0,
				MaxWorkers:   1,
			},
			Workers: 1,
			Name:    config.Name + "-level",
		},
		DataDir: config.DataDir,
	}

	queue := &PersistableChannelUniqueQueue{
		ChannelUniqueQueue: channelUniqueQueue.(*ChannelUniqueQueue),
		closed:             make(chan struct{}),
	}

	levelQueue, err := NewLevelUniqueQueue(func(data ...Data) {
		for _, datum := range data {
			err := queue.Push(datum)
			if err != nil && err != ErrAlreadyInQueue {
				log.Error("Unable push to channelled queue: %v", err)
			}
		}
	}, levelCfg, exemplar)
	if err == nil {
		queue.delayedStarter = delayedStarter{
			internal: levelQueue.(*LevelUniqueQueue),
			name:     config.Name,
		}

		_ = GetManager().Add(queue, PersistableChannelUniqueQueueType, config, exemplar)
		return queue, nil
	}
	if IsErrInvalidConfiguration(err) {
		// Retrying ain't gonna make this any better...
		return nil, ErrInvalidConfiguration{cfg: cfg}
	}

	queue.delayedStarter = delayedStarter{
		cfg:         levelCfg,
		underlying:  LevelUniqueQueueType,
		timeout:     config.Timeout,
		maxAttempts: config.MaxAttempts,
		name:        config.Name,
	}
	_ = GetManager().Add(queue, PersistableChannelUniqueQueueType, config, exemplar)
	return queue, nil
}

// Name returns the name of this queue
func (q *PersistableChannelUniqueQueue) Name() string {
	return q.delayedStarter.name
}

// Push will push the indexer data to queue
func (q *PersistableChannelUniqueQueue) Push(data Data) error {
	return q.PushFunc(data, nil)
}

// PushFunc will push the indexer data to queue
func (q *PersistableChannelUniqueQueue) PushFunc(data Data, fn func() error) error {
	select {
	case <-q.closed:
		return q.internal.(UniqueQueue).PushFunc(data, fn)
	default:
		return q.ChannelUniqueQueue.PushFunc(data, fn)
	}
}

// Has will test if the queue has the data
func (q *PersistableChannelUniqueQueue) Has(data Data) (bool, error) {
	// This is more difficult...
	has, err := q.ChannelUniqueQueue.Has(data)
	if err != nil || has {
		return has, err
	}
	q.lock.Lock()
	defer q.lock.Unlock()
	if q.internal == nil {
		return false, nil
	}
	return q.internal.(UniqueQueue).Has(data)
}

// Run starts to run the queue
func (q *PersistableChannelUniqueQueue) Run(atShutdown, atTerminate func(context.Context, func())) {
	log.Debug("PersistableChannelUniqueQueue: %s Starting", q.delayedStarter.name)

	q.lock.Lock()
	if q.internal == nil {
		err := q.setInternal(atShutdown, func(data ...Data) {
			for _, datum := range data {
				err := q.Push(datum)
				if err != nil && err != ErrAlreadyInQueue {
					log.Error("Unable push to channelled queue: %v", err)
				}
			}
		}, q.exemplar)
		q.lock.Unlock()
		if err != nil {
			log.Fatal("Unable to create internal queue for %s Error: %v", q.Name(), err)
			return
		}
	} else {
		q.lock.Unlock()
	}
	atShutdown(context.Background(), q.Shutdown)
	atTerminate(context.Background(), q.Terminate)

	// Just run the level queue - we shut it down later
	go q.internal.Run(func(_ context.Context, _ func()) {}, func(_ context.Context, _ func()) {})

	go func() {
		_ = q.ChannelUniqueQueue.AddWorkers(q.workers, 0)
	}()

	log.Trace("PersistableChannelUniqueQueue: %s Waiting til closed", q.delayedStarter.name)
	<-q.closed
	log.Trace("PersistableChannelUniqueQueue: %s Cancelling pools", q.delayedStarter.name)
	q.internal.(*LevelUniqueQueue).cancel()
	q.ChannelUniqueQueue.cancel()
	log.Trace("PersistableChannelUniqueQueue: %s Waiting til done", q.delayedStarter.name)
	q.ChannelUniqueQueue.Wait()
	q.internal.(*LevelUniqueQueue).Wait()
	// Redirect all remaining data in the chan to the internal channel
	go func() {
		log.Trace("PersistableChannelUniqueQueue: %s Redirecting remaining data", q.delayedStarter.name)
		for data := range q.ChannelUniqueQueue.dataChan {
			_ = q.internal.Push(data)
		}
		log.Trace("PersistableChannelUniqueQueue: %s Done Redirecting remaining data", q.delayedStarter.name)
	}()
	log.Trace("PersistableChannelUniqueQueue: %s Done main loop", q.delayedStarter.name)
}

// Flush flushes the queue
func (q *PersistableChannelUniqueQueue) Flush(timeout time.Duration) error {
	return q.ChannelUniqueQueue.Flush(timeout)
}

// Shutdown processing this queue
func (q *PersistableChannelUniqueQueue) Shutdown() {
	log.Trace("PersistableChannelUniqueQueue: %s Shutting down", q.delayedStarter.name)
	q.lock.Lock()
	defer q.lock.Unlock()
	select {
	case <-q.closed:
	default:
		if q.internal != nil {
			q.internal.(*LevelUniqueQueue).Shutdown()
		}
		close(q.closed)
	}
	log.Debug("PersistableChannelUniqueQueue: %s Shutdown", q.delayedStarter.name)
}

// Terminate this queue and close the queue
func (q *PersistableChannelUniqueQueue) Terminate() {
	log.Trace("PersistableChannelUniqueQueue: %s Terminating", q.delayedStarter.name)
	q.Shutdown()
	q.lock.Lock()
	defer q.lock.Unlock()
	if q.internal != nil {
		q.internal.(*LevelUniqueQueue).Terminate()
	}
	log.Debug("PersistableChannelUniqueQueue: %s Terminated", q.delayedStarter.name)
}

func init() {
	queuesMap[PersistableChannelUniqueQueueType] = NewPersistableChannelUniqueQueue
}
