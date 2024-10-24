package pkg

import (
	"fmt"
	"github.com/jkaninda/goma-gateway/internal/logger"
	"github.com/jkaninda/goma-gateway/util"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"os"
)

// Middleware defined the route middleware
type Middleware struct {
	//Path contains the protected route path
	Path string `yaml:"path"`
	// Http authentication using HTTP GET method
	//
	//Http contains the authentication details
	Http struct {
		// URL contains the authentication URL, it supports HTTP GET method only.
		URL string `yaml:"url"`
		// RequiredHeaders , contains required before sending request to the backend.
		RequiredHeaders []string `yaml:"requiredHeaders,omitempty"`
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
	} `yaml:"http"`
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
	// HealthCheck Defines the backend is health check
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
	ListenAddr string `yaml:"listenAddr"`
	// WriteTimeout defines proxy write timeout
	WriteTimeout int `yaml:"writeTimeout"`
	// ReadTimeout defines proxy read timeout
	ReadTimeout int `yaml:"readTimeout"`
	// IdleTimeout defines proxy idle timeout
	IdleTimeout int `yaml:"idleTimeout"`
	// RateLimiter Defines routes rateLimiter
	RateLimiter int `yaml:"rateLimiter"`
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

// config reads config file and returns Gateway
func loadConf(configFile string) (*Gateway, error) {
	if util.FileExists(configFile) {
		buf, err := os.ReadFile(configFile)
		if err != nil {
			return nil, err
		}

		c := &GatewayConfig{}
		err = yaml.Unmarshal(buf, c)
		if err != nil {
			return nil, fmt.Errorf("in file %q: %w", configFile, err)
		}
		return &c.GatewayConfig, err
	}
	logger.Error("configuration file not found: %v", configFile)
	logger.Info("Generating new configuration file...")
	initConfig(ConfigFile)
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
	return &c.GatewayConfig, err
	//return nil, fmt.Errorf("configuration file not found: %v", configFile)
}
func getConfigFile() string {
	return util.GetStringEnv("GOMA_PROXY_CONFIG_FILE", ConfigFile)
}
func InitConfig(cmd *cobra.Command) {
	configFile, _ := cmd.Flags().GetString("config")
	if configFile == "" {
		configFile = getConfigFile()
	}
	initConfig(configFile)
	return

}
func initConfig(configFile string) {
	if configFile == "" {
		configFile = getConfigFile()
	}
	conf := &GatewayConfig{
		GatewayConfig: Gateway{
			ListenAddr:   "0.0.0.0:8080",
			WriteTimeout: 15,
			ReadTimeout:  15,
			IdleTimeout:  60,
			Cors: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Cors":    "*",
				"Access-Control-Allow-Methods": "*",
			},
			Routes: []Route{
				{
					Name:        "HealthCheck",
					Path:        "/healthy",
					Destination: "http://localhost:8080",
					Rewrite:     "/health",
					HealthCheck: "",
					Cors:        map[string]string{},
					Middlewares: []Middleware{
						{
							Path: "/admin",
						},
					},
				},
				{
					Name:        "Hello",
					Path:        "/hello",
					Destination: "http://localhost:8080",
					Rewrite:     "/",
					HealthCheck: "",
					Middlewares: []Middleware{
						{},
						{Path: ""},
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
