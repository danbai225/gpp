package main

import (
	"client/backend/config"
	"context"
	"fmt"
	box "github.com/sagernet/sing-box"
	"github.com/sagernet/sing-box/option"
	"net/netip"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	loadConfig, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	gamePeer := loadConfig.PeerList[0]
	httpPeer := loadConfig.PeerList[0]
	_box, err := Client(gamePeer, httpPeer)
	if err != nil {
		panic(err)
	}
	err = _box.Start()
	if err != nil {
		panic(err)
	}
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	s := <-sigCh
	fmt.Printf("Received signal: %v\n", s)
	fmt.Println("Exiting...")
	os.Exit(0)
}
func Client(gamePeer, httpPeer *config.Peer) (*box.Box, error) {
	fmt.Println(gamePeer.Addr, gamePeer.UUID, gamePeer.Port)
	home, _ := os.UserHomeDir()
	options := box.Options{
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
						Address: "local",
						Detour:  "dns_out",
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
							netip.MustParsePrefix("172.19.225.1/30"),
						},
						AutoRoute:              true,
						StrictRoute:            false,
						EndpointIndependentNat: true,
						UDPTimeout:             option.UDPTimeoutCompat(time.Second * 300),
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
							InboundOptions: option.InboundOptions{
								SniffEnabled: true,
							},
						},
					},
				},
			},
			Route: &option.RouteOptions{
				AutoDetectInterface: true,
				GeoIP: &option.GeoIPOptions{
					Path:           fmt.Sprintf("%s%c%s%c%s", home, os.PathSeparator, ".gpp", os.PathSeparator, "geoip.db"),
					DownloadURL:    "https://github.com/SagerNet/sing-geoip/releases/latest/download/geoip.db",
					DownloadDetour: "direct",
				},
				Geosite: &option.GeositeOptions{
					Path:           fmt.Sprintf("%s%c%s%c%s", home, os.PathSeparator, ".gpp", os.PathSeparator, "geosite.db"),
					DownloadURL:    "https://github.com/SagerNet/sing-geosite/releases/latest/download/geosite.db",
					DownloadDetour: "direct",
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
					},
					{
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
							Server:     gamePeer.Addr,
							ServerPort: gamePeer.Port,
						},
						UUID: gamePeer.UUID,
						Multiplex: &option.OutboundMultiplexOptions{
							Enabled:        true,
							Protocol:       "h2mux",
							MaxConnections: 4,
							MinStreams:     4,
							MaxStreams:     0,
							Padding:        false,
						},
						Transport: &option.V2RayTransportOptions{
							Type: "ws",
							WebsocketOptions: option.V2RayWebsocketOptions{
								Path:                fmt.Sprintf("/%s", gamePeer.UUID),
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
	}
	if httpPeer.Addr != "" {
		out := "http-out"
		if httpPeer.Addr != "direct" {
			options.Options.Outbounds = append(options.Options.Outbounds, option.Outbound{
				Type: "vless",
				Tag:  "http-out",
				VLESSOptions: option.VLESSOutboundOptions{
					ServerOptions: option.ServerOptions{
						Server:     httpPeer.Addr,
						ServerPort: httpPeer.Port,
					},
					UUID: httpPeer.UUID,
					Multiplex: &option.OutboundMultiplexOptions{
						Enabled:        true,
						Protocol:       "h2mux",
						MaxConnections: 4,
						MinStreams:     4,
						MaxStreams:     0,
						Padding:        false,
					},
					Transport: &option.V2RayTransportOptions{
						Type: "ws",
						WebsocketOptions: option.V2RayWebsocketOptions{
							Path:                fmt.Sprintf("/%s", httpPeer.UUID),
							MaxEarlyData:        2048,
							EarlyDataHeaderName: "Sec-WebSocket-Protocol",
						},
					},
				},
			})
		} else {
			out = "direct"
		}
		options.Options.Route.Rules = append(options.Options.Route.Rules, option.Rule{Type: "default", DefaultOptions: option.DefaultRule{Protocol: option.Listable[string]{"http", "quic", "tls"}, Outbound: out}})
	}
	var instance, err = box.New(options)
	if err != nil {
		return nil, err
	}
	return instance, nil
}
