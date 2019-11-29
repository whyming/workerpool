package workerpool

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestTryTimes(t *testing.T) {
	var a, b, c int64
	geta := func() int64 {
		a++
		return a
	}

	job1 := func() error {
		b = 100
		return nil
	}

	job2 := func() error {
		if geta() > 2 {
			return nil
		}
		return errors.New("xxxx")
	}

	job3 := func() error {
		c = 100
		return nil
	}
	ctx := context.Background()
	g := NewPool(2, 5).NewGroup().SetConf(&Conf{
		RTimes:   3,
		Interval: 100 * time.Millisecond,
		PanicHD: func(e interface{}) {
			panic(e)
		},
	})
	g.AddJob(ctx, job1)
	g.AddJob(ctx, job2)
	g.AddJob(ctx, job3)
	g.Done()
	if a != 3 {
		t.Error("not retry")
	}
	if b != 100 {
		t.Error("b is not 100")
	}
	if c != 100 {
		t.Error("b is not 100")
	}
}
