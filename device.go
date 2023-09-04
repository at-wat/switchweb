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
	ID   string
	Act  func(ctx context.Context) error
}
