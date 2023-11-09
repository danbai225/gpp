package core

import (
	"context"
	box "github.com/sagernet/sing-box"
	"github.com/sagernet/sing-box/option"
	"net/netip"
	"time"
)

func Server() error {
	var instance, err = box.New(box.Options{
		Context: context.Background(),
		Options: option.Options{
			Log: &option.LogOptions{
				Disabled:     false,
				Level:        "info",
				Output:       "run.log",
				Timestamp:    true,
				DisableColor: true,
			},
			DNS: &option.DNSOptions{
				Servers: []option.DNSServerOptions{
					{
						Tag:     "ali",
						Address: "223.5.5.5",
					},
				},
				Rules:          []option.DNSRule{},
				Final:          "",
				ReverseMapping: false,
				FakeIP:         nil,
				DNSClientOptions: option.DNSClientOptions{
					DisableCache:     false,
					DisableExpire:    false,
					IndependentCache: false,
				},
			},
			NTP: &option.NTPOptions{
				Enabled:       true,
				Interval:      option.Duration(time.Minute * 30),
				WriteToSystem: false,
				ServerOptions: option.ServerOptions{
					Server:     "time.apple.com",
					ServerPort: 123,
				},
				DialerOptions: option.DialerOptions{},
			},
			Inbounds: []option.Inbound{
				{
					Type: "vless",
					Tag:  "vless-in",
					VLESSOptions: option.VLESSInboundOptions{
						ListenOptions: option.ListenOptions{
							Listen:     option.NewListenAddress(netip.AddrFrom4([4]byte([]byte{0, 0, 0, 0}))),
							ListenPort: 5123,
						},
						Users: []option.VLESSUser{
							{
								Name: "danbai",
								UUID: "badb17ef-eb22-4e03-9b17-efeb224e03e7",
							},
						},
						TLS: nil,
						Transport: &option.V2RayTransportOptions{
							Type: "ws",
							WebsocketOptions: option.V2RayWebsocketOptions{
								Path:                "/test",
								Headers:             nil,
								MaxEarlyData:        2048,
								EarlyDataHeaderName: "Sec-WebSocket-Protocol",
							},
						},
					},
				},
			},
			Outbounds: []option.Outbound{
				{
					Type: "direct",
					Tag:  "direct-out",
				},
			},
		},
	})
	if err != nil {
		return err
	}
	err = instance.Start()
	if err != nil {
		return err
	}
	return nil
}
