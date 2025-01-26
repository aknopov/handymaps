package expiry

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	maxCapacity = 123
	ttl         = time.Duration(777000000)
)

func TestCreation(t *testing.T) {
	assertT := assert.New(t)

	em := NewExpiryMap[string, int]().
		WithMaxCapacity(maxCapacity).
		Expirefter(ttl)

	assertT.NotNil(t, em)
	assertT.Equal(maxCapacity, em.Capacity())
	assertT.Equal(ttl, em.ExpiringAfter())
	assertT.NotNil(em.loader)
	assertT.Equal(0, em.listeners.size())
}

func TestLoader(t *testing.T) {
	assertT := assert.New(t)

	em := NewExpiryMap[string, int]().
		WithMaxCapacity(maxCapacity).
		Expirefter(ttl)
	orgLoader := em.loader
	assertT.NotNil(orgLoader)

	em.WithLoader(func(key string) (int, error) { return len(key), nil })
	assertT.NotEqual(fmt.Sprintf("%v", &orgLoader), fmt.Sprintf("%v", &em.loader))
}

func listener1(ev EventType, key string, val int, err error) {
	fmt.Printf("1: Received event: %v, key=%v, val=%v, err=%v\n", ev, key, val, err)
}

func listener2(ev EventType, key string, val int, err error) {
	fmt.Printf("2: Received event: %v, key=%v, val=%v, err=%v\n", ev, key, val, err)
}

func TestListeners(t *testing.T) {
	assertT := assert.New(t)

	em := NewExpiryMap[string, int]().
		WithMaxCapacity(maxCapacity).
		Expirefter(ttl)

	var wrapper1 = ListenerWarapper{f: listener1}
	var wrapper2 = ListenerWarapper{f: listener2}

	em.AddListener(&wrapper1)
	assertT.Equal(1, em.listeners.size())
	assertT.True(em.listeners.contains(&wrapper1))

	em.AddListener(&wrapper1)
	assertT.Equal(1, em.listeners.size())
	assertT.True(em.listeners.contains(&wrapper1))

	em.AddListener(&wrapper2)
	assertT.Equal(2, em.listeners.size())
	assertT.True(em.listeners.contains(&wrapper2))

	em.RemoveListener(&wrapper1)
	assertT.Equal(1, em.listeners.size())
	assertT.False(em.listeners.contains(&wrapper1))
	assertT.True(em.listeners.contains(&wrapper2))
}


func TestTimers(t *testing.T) {
	timer1 := time.NewTimer(2 * time.Second)
	<-timer1.C
	fmt.Println("Timer 1 fired")

	timer2 := time.NewTimer(time.Second)
	go func() {
		<-timer2.C
		fmt.Println("Timer 2 fired") // <- Will not see
	}()
	stop2 := timer2.Stop()
	if stop2 {
		fmt.Println("Timer 2 stopped")
	}
	time.Sleep(2500 * time.Millisecond)
}

