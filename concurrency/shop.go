package concurrency

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

type itemId string

type uerReq struct {
	Name     string `json:"name"`
	Item     itemId `json:"item"`
	Quantity uint   `json:"quantity"`
	Price    uint   `json:"price"`
}
type itemsList []string

var Items = itemsList{"phone", "laptop", "keyboard"}

func RequestGenerator(ctx context.Context, n int, name string, items itemsList, ch1 chan uerReq) {
	for i := 0; i < n; i++ {
		newRequest := uerReq{
			Name:     name,
			Item:     itemId(items[rand.Intn(len(items))]),
			Quantity: uint(rand.Intn(5)),
			Price:    uint(rand.Intn(100)),
		}
		fmt.Printf("sum:%d\n", (newRequest.Price * newRequest.Quantity))
		ch1 <- newRequest
	}
	close(ch1)
}

func ProcessesRequests(ctx context.Context, ch1 chan uerReq) {
	var total int
	for newRequest := range ch1 {
		total += int(newRequest.Price * newRequest.Quantity)

	}

	fmt.Println("total:", total)

}

func ShopProg() {
	// Створення контексту з дедлайном 2 секунди
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	ch1 := make(chan uerReq)
	done := make(chan bool)

	go RequestGenerator(ctx, 30, "name", Items, ch1)
	go func() {
		ProcessesRequests(ctx, ch1)
		done <- false
	}()

	select {
	case <-ctx.Done():

		err := ctx.Err()
		if err == context.DeadlineExceeded {
			fmt.Println("Операція не встигла завершитися протягом дедлайну")
		} else if err == context.Canceled {
			fmt.Println("Операція була скасована")
		}
	case <-done:
	}
}
