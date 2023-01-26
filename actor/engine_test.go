package actor

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type dummy struct{}

func newDummy() Receiver {
	return &dummy{}
}
func (d *dummy) Receive(ctx *Context) {
	switch msg := ctx.Message().(type) {
	case string:
		_ = msg
		fmt.Println("receiving:", ctx.Message())
		ctx.Respond("ksjfkdjfkdjfkdfjkdfjkdf")
	default:
	}
}

func TestSpawn(t *testing.T) {
	e := NewEngine()
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			tag := strconv.Itoa(i)
			pid := e.Spawn(newDummy, "dummy", tag)
			e.Send(pid, 1)
			wg.Done()
		}(i)
	}

	wg.Wait()
}

func TestSpawnPID(t *testing.T) {
	e := NewEngine()

	pid := e.Spawn(newDummy, "dummy", "1")
	assert.Equal(t, "local/dummy/1", pid.String())
}

func TestPoison(t *testing.T) {
	e := NewEngine()

	for i := 0; i < 4; i++ {
		tag := strconv.Itoa(i)
		pid := e.Spawn(newDummy, "dummy", tag)
		e.Poison(pid)
		assert.Nil(t, e.registry.get(pid))
	}
}

func TestXxx(t *testing.T) {
	e := NewEngine()
	pid := e.Spawn(newDummy, "dummy")

	//for i := 0; i < b.N; i++ {
	e.Send(pid, pid)
	//}
}

func BenchmarkSendMessageLocal(b *testing.B) {
	e := NewEngine()
	pid := e.Spawn(newDummy, "dummy")

	for i := 0; i < b.N; i++ {
		e.Send(pid, pid)
	}
}

func TestRequestResponse(t *testing.T) {
	e := NewEngine()
	pid := e.Spawn(newDummy, "dummy")
	resp := e.Request(pid, "foo", time.Millisecond)
	res, err := resp.Result()
	assert.Nil(t, err)
	fmt.Println(resp.PID())
	fmt.Println(res)
	// deadletter
	fmt.Println(e.registry.get(resp.pid).PID())
}
