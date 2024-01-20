package core

import (
	"context"
	"fmt"
	box "github.com/sagernet/sing-box"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common/auth"
	"net/netip"
	"os"
)

func Client(conf Config) (*box.Box, error) {
	home, _ := os.UserHomeDir()
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
						Detour:  "direct",
					},
				},
				Rules:          []option.DNSRule{},
				Final:          "ali",
				ReverseMapping: false,
				FakeIP:         nil,
				DNSClientOptions: option.DNSClientOptions{
					DisableCache:     false,
					DisableExpire:    false,
					IndependentCache: false,
				},
			},
			Inbounds: []option.Inbound{
				{
					Type: "tun",
					Tag:  "tun-in",
					TunOptions: option.TunInboundOptions{
						InterfaceName: "utun225",
						MTU:           9000,
						Inet4Address: option.Listable[netip.Prefix]{
							netip.MustParsePrefix("172.19.0.1/30"),
						},
						AutoRoute:              true,
						StrictRoute:            false,
						EndpointIndependentNat: true,
						UDPTimeout:             300,
						Stack:                  "system",
						InboundOptions: option.InboundOptions{
							SniffEnabled: true,
						},
					},
				},
				{
					Type: "socks",
					Tag:  "socks-in",
					SocksOptions: option.SocksInboundOptions{
						ListenOptions: option.ListenOptions{
							Listen:     option.NewListenAddress(netip.MustParseAddr("0.0.0.0")),
							ListenPort: 5123,
						},
						Users: []auth.User{
							{
								Username: "admin",
								Password: conf.UUID,
							},
						},
					},
				},
			},
			Route: &option.RouteOptions{
				AutoDetectInterface: true,
				GeoIP: &option.GeoIPOptions{
					Path:        fmt.Sprintf("%s%c%s%c%s", home, os.PathSeparator, ".gpp", os.PathSeparator, "data-a"),
					DownloadURL: "https://ghps.cc/https://github.com/SagerNet/sing-geoip/releases/latest/download/geoip.db",
				},
				Geosite: &option.GeositeOptions{
					Path:        fmt.Sprintf("%s%c%s%c%s", home, os.PathSeparator, ".gpp", os.PathSeparator, "data-b"),
					DownloadURL: "https://ghps.cc/https://github.com/SagerNet/sing-geosite/releases/latest/download/geosite.db",
				},
				Rules: []option.Rule{
					{
						Type: "default",
						DefaultOptions: option.DefaultRule{
							Protocol: option.Listable[string]{"dns"},
							Outbound: "dns_out",
						},
					},
					{
						Type: "default",
						DefaultOptions: option.DefaultRule{
							Inbound:  option.Listable[string]{"dns_in"},
							Outbound: "dns_out",
						},
					},
					{
						Type: "default",
						DefaultOptions: option.DefaultRule{
							Geosite:  option.Listable[string]{"cn"},
							Outbound: "direct",
						},
					}, {
						Type: "default",
						DefaultOptions: option.DefaultRule{
							GeoIP:    option.Listable[string]{"cn", "private"},
							Outbound: "direct",
						},
					},
				},
			},
			Outbounds: []option.Outbound{
				{
					Type: "vless",
					Tag:  "vless-out",
					VLESSOptions: option.VLESSOutboundOptions{
						ServerOptions: option.ServerOptions{
							Server:     conf.Addr,
							ServerPort: conf.Port,
						},
						UUID: conf.UUID,
						Multiplex: &option.MultiplexOptions{
							Enabled:        true,
							Protocol:       "smux",
							MaxConnections: 5,
							MinStreams:     1,
							MaxStreams:     10,
							Padding:        false,
						},
						Transport: &option.V2RayTransportOptions{
							Type: "ws",
							WebsocketOptions: option.V2RayWebsocketOptions{
								Path:                "/test",
								MaxEarlyData:        2048,
								EarlyDataHeaderName: "Sec-WebSocket-Protocol",
							},
						},
					},
				},
				{
					Type: "block",
					Tag:  "block",
				},
				{
					Type: "direct",
					Tag:  "direct",
				}, {
					Type: "dns",
					Tag:  "dns_out",
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return instance, nil
}
