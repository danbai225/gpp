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
				Method:   "aes-256-gcm",
				Password: peer.UUID,
				UDPOverTCP: &option.UDPOverTCPOptions{
					Enabled: true,
					Version: 2,
				},
				Multiplex: &option.OutboundMultiplexOptions{
					Enabled:        true,
					Protocol:       "h2mux",
					MaxConnections: 16,
					MinStreams:     32,
					Padding:        false,
				},
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
				Multiplex: &option.OutboundMultiplexOptions{
					Enabled:        true,
					Protocol:       "h2mux",
					MaxConnections: 16,
					MinStreams:     32,
					Padding:        false,
				},
			},
		}
	}
	out.Tag = uuid.New().String()
	return out
}
func Client(gamePeer, httpPeer *config.Peer, proxyDNS, localDNS string, rules []option.Rule) (*box.Box, error) {
	home, _ := os.UserHomeDir()
	proxyOut := getOUt(gamePeer)
	proxyOut.Tag = "proxy"
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
						Address:  proxyDNS,
						Detour:   "proxy",
						Strategy: option.DomainStrategy(dns.DomainStrategyUseIPv4),
					},
					{
						Tag:      "localDns",
						Address:  localDNS,
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
						Address: option.Listable[netip.Prefix]{
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
					DownloadURL:    "https://ghp.ci/https://github.com/SagerNet/sing-geoip/releases/latest/download/geoip.db",
					DownloadDetour: "direct",
				},
				Geosite: &option.GeositeOptions{
					Path:           fmt.Sprintf("%s%c%s%c%s", home, os.PathSeparator, ".gpp", os.PathSeparator, "geosite.db"),
					DownloadURL:    "https://ghp.ci/https://github.com/SagerNet/sing-geosite/releases/latest/download/geosite.db",
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
				proxyOut,
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
		{
			Type: "default",
			DefaultOptions: option.DefaultRule{
				DomainSuffix: option.Listable[string]{"vivox.com",
					"cm.steampowered.com",
					"steamchina.com",
					"steamcontent.com",
					"steamserver.net",
					"steamusercontent.com",
				},
				Outbound: "direct",
			},
		}, {
			Type: "default",
			DefaultOptions: option.DefaultRule{
				SourceIPCIDR: option.Listable[string]{"63.251.140.0/24",
					"69.25.124.0/23",
					"70.42.8.0/24",
					"70.42.198.0/23",
					"74.201.102.0/23",
					"74.201.106.0/23",
					"74.201.105.108/30",
					"85.236.96.0/21",
					"188.42.95.0/24",
					"188.42.147.0/24",
					"216.52.53.0/24"},
				Outbound: "direct",
			},
		}, {
			Type: "default",
			DefaultOptions: option.DefaultRule{
				Domain: option.Listable[string]{"csgo.wmsj.cn",
					"dl.steam.clngaa.com",
					"dl.steam.ksyna.com",
					"dota2.wmsj.cn",
					"st.dl.bscstorage.net",
					"st.dl.eccdnx.com",
					"st.dl.pinyuncloud.com",
					"steampipe.steamcontent.tnkjmec.com",
					"steampowered.com.8686c.com",
					"steamstatic.com.8686c.com",
					"wmsjsteam.com",
					"xz.pphimalayanrt.com"},
				Outbound: "direct",
			},
		},
	}...)
	options.Options.Route.Rules = append(options.Options.Route.Rules, rules...)
	// http
	if httpPeer != nil && httpPeer.Name != gamePeer.Name {
		out := getOUt(httpPeer)
		options.Options.Outbounds = append(options.Options.Outbounds, out)
		options.Options.Route.Rules = append(options.Options.Route.Rules, option.Rule{Type: "default", DefaultOptions: option.DefaultRule{Protocol: option.Listable[string]{"http"}, Outbound: out.Tag}})
		options.Options.Route.Rules = append(options.Options.Route.Rules, option.Rule{Type: "default", DefaultOptions: option.DefaultRule{Network: option.Listable[string]{"tcp"}, Port: []uint16{80, 443, 8080, 8443}, Outbound: out.Tag}})
	}
	var instance, err = box.New(options)
	if err != nil {
		return nil, err
	}
	return instance, nil
}
