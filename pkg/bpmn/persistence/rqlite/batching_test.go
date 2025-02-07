package rqlite

import (
	"bytes"
	"log"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/rqlite/rqlite/v8/store"
	"github.com/stretchr/testify/assert"
)

var persistence *BpmnEnginePersistenceRqlite
var executedQueries []string
var logBuf bytes.Buffer

func TestMain(m *testing.M) {
	// Setup logging
	log.SetOutput(&logBuf)
	defer log.SetOutput(os.Stderr)

	setup()
	os.Exit(m.Run())
}

func setup() {

	// Mock storage
	var executedMutex sync.Mutex

	// Mock execute function
	mockExecute := func(queries []string, conn *store.Store) {
		executedMutex.Lock()
		executedQueries = append(executedQueries, queries...)
		executedMutex.Unlock()
	}

	// Create an instance with the mock function
	persistence = &BpmnEnginePersistenceRqlite{
		queryChan:   make(chan Query, 10),
		flushChan:   make(chan int64),
		executeFunc: mockExecute, // Inject the mock function
		ctx:         &RqliteContext{Str: nil},
	}

	// Start the query processor
	go persistence.queryProcessor()
}

func Test_Batching_Flush_should_execute_queries(t *testing.T) {

	// Send queries
	persistence.queryChan <- Query{Key: 1, SQL: "INSERT INTO test_table VALUES (1)"}
	persistence.queryChan <- Query{Key: 1, SQL: "INSERT INTO test_table VALUES (2)"}

	// Allow time for queries to be queued
	time.Sleep(500 * time.Millisecond)

	// Flush queries
	persistence.flushChan <- 1

	// Allow time for execution
	time.Sleep(500 * time.Millisecond)

	// Verify that all queries were executed
	assert.Equal(t, 1, len(executedQueries), "Expected 1 batch to be executed")
	assert.Contains(t, executedQueries[0], "INSERT INTO test_table VALUES (1)", "Expected query to be executed")
	assert.Contains(t, executedQueries[0], "INSERT INTO test_table VALUES (2)", "Expected query to be executed")
}

func Test_Batching_emptyFlush_should_execte_nothing(t *testing.T) {
	persistence.flushChan <- 1

	// Allow time for execution
	time.Sleep(500 * time.Millisecond)

	// Check log output
	assert.Contains(t, logBuf.String(), "DEBUG No queries to flush.", "Expected log message to be printed")

}
