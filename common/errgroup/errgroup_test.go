package errgroup

import (

	// "errors"

	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestErrGroup(t *testing.T) {
	g := WithCancel(context.Background())
	g.GOMAXPROCS(5)
	// // task1
	g.Go(func(ctx context.Context) error {
		time.Sleep(time.Second * 1)
		return errors.New("err")
	})
	// // task2
	g.Go(func(ctx context.Context) error {
		time.Sleep(time.Second * 1)
		select {
		case <-ctx.Done():
			fmt.Println("task is cancel")
		default:
			fmt.Println("task success")
		}
		return nil
	})
	g.Go(func(ctx context.Context) error {
		time.Sleep(time.Second * 1)
		select {
		case <-ctx.Done():
			fmt.Println("task is cancel")
		default:
			fmt.Println("task success")
		}
		return nil
	})
	g.Go(func(ctx context.Context) error {
		time.Sleep(time.Second * 1)
		select {
		case <-ctx.Done():
			fmt.Println("task is cancel")
		default:
			fmt.Println("task success")
		}
		return nil
	})
	g.Go(func(ctx context.Context) error {
		time.Sleep(time.Second * 1)
		select {
		case <-ctx.Done():
			fmt.Println("task is cancel")
		default:
			fmt.Println("task success")
		}
		return nil
	})

	err := g.Wait()
	if err != nil {
		t.Error(err)
	}
	// ctx1, _ := context.WithCancel(context.Background())

	// ctx2, _ := context.WithCancel(ctx1)

	// ctx3, cancel4 := context.WithCancel(ctx2)

	// cancel4()

	// fmt.Println(ctx3.Err())
	// fmt.Println(ctx2.Err())
	// fmt.Println(ctx1.Err())

}
