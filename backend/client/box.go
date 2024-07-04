package client

import (
	"context"
	"fmt"
	"github.com/danbai225/gpp/backend/config"
	"github.com/google/uuid"
	box "github.com/sagernet/sing-box"
	"github.com/sagernet/sing-box/option"
	dns "github.com/sagernet/sing-dns"
	"net/netip"
	"os"
	"time"
)

func getOUt(peer *config.Peer) option.Outbound {
	var out option.Outbound
	switch peer.Protocol {
	case "shadowsocks":
		out = option.Outbound{
			Type: "shadowsocks",
			ShadowsocksOptions: option.ShadowsocksOutboundOptions{
				ServerOptions: option.ServerOptions{
					Server:     peer.Addr,
					ServerPort: peer.Port,
				},
				Method:   "xchacha20-ietf-poly1305",
				Password: peer.UUID,
			},
		}
	case "socks":
		out = option.Outbound{
			Type: "socks",
			SocksOptions: option.SocksOutboundOptions{
				ServerOptions: option.ServerOptions{
					Server:     peer.Addr,
					ServerPort: peer.Port,
				},
				Username: "gpp",
				Password: peer.UUID,
				UDPOverTCP: &option.UDPOverTCPOptions{
					Enabled: true,
					Version: 2,
				},
			},
		}
	case "hysteria2":
		out = option.Outbound{
			Type: "hysteria2",
			Hysteria2Options: option.Hysteria2OutboundOptions{
				ServerOptions: option.ServerOptions{
					Server:     peer.Addr,
					ServerPort: peer.Port,
				},
				Password: peer.UUID,
				OutboundTLSOptionsContainer: option.OutboundTLSOptionsContainer{
					TLS: &option.OutboundTLSOptions{
						Enabled:    true,
						ServerName: "gpp",
						Insecure:   true,
						ALPN:       option.Listable[string]{"h3"},
					},
				},
				BrutalDebug: false,
			},
		}
	case "direct":
		out = option.Outbound{
			Type: "direct",
		}
	default:
		out = option.Outbound{
			Type: "vless",
			VLESSOptions: option.VLESSOutboundOptions{
				ServerOptions: option.ServerOptions{
					Server:     peer.Addr,
					ServerPort: peer.Port,
				},
				UUID: peer.UUID,
				Transport: &option.V2RayTransportOptions{
					Type: "ws",
					WebsocketOptions: option.V2RayWebsocketOptions{
						Path:                fmt.Sprintf("/%s", peer.UUID),
						MaxEarlyData:        2048,
						EarlyDataHeaderName: "Sec-WebSocket-Protocol",
					},
				},
			},
		}
	}
	out.Tag = uuid.New().String()
	return out
}
func Client(gamePeer, httpPeer *config.Peer, processList []string) (*box.Box, error) {
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
						Tag:      "proxyDns",
						Address:  "8.8.8.8",
						Strategy: option.DomainStrategy(dns.DomainStrategyUseIPv4),
					},
					{
						Tag:      "localDns",
						Address:  "223.5.5.5",
						Detour:   "direct",
						Strategy: option.DomainStrategy(dns.DomainStrategyUseIPv4),
					},
					{
						Tag:      "block",
						Address:  "rcode://success",
						Strategy: option.DomainStrategy(dns.DomainStrategyUseIPv4),
					},
				},
				Rules: []option.DNSRule{
					{
						Type: "default",
						DefaultOptions: option.DefaultDNSRule{
							Server: "localDns",
							Domain: []string{
								"ghproxy.com",
								"cdn.jsdelivr.net",
								"testingcf.jsdelivr.net",
							},
						},
					},
					{
						Type: "default",
						DefaultOptions: option.DefaultDNSRule{
							Server: "localDns",
							Geosite: []string{
								"cn",
							},
						},
					},
					{
						Type: "default",
						DefaultOptions: option.DefaultDNSRule{
							Server: "proxyDns",
							Geosite: []string{
								"geolocation-!cn",
							},
						},
					},
				},
				DNSClientOptions: option.DNSClientOptions{
					DisableCache: true,
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
							netip.MustParsePrefix("172.225.0.1/30"),
						},
						AutoRoute:              true,
						StrictRoute:            true,
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
					DownloadURL:    "https://mirror.ghproxy.com/https://github.com/SagerNet/sing-geoip/releases/latest/download/geoip.db",
					DownloadDetour: "direct",
				},
				Geosite: &option.GeositeOptions{
					Path:           fmt.Sprintf("%s%c%s%c%s", home, os.PathSeparator, ".gpp", os.PathSeparator, "geosite.db"),
					DownloadURL:    "https://mirror.ghproxy.com/https://github.com/SagerNet/sing-geosite/releases/latest/download/geosite.db",
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
				},
			},
			Outbounds: []option.Outbound{
				getOUt(gamePeer),
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
	if processList != nil && len(processList) > 0 {
		options.Options.Route.Rules = append(options.Options.Route.Rules, option.Rule{
			Type: "logical",
			LogicalOptions: option.LogicalRule{
				Rules: []option.Rule{
					{
						Type: "default",
						DefaultOptions: option.DefaultRule{
							ProcessPath: processList,
						},
					},
				},
				Outbound: "direct",
				Invert:   true,
			},
		})
	} else {
		options.Options.Route.Rules = append(options.Options.Route.Rules, []option.Rule{
			{
				Type: "default",
				DefaultOptions: option.DefaultRule{
					Network:  option.Listable[string]{"udp"},
					Port:     []uint16{443},
					Outbound: "block",
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
		}...)
		if httpPeer != nil && httpPeer.Name != gamePeer.Name {
			out := getOUt(httpPeer)
			options.Options.Outbounds = append(options.Options.Outbounds, out)
			options.Options.Route.Rules = append(options.Options.Route.Rules, option.Rule{Type: "default", DefaultOptions: option.DefaultRule{Protocol: option.Listable[string]{"http", "quic", "tls"}, Outbound: out.Tag}})
		}
	}
	var instance, err = box.New(options)
	if err != nil {
		return nil, err
	}
	return instance, nil
}
