package httpgovern

import (
	"context"
	"encoding/json"
	"git.4321.sh/feige/flygo/component/server"
	"git.4321.sh/feige/flygo/core/kapp"
	"git.4321.sh/feige/flygo/core/kcfg"
	"git.4321.sh/feige/flygo/core/klog"
	jsoniter "github.com/json-iterator/go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"time"
)

var (
	DefaultServeMux = http.NewServeMux()
	routes          = []string{}
)

const PackageName = "server.httpgovern"

func init() {
	HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		_ = json.NewEncoder(resp).Encode(routes)
	})
	HandleFunc("/debug/pprof/", pprof.Index)
	HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	HandleFunc("/debug/pprof/profile", pprof.Profile)
	HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	HandleFunc("/debug/pprof/trace", pprof.Trace)

	HandleFunc("/config/json", func(w http.ResponseWriter, r *http.Request) {
		encoder := json.NewEncoder(w)
		if r.URL.Query().Get("pretty") == "true" {
			encoder.SetIndent("", "    ")
		}
		_ = encoder.Encode(kcfg.Traverse("."))
	})
	HandleFunc("/config/raw", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(kcfg.RawConfig())
	})
	HandleFunc("/env/info", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_ = jsoniter.NewEncoder(w).Encode(os.Environ())
	})
	HandleFunc("/build/info", func(w http.ResponseWriter, r *http.Request) {
		serverStats := map[string]string{
			"name":         kapp.Name(),
			"appMode":      kapp.AppMode(),
			"appVersion":   kapp.AppVersion(),
			"flygoVersion": kapp.FlygoVersion(),
			"buildUser":    kapp.BuildUser(),
			"buildHost":    kapp.BuildHost(),
			"buildTime":    kapp.BuildTime(),
			"startTime":    kapp.StartTime(),
			"hostName":     kapp.HostName(),
			"goVersion":    kapp.GoVersion(),
		}
		_ = jsoniter.NewEncoder(w).Encode(serverStats)
	})
}

type Component struct {
	name   string
	config *Config
	logger *klog.Component
	*http.Server
	listener net.Listener
}

func newComponent(name string, config *Config, logger *klog.Component) *Component {
	return &Component{
		name:   name,
		logger: logger,
		Server: &http.Server{
			Addr:    config.Address(),
			Handler: DefaultServeMux,
		},
		listener: nil,
		config:   config,
	}
}

func (c *Component) ConfigKey() string {
	return c.name
}

func (c *Component) PackageName() string {
	return PackageName
}

func (c *Component) Start() error {
	go func() {
		time.Sleep(10 * time.Second)
		var listener, err = net.Listen("tcp4", c.config.Address())
		if err != nil {
			klog.Panic("governor start error", klog.FieldErr(err))
		}
		c.listener = listener
		HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
			promhttp.Handler().ServeHTTP(w, r)
		})
		_ = c.Server.Serve(c.listener)
		/*if err == http.ErrServerClosed {
			return nil
		}
		return err*/
	}()
	return nil
}

func (c *Component) Stop() error {
	return c.Server.Close()
}

func (c *Component) GraceShutdown(ctx context.Context) error {
	return c.Server.Shutdown(ctx)
}

func (c *Component) Info() *server.ServiceInfo {
	info := server.ApplyOptions(
		server.WithScheme("http"),
		// server.WithAddress(c.listener.Addr().String()),
		server.WithKind(kapp.ServiceHttpGovern),
	)
	return &info
}

func HandleFunc(pattern string, handler http.HandlerFunc) {
	DefaultServeMux.HandleFunc(pattern, handler)
	routes = append(routes, pattern)
}
