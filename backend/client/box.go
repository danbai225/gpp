package client

import (
	"client/backend/config"
	"context"
	"fmt"
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
			Tag:  fmt.Sprintf("ss-out-%s", peer.UUID),
			ShadowsocksOptions: option.ShadowsocksOutboundOptions{
				ServerOptions: option.ServerOptions{
					Server:     peer.Addr,
					ServerPort: peer.Port,
				},
				Method:   "xchacha20-ietf-poly1305",
				Password: peer.UUID,
				Multiplex: &option.OutboundMultiplexOptions{
					Enabled:        true,
					Protocol:       "h2mux",
					MaxConnections: 4,
					MinStreams:     4,
					MaxStreams:     0,
					Padding:        false,
				},
			},
		}
	case "socks":
		out = option.Outbound{
			Type: "socks",
			Tag:  fmt.Sprintf("socks-out-%s", peer.UUID),
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
	case "direct":
		out = option.Outbound{
			Type: "direct",
			Tag:  "direct-out",
		}
	default:
		out = option.Outbound{
			Type: "vless",
			Tag:  fmt.Sprintf("vless-out-%s", peer.UUID),
			VLESSOptions: option.VLESSOutboundOptions{
				ServerOptions: option.ServerOptions{
					Server:     peer.Addr,
					ServerPort: peer.Port,
				},
				UUID: peer.UUID,
				Multiplex: &option.OutboundMultiplexOptions{
					Enabled:        true,
					Protocol:       "h2mux",
					MaxConnections: 4,
					MinStreams:     4,
					MaxStreams:     0,
					Padding:        false,
				},
			},
		}
	}
	return out
}
func Client(gamePeer, httpPeer *config.Peer) (*box.Box, error) {
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
						Address:  "https://223.5.5.5/dns-query",
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
							Server: "block",
							Geosite: []string{
								"category-ads-all",
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
	if httpPeer != nil && httpPeer.Name != gamePeer.Name {
		out := getOUt(httpPeer)
		options.Options.Outbounds = append(options.Options.Outbounds, out)
		options.Options.Route.Rules = append(options.Options.Route.Rules, option.Rule{Type: "default", DefaultOptions: option.DefaultRule{Protocol: option.Listable[string]{"http", "quic", "tls"}, Outbound: out.Tag}})
	}
	var instance, err = box.New(options)
	if err != nil {
		return nil, err
	}
	return instance, nil
}
