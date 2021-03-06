package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pkg.agungdwiprasetyo.com/candi/codebase/interfaces"
	"pkg.agungdwiprasetyo.com/candi/config/env"
	"pkg.agungdwiprasetyo.com/candi/logger"
	"pkg.agungdwiprasetyo.com/candi/tracer"
)

const timeout time.Duration = 10 * time.Second

// Config app
type Config struct {
	closers []interfaces.Closer
}

// Init app config
func Init(serviceName string) *Config {
	env.Load(serviceName)
	logger.InitZap()
	tracer.InitOpenTracing()
	return &Config{}
}

// LoadFunc load selected dependency with context timeout
func (c *Config) LoadFunc(depsFunc func(context.Context) []interfaces.Closer) {
	// set timeout for init configuration
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result := make(chan error)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				result <- fmt.Errorf("Failed init configuration :=> %v", r)
			}
			close(result)
		}()

		c.closers = depsFunc(ctx)
	}()

	// with timeout to init configuration
	select {
	case <-ctx.Done():
		panic(fmt.Errorf("Timeout to load selected dependencies: %v", ctx.Err()))
	case err := <-result:
		if err != nil {
			panic(err)
		}
		return
	}
}

// Exit close all connection
func (c *Config) Exit() {
	// set timeout for close all configuration
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	fmt.Println()

	errCloseChan := make(chan error)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				errCloseChan <- fmt.Errorf("Failed close connection :=> %v", r)
			}
			close(errCloseChan)
		}()

		for _, cl := range c.closers {
			cl.Disconnect(ctx)
		}
	}()

	// for force exit
	quitSignal := make(chan os.Signal, 1)
	signal.Notify(quitSignal, os.Interrupt, syscall.SIGTERM)

	// with timeout to close all configuration
	select {
	case <-quitSignal:
		fmt.Println("\x1b[31;1mForce exit\x1b[0m")
	case <-ctx.Done():
		panic(fmt.Errorf("Timeout to close all selected dependencies connection: %v", ctx.Err()))
	case err := <-errCloseChan:
		if err != nil {
			panic(err)
		}
		log.Println("\x1b[32;1mSuccess close all config dependency\x1b[0m")
	}
}
