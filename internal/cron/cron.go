package cron

import (
	"runtime"

	"github.com/ICEY1W32/proxypool/config"
	"github.com/ICEY1W32/proxypool/internal/app"
	"github.com/ICEY1W32/proxypool/internal/cache"
	"github.com/ICEY1W32/proxypool/log"
	"github.com/ICEY1W32/proxypool/pkg/geoIp"
	"github.com/ICEY1W32/proxypool/pkg/healthcheck"
	"github.com/ICEY1W32/proxypool/pkg/provider"
	"github.com/jasonlvhit/gocron"
)

func Cron() {
	_ = gocron.Every(config.Config.CrawlInterval).Minutes().Do(crawlTask)
	_ = gocron.Every(config.Config.SpeedTestInterval).Minutes().Do(speedTestTask)
	_ = gocron.Every(config.Config.ActiveInterval).Minutes().Do(frequentSpeedTestTask)
	_ = gocron.Every(1).Day().At("04:30").Do(geoIp.UpdateGeoIP)
	<-gocron.Start()
}

func crawlTask() {
	err := app.InitConfigAndGetters("")
	if err != nil {
		log.Errorln("[cron.go] config parse error: %s", err)
	}
	app.CrawlGo()
	app.Getters = nil
	runtime.GC()
}

func speedTestTask() {
	log.Infoln("Doing speed test task...")
	err := config.Parse("")
	if err != nil {
		log.Errorln("[cron.go] config parse error: %s", err)
	}
	pl := cache.GetProxies("proxies")

	app.SpeedTest(pl)
	cache.SetString("clashproxies", provider.Clash{
		Base: provider.Base{
			Proxies: &pl,
		},
	}.Provide()) // update static string provider
	cache.SetString("surgeproxies", provider.Surge{
		Base: provider.Base{
			Proxies: &pl,
		},
	}.Provide())
	cache.SetString("loonproxies", provider.Loon{
		Base: provider.Base{
			Proxies: &pl,
		},
	}.Provide())
	runtime.GC()
}

func frequentSpeedTestTask() {
	log.Infoln("Doing speed test task for active proxies...")
	err := config.Parse("")
	if err != nil {
		log.Errorln("[cron.go] config parse error: %s", err)
	}
	pl_all := cache.GetProxies("proxies")
	pl := healthcheck.ProxyStats.ReqCountThan(config.Config.ActiveFrequency, pl_all, true)
	if len(pl) > int(config.Config.ActiveMaxNumber) {
		pl = healthcheck.ProxyStats.SortProxiesBySpeed(pl)[:config.Config.ActiveMaxNumber]
	}
	log.Infoln("Active proxies count: %d", len(pl))

	app.SpeedTest(pl)
	cache.SetString("clashproxies", provider.Clash{
		Base: provider.Base{
			Proxies: &pl_all,
		},
	}.Provide()) // update static string provider
	cache.SetString("surgeproxies", provider.Surge{
		Base: provider.Base{
			Proxies: &pl_all,
		},
	}.Provide())
	cache.SetString("loonproxies", provider.Loon{
		Base: provider.Base{
			Proxies: &pl_all,
		},
	}.Provide())
	runtime.GC()
}
