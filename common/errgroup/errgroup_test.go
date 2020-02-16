package errgroup

import (

	// "errors"

	"testing"
)

func TestErrGroup(t *testing.T) {
	// g := WithCancel(context.Background())
	// g.GOMAXPROCS(1)
	// // task1
	// g.Go(func(ctx context.Context) error {
	// 	time.Sleep(time.Second)
	// 	return errors.New("err")
	// })
	// // task2
	// g.Go(func(ctx context.Context) error {
	// 	time.Sleep(time.Second)
	// 	fmt.Println("test0")
	// 	time.Sleep(time.Millisecond * 300)
	// 	fmt.Println("test1")

	// 	time.Sleep(time.Millisecond * 300)
	// 	fmt.Println("test2")
	// 	time.Sleep(time.Millisecond * 300)
	// 	fmt.Println("test3")
	// 	time.Sleep(time.Millisecond * 300)
	// 	fmt.Println("test4")

	// 	select {
	// 	case <-ctx.Done():
	// 		fmt.Println("task is cancel")
	// 	default:
	// 		fmt.Println("task success")
	// 	}
	// 	return nil
	// })

	// err := g.Wait()
	// if err != nil {
	// 	t.Error(err)
	// }
	// ctx1, _ := context.WithCancel(context.Background())

	// ctx2, _ := context.WithCancel(ctx1)

	// ctx3, cancel4 := context.WithCancel(ctx2)

	// cancel4()

	// fmt.Println(ctx3.Err())
	// fmt.Println(ctx2.Err())
	// fmt.Println(ctx1.Err())

}
