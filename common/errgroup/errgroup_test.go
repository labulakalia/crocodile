package errgroup

import (
	"context"
	"errors"
	// "errors"
	"time"

	"fmt"
	"testing"
)

func TestErrGroup(t *testing.T) {
	g := WithCancel(context.Background())
	g.GOMAXPROCS(2)
	// task1
	g.Go(func(ctx context.Context) error {
		select {
		case <- ctx.Done():
			return ctx.Err()
		case <- time.After(time.Second * 2):
			return errors.New("timeout")
		}
	})
	// task2
	g.Go(func(ctx context.Context) error {
		var err error
		select {
		case <- ctx.Done():
			err =  ctx.Err()
			goto Check
		case <- time.After(time.Second * 3):
			err = errors.New("timeout")
			goto Check
		}
		Check:
			fmt.Println("error: ",err)
			return err
	})

	err := g.Wait()
	if err == nil {
		t.Error("err should not nil")
	}
}
