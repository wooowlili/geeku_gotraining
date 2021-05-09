package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	g "github.com/da440dil/go-workgroup"
	gc "github.com/da440dil/go-workgroup/template/context"
	gsh "github.com/da440dil/go-workgroup/template/shutdown"
	gsig "github.com/da440dil/go-workgroup/template/signal"
)

func main() {
	// Create workgroup
	var wg g.Group
	// Add function to cancel execution using os signal
	wg.Add(gsig.New(os.Interrupt))
	// Create http server
	srv := http.Server{Addr: "127.0.0.1:8080"}
	// Add function to start and stop http server
	wg.Add(gsh.New(
		func() error {
			fmt.Printf("Server is about to listen at %v\n", srv.Addr)
			return srv.ListenAndServe()
		},
		func() {
			fmt.Println("Server is about to shutdown")
			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
			defer cancel()

			err := srv.Shutdown(ctx)
			fmt.Printf("Server shutdown with error: %v\n", err)
		},
	))
	// Create context to cancel execution after 5 seconds
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(time.Second * 5)
		fmt.Println("Context cancel")
		cancel()
	}()
	// Add function to cancel execution using context
	wg.Add(gc.New(ctx))
	// Execute each function
	err := wg.Run()
	fmt.Printf("Workgroup quit with error: %v\n", err)

}
