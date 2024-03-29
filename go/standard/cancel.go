package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	fmt.Println("exec sleep 5")

	err := exec.CommandContext(ctx, "sleep", "5").Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("sleep 5 finished")
}

// import (
// 	"context"
// 	"fmt"
// 	"time"
// )

// func main() {
// 	// ctx, cancel := context.WithCancel(context.Background())
// 	// ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
// 	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(2*time.Second))
// 	defer cancel()

// 	go fun(ctx)

// 	time.Sleep(5 * time.Second)
// 	// cancel()
// 	// time.Sleep(5 * time.Second)
// }

// func fun(ctx context.Context) {
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			fmt.Println("canceled")
// 			return
// 		default:
// 			fmt.Println("working")
// 			time.Sleep(1 * time.Second)
// 		}
// 	}
// }
