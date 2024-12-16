package main

import (
	"context"
	"log"
	"os"

	"github.com/nasa9084/go-switchbot/v4"
)

type client struct {
	cli *switchbot.Client
}

func newClient() *client {
	token, ok := os.LookupEnv("SWITCHBOT_TOKEN")
	if !ok {
		log.Fatal("SWITCHBOT_TOKEN not set")
	}
	secret, ok := os.LookupEnv("SWITCHBOT_CLIENT_SECRET")
	if !ok {
		log.Fatal("SWITCHBOT_CLIENT_SECRET not set")
	}

	return &client{
		cli: switchbot.New(token, secret),
	}
}

func (c *client) list(ctx context.Context) []Device {
	var devs []Device

	pdev, idev, err := c.cli.Device().List(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for _, d := range pdev {
		d := d
		log.Printf("%s: [%s] %s\n", d.ID, d.Type, d.Name)
		var acts []Action
		switch d.Type {
		case switchbot.Bot, switchbot.Plug:
			acts = []Action{
				{
					Name: "ON",
					ID:   d.ID + ".on",
					Icon: "plug",
					Act: func(ctx context.Context) error {
						return c.cli.Device().Command(ctx, d.ID, switchbot.TurnOnCommand())
					},
				},
				{
					Name: "OFF",
					ID:   d.ID + ".off",
					Icon: "power-off",
					Act: func(ctx context.Context) error {
						return c.cli.Device().Command(ctx, d.ID, switchbot.TurnOffCommand())
					},
				},
			}
		case switchbot.Curtain:
			acts = []Action{
				{
					Name: "Open",
					ID:   d.ID + ".open",
					Icon: "door-open",
					Act: func(ctx context.Context) error {
						return c.cli.Device().Command(ctx, d.ID, switchbot.TurnOnCommand())
					},
				},
				{
					Name: "Close",
					ID:   d.ID + ".close",
					Icon: "door-closed",
					Act: func(ctx context.Context) error {
						return c.cli.Device().Command(ctx, d.ID, switchbot.TurnOffCommand())
					},
				},
			}
		}
		if len(acts) == 0 {
			continue
		}
		devs = append(devs, Device{
			ID:      d.ID,
			Name:    d.Name,
			Actions: acts,
		})
	}
	for _, d := range idev {
		d := d
		log.Printf("%s: [%s] %s\n", d.ID, d.Type, d.Name)
		var acts []Action
		switch d.Type {
		case switchbot.Others:
		case switchbot.Fan, "DIY Fan":
			acts = []Action{
				{
					Icon: "plug",
					ID:   d.ID + ".on",
					Act: func(ctx context.Context) error {
						return c.cli.Device().Command(ctx, d.ID, switchbot.TurnOnCommand())
					},
				},
				{
					Name: "1",
					ID:   d.ID + ".1",
					Icon: "fan",
					Act: func(ctx context.Context) error {
						return c.cli.Device().Command(ctx, d.ID, switchbot.FanLowSpeedCommand())
					},
				},
				{
					Name: "2",
					ID:   d.ID + ".2",
					Icon: "fan",
					Act: func(ctx context.Context) error {
						return c.cli.Device().Command(ctx, d.ID, switchbot.FanMiddleSpeedCommand())
					},
				},
				{
					Name: "3",
					ID:   d.ID + ".3",
					Icon: "fan",
					Act: func(ctx context.Context) error {
						return c.cli.Device().Command(ctx, d.ID, switchbot.FanHighSpeedCommand())
					},
				},
				{
					Icon: "power-off",
					ID:   d.ID + ".off",
					Act: func(ctx context.Context) error {
						return c.cli.Device().Command(ctx, d.ID, switchbot.TurnOffCommand())
					},
				},
				{
					Icon: "left-right",
					ID:   d.ID + ".lr",
					Act: func(ctx context.Context) error {
						return c.cli.Device().Command(ctx, d.ID, switchbot.FanSwingCommand())
					},
				},
			}
		case switchbot.TV, "DIY TV":
			acts = []Action{
				{
					Icon: "arrow-right-to-bracket",
					ID:   d.ID + ".input",
					Name: "Input",
					Act: func(ctx context.Context) error {
						return c.cli.Device().Command(ctx, d.ID, &switchbot.DeviceCommandRequest{
							Command:     "input",
							Parameter:   "default",
							CommandType: "command",
						})
					},
				},
			}
		}
		if dev, ok := customDevices[d.ID]; ok {
			for _, b := range dev.Buttons {
				b := b
				acts = append(acts, Action{
					Icon: b.Icon,
					ID:   dev.ID + "." + b.Name,
					Name: b.Name,
					Act: func(ctx context.Context) error {
						return c.cli.Device().Command(ctx, dev.ID, &switchbot.DeviceCommandRequest{
							Command:     b.Name,
							Parameter:   "default",
							CommandType: "customize",
						})
					},
				})
			}
		}
		if len(acts) == 0 {
			continue
		}
		devs = append(devs, Device{
			ID:      d.ID,
			Name:    d.Name,
			Actions: acts,
		})
	}
	return devs
}
