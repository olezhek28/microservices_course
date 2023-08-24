package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	ticker := time.NewTicker(1 * time.Second)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("done")
			return
		case <-ticker.C:
			fmt.Println("tick")
		}
	}
}
