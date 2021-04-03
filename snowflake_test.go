package snowflake

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	sf = New(111)
}

func TestUniqueID(t *testing.T) {
	for i := 0; i < 10000000; i++ { // 1000w
		_, err := Next()
		assert.Equal(t, err, nil)
	}
	// 150w/s
}

func BenchmarkGenerateID(t *testing.B) {
	for i := 0; i < t.N; i++ {
		_, err := sf.Next()
		assert.Equal(t, err, nil)
	}

	// 	BenchmarkGenerateID
	//  BenchmarkGenerateID-12    	 1436590	       727 ns/op	     216 B/op	       2 allocs/op
}

func TestBench(t *testing.T) {
	workerNum := 5
	cycle := 100 * 10000 // 100w

	m := make(map[int64]bool, workerNum*cycle)
	lock := sync.Mutex{}
	now := time.Now()

	sf = New(100)

	wg := sync.WaitGroup{}
	for i := 0; i < workerNum; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for i := 0; i < cycle; i++ {
				id, err := Next()
				assert.Equal(t, err, nil)

				lock.Lock()
				_, ok := m[id]
				assert.Equal(t, ok, false)
				m[id] = true
				lock.Unlock()
			}
		}()
	}
	wg.Wait()

	t.Log(time.Since(now).String())
	assert.Equal(t, len(m), workerNum*cycle)
}

func TestDiffTime(t *testing.T) {
	id, err := sf.Next()
	assert.Equal(t, nil, err)
	t.Log(id, sf.GetTimeFromID(id))

	count := 500 * 10000 // 500w
	for i := 0; i < count; i++ {
		nid, err := sf.Next()
		if nid == id {
			t.Error("id cause conflict")
			return
		}

		id = nid
		assert.Equal(t, nil, err)

		ms := GetTimeFromID(id)
		now := int64(time.Now().UnixNano() / 1000 / 1000)

		diff := now - ms
		assert.GreaterOrEqual(t, int64(5), diff) // diff < 5 ms
	}
	t.Log("end")
}
