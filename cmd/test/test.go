package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/valuetechdev/minimal-server-wrapper-go/vserver"
)

func main() {
	server := vserver.New(&vserver.ServerOptions{
		Addr: ":8080",
	})

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	// simulating the container shutting down
	go func() {
		pid := os.Getpid()

		time.Sleep(5 * time.Second)
		syscall.Kill(pid, syscall.SIGTERM)
	}()

	server.AddRoute("GET /", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(3 * time.Second)
		fmt.Fprint(w, "hello world")
		w.WriteHeader(http.StatusOK)
	}))

	server.AddMiddleware(func(h http.Handler) http.Handler {
		fmt.Println("request received")

		return h
	})

	go func() {
		err := server.Serve()

		fmt.Printf("%s\n", err)
	}()

	<-ctx.Done()

	err := server.Shutdown(context.Background())
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	} else {
		fmt.Print("successful shutdown\n")
		os.Exit(0)
	}
}
