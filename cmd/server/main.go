package main

import (
	"context"
	"os"
)

const (
	ExitOK = iota
	ExitError
)

func main() {
	if err := run(context.Background()); err != nil {
		os.Exit(ExitError)
	}
	os.Exit(ExitOK)
}

func run(ctx context.Context) error {
	return nil
}