package core

import (
	"context"
	"fmt"
	box "github.com/sagernet/sing-box"
	"github.com/sagernet/sing-box/option"
	"net/netip"
)

func Server(conf Peer) error {
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
			Inbounds: []option.Inbound{
				{
					Type: "vless",
					Tag:  "vless-in",
					VLESSOptions: option.VLESSInboundOptions{
						ListenOptions: option.ListenOptions{
							Listen:     option.NewListenAddress(netip.MustParseAddr(conf.Addr)),
							ListenPort: conf.Port,
						},
						Users: []option.VLESSUser{
							{
								Name: "gpp",
								UUID: conf.UUID,
							},
						},
						Transport: &option.V2RayTransportOptions{
							Type: "ws",
							WebsocketOptions: option.V2RayWebsocketOptions{
								Path:                fmt.Sprintf("/%s", conf.UUID),
								Headers:             nil,
								MaxEarlyData:        2048,
								EarlyDataHeaderName: "Sec-WebSocket-Protocol",
							},
						},
						Multiplex: &option.InboundMultiplexOptions{
							Enabled: true,
							Padding: false,
							Brutal:  nil,
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
