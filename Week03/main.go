package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"golang.org/x/sync/errgroup"
)

func main() {
	eg := errgroup.Group{}

	// handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		fmt.Fprintln(w, "hello")
	})

	// server
	srv := http.Server{
		Addr:    "127.0.0.1",
		Handler: handler,
	}

	eg.Go(func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); nil != err {
			log.Fatalf("server shutdown failed, err: %v\n", err)
		}
		log.Println("server gracefully shutdown")
	}())

	// serve
	eg.Go(func() error {
		go func() {
			err := srv.ListenAndServe()
			if http.ErrServerClosed != err {
				log.Fatalf("server not gracefully shutdown, err :%v\n", err)
			}

		}()
	})

	eg.Wait()
}

