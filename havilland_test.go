package havilland

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestFlight(t *testing.T) {
	var a Airline
	v, err := a.Fly("LH492", func() (interface{}, error) {
		return "FRA/YVR", nil
	})
	if v.(string) != "FRA/YVR" {
		t.Errorf("Wrong value returned. %q != %q", v, "FRA/YVR")
	}
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestFlightError(t *testing.T) {
	var a Airline
	stop := errors.New("Need extra stop")
	v, err := a.Fly("LH492", func() (interface{}, error) {
		return nil, stop
	})
	if err != stop {
		t.Errorf("Wrong error returned. %v != %v ", err, stop)
	}
	if v != nil {
		t.Errorf("Return value is nil. %#v", v)
	}
}

func TestSingleFlight(t *testing.T) {
	var a Airline
	var flights int32
	c := make(chan string)

	fn := func() (interface{}, error) {
		atomic.AddInt32(&flights, 1)
		return <-c, nil
	}

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			v, err := a.Fly("LH492", fn)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if v.(string) != "FRA/YVR" {
				t.Errorf("Wrong value returned. %q != %q", v, "FRA/YVR")
			}
			wg.Done()
		}()
	}

	time.Sleep(100 * time.Millisecond)
	c <- "FRA/YVR"
	wg.Wait()

	if got := atomic.LoadInt32(&flights); got != 1 {
		t.Errorf("More than one flight has been run: %d", got)
	}
}
