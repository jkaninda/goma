package pkg

import (
	"fmt"
	"github.com/jkaninda/goma/internal/logger"
	"github.com/jkaninda/goma/util"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"os"
)

var cfg *Gateway

type Config struct {
	file string
}
type BasicMiddle struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// HttpMiddle authentication using HTTP GET method
//
// HttpMiddle contains the authentication details
type HttpMiddle struct {
	// URL contains the authentication URL, it supports HTTP GET method only.
	URL string `yaml:"url"`
	// RequiredHeaders , contains required before sending request to the backend.
	RequiredHeaders []string `yaml:"requiredHeaders"`
	// Headers Add header to the backend from Authentication request's header, depending on your requirements.
	// Key is Http's response header Key, and value  is the backend Request's header Key.
	// In case you want to get headers from Authentication service and inject them to backend request's headers.
	Headers map[string]string `yaml:"headers"`
	// Params same as Headers, contains the request params.
	//
	// Gets authentication headers from authentication request and inject them as request params to the backend.
	//
	// Key is Http's response header Key, and value  is the backend Request's request param Key.
	//
	// In case you want to get headers from Authentication service and inject them to next request's params.
	//
	//e.g: Header X-Auth-UserId to query userId
	Params map[string]string `yaml:"params"`
}

// Middleware defined the route middleware
type Middleware struct {
	//Path contains the protected route path
	Path string `yaml:"path"`
	// Http authentication using HTTP GET method
	//
	// Http contains the authentication details
	Http HttpMiddle `yaml:"http"`
	// Basic contains basic-auth authentication details
	//
	// Protects a route path with Basic-Auth
	Basic BasicMiddle `yaml:"basic"`
}

// Route defines gateway route
type Route struct {
	// Name defines route name
	Name string `yaml:"name"`
	// Path defines route path
	Path string `yaml:"path"`
	// Rewrite rewrites route path to desired path
	//
	// E.g. /cart to / => It will rewrite /cart path to /
	Rewrite string `yaml:"rewrite"`
	// Destination Defines backend URL
	Destination string `yaml:"destination"`
	// Cors contains the route cors headers
	Cors map[string]string `yaml:"cors"`
	// DisableHeaderXForward Disable X-forwarded header.
	//
	// [X-Forwarded-Host, X-Forwarded-For, Host, Scheme ]
	//
	// It will not match the backend route
	DisableHeaderXForward bool `yaml:"disableHeaderXForward"`
	// HealthCheck Defines the backend is health check PATH
	HealthCheck string `yaml:"healthCheck"`
	// Blocklist Defines route blacklist
	Blocklist []string `yaml:"blocklist"`
	// Middlewares Defines route middleware
	Middlewares []Middleware `yaml:"middlewares"`
}

// Gateway contains Goma Proxy Gateway's configs
type Gateway struct {
	// ListenAddr Defines the server listenAddr
	//
	//e.g: localhost:8080
	ListenAddr string `yaml:"listenAddr" env:"GOMA_LISTEN_ADDR, overwrite"`
	// WriteTimeout defines proxy write timeout
	WriteTimeout int `yaml:"writeTimeout" env:"GOMA_WRITE_TIMEOUT, overwrite"`
	// ReadTimeout defines proxy read timeout
	ReadTimeout int `yaml:"readTimeout" env:"GOMA_READ_TIMEOUT, overwrite"`
	// IdleTimeout defines proxy idle timeout
	IdleTimeout int `yaml:"idleTimeout" env:"GOMA_IDLE_TIMEOUT, overwrite"`
	// RateLimiter Defines routes rateLimiter
	RateLimiter                  int    `yaml:"rateLimiter" env:"GOMA_RATE_LIMITER, overwrite"`
	AccessLog                    string `yaml:"accessLog" env:"GOMA_ACCESS_LOG, overwrite"`
	ErrorLog                     string `yaml:"errorLog" env:"GOMA_ERROR_LOG=, overwrite"`
	DisableRouteHealthCheckError bool   `yaml:"disableRouteHealthCheckError"`
	//Disable dispelling routes on start
	DisableDisplayRouteOnStart bool `yaml:"disableDisplayRouteOnStart"`
	// Cors contains the proxy headers
	//
	//e.g:
	//
	//Access-Control-Allow-Origin: '*'
	//
	//    Access-Control-Allow-Methods: 'GET, POST, PUT, DELETE, OPTIONS'
	//
	//    Access-Control-Allow-Cors: 'Content-Type, Authorization'
	Cors map[string]string `yaml:"cors"`
	// Routes defines the proxy routes
	Routes []Route `yaml:"routes"`
}
type GatewayConfig struct {
	GatewayConfig Gateway `yaml:"gateway"`
}

// ErrorResponse represents the structure of the JSON error response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}
type GatewayServer struct {
	gateway         Gateway
	HealthyCallback func()
	stop            chan struct{} // channel for waiting shutdown
}

// New reads config file and returns Gateway
func (GatewayServer) New(configFile string) (*GatewayServer, error) {
	if util.FileExists(configFile) {
		buf, err := os.ReadFile(configFile)
		if err != nil {
			return nil, err
		}
		util.SetEnv("GOMA_CONFIG_FILE", configFile)
		c := &GatewayConfig{}
		err = yaml.Unmarshal(buf, c)
		if err != nil {
			return nil, fmt.Errorf("in file %q: %w", configFile, err)
		}
		return &GatewayServer{
			stop:    make(chan struct{}),
			gateway: c.GatewayConfig,
			HealthyCallback: func() {
				logger.Info("Healthcheck call back")
			},
		}, nil
	}
	logger.Error("configuration file not found: %v", configFile)
	logger.Info("Generating new configuration file...")
	initConfig(ConfigFile)
	util.SetEnv("GOMA_CONFIG_FILE", ConfigFile)
	buf, err := os.ReadFile(ConfigFile)
	if err != nil {
		return nil, err
	}
	c := &GatewayConfig{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %w", ConfigFile, err)
	}
	logger.Info("Generating new configuration file...done")
	logger.Info("Starting server with default configuration")
	return &GatewayServer{
		stop:    make(chan struct{}),
		gateway: c.GatewayConfig,
	}, nil
	//return nil, fmt.Errorf("configuration file not found: %v", configFile)
}
func GetConfigPaths() string {
	return util.GetStringEnv("GOMAY_CONFIG_FILE", ConfigFile)
}
func InitConfig(cmd *cobra.Command) {
	configFile, _ := cmd.Flags().GetString("output")
	if configFile == "" {
		configFile = GetConfigPaths()
	}
	initConfig(configFile)
	return

}
func initConfig(configFile string) {
	if configFile == "" {
		configFile = GetConfigPaths()
	}
	conf := &GatewayConfig{
		GatewayConfig: Gateway{
			ListenAddr:                   "0.0.0.0:80",
			WriteTimeout:                 15,
			ReadTimeout:                  15,
			IdleTimeout:                  60,
			AccessLog:                    "/dev/Stdout",
			ErrorLog:                     "/dev/stderr",
			DisableRouteHealthCheckError: false,
			DisableDisplayRouteOnStart:   false,
			Cors: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Cors":    "*",
				"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
			},
			Routes: []Route{
				{
					Name:        "HealthCheck",
					Path:        "/healthy",
					Destination: "http://localhost:8080",
					Rewrite:     "/health",
					HealthCheck: "",
					Cors: map[string]string{
						"Access-Control-Allow-Origin":  "*",
						"Access-Control-Allow-Cors":    "*",
						"Access-Control-Allow-Methods": "GET, OPTIONS",
					},
					Middlewares: []Middleware{
						{
							Path:  "/admin",
							Http:  HttpMiddle{},
							Basic: BasicMiddle{},
						},
					},
				},
				{
					Name:        "Basic auth",
					Path:        "/basic",
					Destination: "http://localhost:8080",
					Rewrite:     "/health",
					HealthCheck: "",
					Middlewares: []Middleware{
						{},
						{Path: "",
							Basic: BasicMiddle{
								Username: "goma",
								Password: "goma",
							},
						},
					},
					Blocklist: []string{},
				},
			},
		},
	}
	yamlData, err := yaml.Marshal(&conf)
	if err != nil {
		logger.Fatal("Error serializing configuration %v", err.Error())
	}
	err = os.WriteFile(configFile, yamlData, 0644)
	if err != nil {
		logger.Fatal("Unable to write config file %s", err)
	}
	logger.Info("Configuration file has been initialized successfully")
}
func Get() *Gateway {
	if cfg == nil {
		c := &Gateway{}
		c.Setup(GetConfigPaths())
		cfg = c
	}
	return cfg
}
func (Gateway) Setup(conf string) *Gateway {
	if util.FileExists(conf) {
		buf, err := os.ReadFile(conf)
		if err != nil {
			return &Gateway{}
		}
		util.SetEnv("GOMA_CONFIG_FILE", conf)
		c := &GatewayConfig{}
		err = yaml.Unmarshal(buf, c)
		if err != nil {
			logger.Fatal("Error loading configuration %v", err.Error())
		}
		return &c.GatewayConfig
	}
	return &Gateway{}

}
