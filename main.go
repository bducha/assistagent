package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bducha/assistagent/agent"
)

// Dev consts
const (
	OBJECT_ID = "pc-fixe"
)

func main() {

	agent := agent.NewAgent()

	fmt.Println("AssistAgent started")
	fmt.Println()

	


	ctx, cancel := context.WithCancel(context.Background())

	exitChan := make(chan struct{})

	go func() {
		agent.Start(ctx)
		close(exitChan)
	}()

	signalChan := make(chan os.Signal, 1)

	signal.Notify(
		signalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)

	// waiting for signal to stop the agent
	<-signalChan
	fmt.Println("Stopping...")
	cancel()
	// waiting for agent to stop
	<-exitChan
}