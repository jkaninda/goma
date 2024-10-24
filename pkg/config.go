package pkg

import (
	"fmt"
	"github.com/jkaninda/goma-gateway/util"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Middleware struct {
	Path        string `yaml:"path"`
	AuthRequest struct {
		URL     string            `yaml:"url"`
		Headers map[string]string `yaml:"headers"`
		Params  map[string]string `yaml:"params"`
	} `yaml:"authRequest"`
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
	// Target Defines route blacklist
	Target      string `yaml:"target"`
	HealthCheck string `yaml:"healthCheck"`
	// Blocklist Defines route blacklist
	Blocklist []string `yaml:"blocklist"`

	// Middlewares Defines route middleware
	Middlewares []Middleware `yaml:"middlewares"`
}
type Gateway struct {
	ListenAddr   string `yaml:"listenAddr"`
	WriteTimeout int    `yaml:"writeTimeout"`
	ReadTimeout  int    `yaml:"readTimeout"`
	IdleTimeout  int    `yaml:"idleTimeout"`
	// RateLimiter Defines routes rateLimiter
	RateLimiter int               `yaml:"rateLimiter"`
	Headers     map[string]string `yaml:"headers"`
	Routes      []Route           `yaml:"routes"`
}
type GatewayConfig struct {
	GatewayConfig Gateway `yaml:"gateway"`
}

// ErrorResponse represents the structure of the JSON error response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	Data interface{} `json:"data"`
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
	return nil, fmt.Errorf("configuration file not found: %v", configFile)
}
func InitConfig() {
	initConfig()
	return

}
func initConfig() {
	conf := &GatewayConfig{
		GatewayConfig: Gateway{
			ListenAddr:   "localhost:8080",
			WriteTimeout: 15,
			ReadTimeout:  15,
			IdleTimeout:  60,
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Headers": "*",
				"Access-Control-Allow-Methods": "*",
			},
			Routes: []Route{
				{
					Name:        "HealthCheck",
					Path:        "/healthy",
					Target:      "http://localhost:8080",
					Rewrite:     "/",
					HealthCheck: "",
					Middlewares: []Middleware{
						{
							Path: "/admin",
							//AuthRequest: "",
						},
					},
				},
			},
		},
	}
	yamlData, err := yaml.Marshal(&conf)
	if err != nil {
		util.Fatal("Error %v", err.Error())
	}
	err = os.WriteFile("./data/config.yaml", yamlData, 0644)
	if err != nil {
		util.Fatal("Unable to write data into the file")
	}
	log.Println("Configuration file has been initialized successfully")
}
