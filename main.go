package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	gholtcontext "github.com/gholt/context"
	"golang.org/x/net/context"
)

const COUNT = 10000000

func main() {
	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)
	beginMemSys := memStats.Sys
	beginMemMallocs := memStats.Mallocs
	switch os.Args[1] {
	case "noTimeout":
		noTimeout()
	case "timeAfter":
		timeAfter()
	case "ticker":
		ticker()
	case "timer":
		timer()
	case "context1":
		context1()
	case "context2":
		context2()
	}
	runtime.ReadMemStats(memStats)
	endMemSys := memStats.Sys
	endMemMallocs := memStats.Mallocs
	fmt.Printf("    %d sys, %d mallocs\n", endMemSys-beginMemSys, endMemMallocs-beginMemMallocs)
}

func noTimeout() {
	begin := time.Now()
	type message struct{}
	messageChan := make(chan *message)
	doneChan := make(chan struct{})
	go func() {
		for {
			if <-messageChan == nil {
				doneChan <- struct{}{}
				return
			}
		}
	}()
	for i := 0; i < COUNT; i++ {
		messageChan <- &message{}
	}
	close(messageChan)
	<-doneChan
	elapsed := time.Since(begin)
	fmt.Printf("noTimeout: %dns/message, %s elapsed\n", int64(elapsed)/COUNT, elapsed)
}

func timeAfter() {
	begin := time.Now()
	type message struct{}
	messageChan := make(chan *message)
	doneChan := make(chan struct{})
	go func() {
		for {
			if <-messageChan == nil {
				doneChan <- struct{}{}
				return
			}
		}
	}()
	var timeouts int
	for i := 0; i < COUNT; i++ {
		select {
		case messageChan <- &message{}:
		case <-time.After(time.Second):
			timeouts++
		}
	}
	close(messageChan)
	<-doneChan
	elapsed := time.Since(begin)
	fmt.Printf("timeAfter: %dns/message, %s elapsed, %d timeouts\n", int64(elapsed)/COUNT, elapsed, timeouts)
}

func ticker() {
	begin := time.Now()
	type message struct{}
	messageChan := make(chan *message)
	doneChan := make(chan struct{})
	go func() {
		for {
			if <-messageChan == nil {
				doneChan <- struct{}{}
				return
			}
		}
	}()
	ticker := time.NewTicker(time.Second)
	var timeouts int
	for i := 0; i < COUNT; i++ {
		for j := 0; j < 2; j++ {
			select {
			case messageChan <- &message{}:
				j = 2
			case <-ticker.C:
				if j == 1 {
					timeouts++
				}
			}
		}
	}
	ticker.Stop()
	close(messageChan)
	<-doneChan
	elapsed := time.Since(begin)
	fmt.Printf("ticker:    %dns/message, %s elapsed, %d timeouts\n", int64(elapsed)/COUNT, elapsed, timeouts)
}

func timer() {
	begin := time.Now()
	type message struct{}
	messageChan := make(chan *message)
	doneChan := make(chan struct{})
	go func() {
		for {
			if <-messageChan == nil {
				doneChan <- struct{}{}
				return
			}
		}
	}()
	timer := time.NewTimer(time.Second)
	if !timer.Stop() {
		<-timer.C
	}
	var timeouts int
	for i := 0; i < COUNT; i++ {
		timer.Reset(time.Second)
		select {
		case messageChan <- &message{}:
			if !timer.Stop() {
				<-timer.C
			}
		case <-timer.C:
			timeouts++
		}
	}
	timer.Stop()
	close(messageChan)
	<-doneChan
	elapsed := time.Since(begin)
	fmt.Printf("timer:     %dns/message, %s elapsed, %d timeouts\n", int64(elapsed)/COUNT, elapsed, timeouts)
}

func context1() {
	begin := time.Now()
	type message struct{}
	messageChan := make(chan *message)
	doneChan := make(chan struct{})
	go func() {
		for {
			if <-messageChan == nil {
				doneChan <- struct{}{}
				return
			}
		}
	}()
	bgctx := context.Background()
	var timeouts int
	for i := 0; i < COUNT; i++ {
		ctx, cancel := context.WithTimeout(bgctx, time.Second)
		select {
		case messageChan <- &message{}:
		case <-ctx.Done():
			timeouts++
		}
		cancel()
	}
	close(messageChan)
	<-doneChan
	elapsed := time.Since(begin)
	fmt.Printf("context1:  %dns/message, %s elapsed, %d timeouts\n", int64(elapsed)/COUNT, elapsed, timeouts)
}

func context2() {
	begin := time.Now()
	type message struct{}
	messageChan := make(chan *message)
	doneChan := make(chan struct{})
	go func() {
		for {
			if <-messageChan == nil {
				doneChan <- struct{}{}
				return
			}
		}
	}()
	ctx := gholtcontext.New(time.Second)
	var timeouts int
	for i := 0; i < COUNT; i++ {
		ctx.Reinit(time.Second)
		select {
		case messageChan <- &message{}:
		case <-ctx.Done():
			timeouts++
		}
	}
	close(messageChan)
	<-doneChan
	elapsed := time.Since(begin)
	fmt.Printf("context2:  %dns/message, %s elapsed, %d timeouts\n", int64(elapsed)/COUNT, elapsed, timeouts)
}
