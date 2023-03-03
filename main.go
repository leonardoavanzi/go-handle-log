package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

func main() {
	start := time.Now()
	userId := 200
	ctx := context.WithValue(context.Background(), "userId", userId)

	value, err := fetchUserData(ctx, userId)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("result:", value)
	fmt.Println("took:", time.Since(start))

}

type Response struct {
	value int
	err   error
}

func fetchUserData(ctx context.Context, userId int) (int, error) {
	userData := ctx.Value("userId")
	fmt.Println(userData)
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*200)
	defer cancel()

	res := make(chan Response)

	// add go routines
	go func() {
		value, err := fetchDataFromApiWhichCanBeSlow()
		res <- Response{
			value: value,
			err:   err,
		}

	}()

	for {
		select {
		case <-ctx.Done():
			return 0, fmt.Errorf("fetching data from third party tooks to long...")
		case resp := <-res:
			return resp.value, resp.err
		}
	}

}

func fetchDataFromApiWhichCanBeSlow() (int, error) {
	time.Sleep(time.Millisecond * 500) // error
	time.Sleep(time.Millisecond * 150) // ok

	return 777, nil
}
