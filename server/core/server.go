package core

import (
	"context"
	"fmt"
	box "github.com/sagernet/sing-box"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common/auth"
	"net/netip"
)

func Server(conf Peer) error {
	var in option.Inbound
	switch conf.Protocol {
	case "shadowsocks":
		in = option.Inbound{
			Type: "shadowsocks",
			Tag:  "ss-in",
			ShadowsocksOptions: option.ShadowsocksInboundOptions{
				ListenOptions: option.ListenOptions{
					Listen:     option.NewListenAddress(netip.MustParseAddr(conf.Addr)),
					ListenPort: conf.Port,
				},
				Method:   "xchacha20-ietf-poly1305",
				Password: conf.UUID,
				Multiplex: &option.InboundMultiplexOptions{
					Enabled: true,
					Padding: false,
					Brutal:  nil,
				},
			},
		}
	case "socks":
		in = option.Inbound{
			Type: "socks",
			Tag:  "socks-in",
			SocksOptions: option.SocksInboundOptions{
				ListenOptions: option.ListenOptions{
					Listen:     option.NewListenAddress(netip.MustParseAddr(conf.Addr)),
					ListenPort: conf.Port,
				},
				Users: []auth.User{
					{
						Username: "gpp",
						Password: conf.UUID,
					},
				},
			},
		}
	default:
		in = option.Inbound{
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
		}
	}
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
			Inbounds: []option.Inbound{in},
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
