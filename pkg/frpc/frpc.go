package frpc

import (
	"context"
	"github.com/fatedier/frp/client"
	v1 "github.com/fatedier/frp/pkg/config/v1"
	"github.com/fatedier/frp/pkg/util/log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	cfg := &v1.ClientCommonConfig{
		ServerAddr: "139.9.150.125",
		ServerPort: 8222,
		User:       "admin",
		Auth: v1.AuthClientConfig{
			Token: "jxkj123",
		},
	}

	proxyCfgs := []v1.ProxyConfigurer{
		&v1.TCPProxyConfig{
			ProxyBaseConfig: v1.ProxyBaseConfig{
				Name: "test",
				Type: "tcp",
				ProxyBackend: v1.ProxyBackend{
					LocalIP:   "127.0.0.1",
					LocalPort: 5080,
				},
			},
			RemotePort: 8787,
		},
	}

	visitorCfgs := []v1.VisitorConfigurer{}

	err := startService(cfg, proxyCfgs, visitorCfgs, "")
	if err != nil {
		return
	}
}

func startService(
	cfg *v1.ClientCommonConfig,
	proxyCfgs []v1.ProxyConfigurer,
	visitorCfgs []v1.VisitorConfigurer,
	cfgFile string,
) error {
	log.InitLogger(cfg.Log.To, cfg.Log.Level, int(cfg.Log.MaxDays), cfg.Log.DisablePrintColor)

	if cfgFile != "" {
		log.Infof("start frpc service for config file [%s]", cfgFile)
		defer log.Infof("frpc service for config file [%s] stopped", cfgFile)
	}
	svr, err := client.NewService(client.ServiceOptions{
		Common:         cfg,
		ProxyCfgs:      proxyCfgs,
		VisitorCfgs:    visitorCfgs,
		ConfigFilePath: cfgFile,
	})
	if err != nil {
		return err
	}

	shouldGracefulClose := cfg.Transport.Protocol == "kcp" || cfg.Transport.Protocol == "quic"
	// Capture the exit signal if we use kcp or quic.
	if shouldGracefulClose {
		go handleTermSignal(svr)
	}
	return svr.Run(context.Background())
}

func handleTermSignal(svr *client.Service) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	svr.GracefulClose(500 * time.Millisecond)
}
