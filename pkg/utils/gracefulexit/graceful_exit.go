package gracefulexit

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func TerminateApp(ctx context.Context) {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case v := <-quit:
		fmt.Printf("signal.Notify: %v \n", v)
		os.Exit(1)
	case done := <-ctx.Done():
		fmt.Printf("ctx.Done: %v \n", done)
	}

}
