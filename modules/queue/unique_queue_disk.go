// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package queue

import (
	"github.com/sitename/sitename/modules/nosql"

	"gitea.com/lunny/levelqueue"
)

// LevelUniqueQueueType is the type for level queue
const LevelUniqueQueueType Type = "unique-level"

// LevelUniqueQueueConfiguration is the configuration for a LevelUniqueQueue
type LevelUniqueQueueConfiguration struct {
	ByteFIFOQueueConfiguration
	DataDir          string
	ConnectionString string
	QueueName        string
}

// LevelUniqueQueue implements a disk library queue
type LevelUniqueQueue struct {
	*ByteFIFOUniqueQueue
}

// NewLevelUniqueQueue creates a ledis local queue
//
// Please note that this Queue does not guarantee that a particular
// task cannot be processed twice or more at the same time. Uniqueness is
// only guaranteed whilst the task is waiting in the queue.
func NewLevelUniqueQueue(handle HandlerFunc, cfg, exemplar interface{}) (Queue, error) {
	configInterface, err := toConfig(LevelUniqueQueueConfiguration{}, cfg)
	if err != nil {
		return nil, err
	}
	config := configInterface.(LevelUniqueQueueConfiguration)

	if len(config.ConnectionString) == 0 {
		config.ConnectionString = config.DataDir
	}

	byteFIFO, err := NewLevelUniqueQueueByteFIFO(config.ConnectionString, config.QueueName)
	if err != nil {
		return nil, err
	}

	byteFIFOQueue, err := NewByteFIFOUniqueQueue(LevelUniqueQueueType, byteFIFO, handle, config.ByteFIFOQueueConfiguration, exemplar)
	if err != nil {
		return nil, err
	}

	queue := &LevelUniqueQueue{
		ByteFIFOUniqueQueue: byteFIFOQueue,
	}
	queue.qid = GetManager().Add(queue, LevelUniqueQueueType, config, exemplar)
	return queue, nil
}

var _ UniqueByteFIFO = &LevelUniqueQueueByteFIFO{}

// LevelUniqueQueueByteFIFO represents a ByteFIFO formed from a LevelUniqueQueue
type LevelUniqueQueueByteFIFO struct {
	internal   *levelqueue.UniqueQueue
	connection string
}

// NewLevelUniqueQueueByteFIFO creates a new ByteFIFO formed from a LevelUniqueQueue
func NewLevelUniqueQueueByteFIFO(connection, prefix string) (*LevelUniqueQueueByteFIFO, error) {
	db, err := nosql.GetManager().GetLevelDB(connection)
	if err != nil {
		return nil, err
	}

	internal, err := levelqueue.NewUniqueQueue(db, []byte(prefix), []byte(prefix+"-unique"), false)
	if err != nil {
		return nil, err
	}

	return &LevelUniqueQueueByteFIFO{
		connection: connection,
		internal:   internal,
	}, nil
}

// PushFunc pushes data to the end of the fifo and calls the callback if it is added
func (fifo *LevelUniqueQueueByteFIFO) PushFunc(data []byte, fn func() error) error {
	return fifo.internal.LPushFunc(data, fn)
}

// Pop pops data from the start of the fifo
func (fifo *LevelUniqueQueueByteFIFO) Pop() ([]byte, error) {
	data, err := fifo.internal.RPop()
	if err != nil && err != levelqueue.ErrNotFound {
		return nil, err
	}
	return data, nil
}

// Len returns the length of the fifo
func (fifo *LevelUniqueQueueByteFIFO) Len() int64 {
	return fifo.internal.Len()
}

// Has returns whether the fifo contains this data
func (fifo *LevelUniqueQueueByteFIFO) Has(data []byte) (bool, error) {
	return fifo.internal.Has(data)
}

// Close this fifo
func (fifo *LevelUniqueQueueByteFIFO) Close() error {
	err := fifo.internal.Close()
	_ = nosql.GetManager().CloseLevelDB(fifo.connection)
	return err
}

func init() {
	queuesMap[LevelUniqueQueueType] = NewLevelUniqueQueue
}
