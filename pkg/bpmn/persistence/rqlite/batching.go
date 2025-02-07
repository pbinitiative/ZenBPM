package rqlite

import (
	"github.com/pbinitiative/zenbpm/internal/log"
	"github.com/rqlite/rqlite/v8/store"
)

type Query struct {
	Key int64
	SQL string
}

// queryProcessor listens for queries and flush triggers
func (persistence *BpmnEnginePersistenceRqlite) queryProcessor() {
	data := make(map[int64]string)

	for {
		select {
		case query := <-persistence.queryChan:
			// Append query to batch
			persistence.mutex.Lock()
			data[query.Key] += query.SQL
			persistence.mutex.Unlock()

		case key := <-persistence.flushChan:
			// Execute batched queries
			persistence.mutex.Lock()
			queue := data[key]
			if len(queue) > 0 {
				persistence.executeFunc([]string{queue}, persistence.ctx.Str)
				delete(data, key)
			} else {
				log.Debug("No queries to flush.")
			}
			persistence.mutex.Unlock()
		}
	}
}

// Flushes the transaction effectively writing the sql commands batch to the database cluster
func (persistence *BpmnEnginePersistenceRqlite) FlushTransaction(key int64) error {
	persistence.flushChan <- key
	return nil
}

// Default execute function (production)
func defaultExecute(queries []string, str *store.Store) {
	// Your actual database execution logic
	execute(queries, str)
}
