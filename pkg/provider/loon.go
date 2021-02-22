package provider

import (
	"strings"

	"github.com/Sansui233/proxypool/pkg/tool"

	"github.com/Sansui233/proxypool/pkg/proxy"
)

// Loon provides functions that make proxies support clash client
type Loon struct {
	Base
}

// Provide of Surge generates proxy list supported by surge
func (s Loon) Provide() string {
	s.preFilter()

	var resultBuilder strings.Builder
	for _, p := range *s.Proxies {
		if checkLoonSupport(p) {
			resultBuilder.WriteString(p.ToLoon() + "\n")
		}
	}
	return resultBuilder.String()
}

func checkLoonSupport(p proxy.Proxy) bool {
	switch p.(type) {
	case *proxy.ShadowsocksR:
		ssr := p.(*proxy.ShadowsocksR)
		if tool.CheckInList(proxy.SSRCipherList, ssr.Cipher) && tool.CheckInList(ssrProtocolList, ssr.Protocol) && tool.CheckInList(ssrObfsList, ssr.Obfs) {
			return true
		}
	case *proxy.Vmess:
		return true
	case *proxy.Shadowsocks:
		ss := p.(*proxy.Shadowsocks)
		if tool.CheckInList(proxy.SSCipherList, ss.Cipher) {
			return true
		}
	case *proxy.Trojan:
		return true
	default:
		return false
	}
	return false
}
