package mdata

import (
	"time"

	"github.com/raintank/raintank-metric/metric_tank/mdata/chunk"
)

type ChunkReadRequest struct {
	month   uint32
	sortKey uint32
	q       string
	p       []interface{}
	out     chan outcome
}

type ChunkWriteRequest struct {
	key       string
	chunk     *chunk.Chunk
	ttl       uint32
	timestamp time.Time
}
