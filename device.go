package main

import (
	"context"
)

type Device struct {
	ID      string
	Name    string
	Actions []Action
}

type Action struct {
	Name string
	Icon string
	ID   int
	Act  func(ctx context.Context) error
}
