package pkg

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"testing"
)

const MidName = "google-jwt"

var rules = []string{"fake", "jwt", "google-jwt"}

func TestMiddleware(t *testing.T) {
	TestInit(t)
	middlewares := []Middleware{
		{
			Name: "basic-auth",
			Type: "basic",
			Rule: BasicRule{
				Username: "goma",
				Password: "goma",
			},
		}, {
			Name: MidName,
			Type: "jwt",
			Rule: JWTRuler{
				URL:     "https://www.googleapis.com/auth/userinfo.email",
				Headers: map[string]string{},
				Params:  map[string]string{},
			},
		},
	}
	yamlData, err := yaml.Marshal(&middlewares)
	if err != nil {
		t.Fatalf("Error serializing configuration %v", err.Error())
	}
	err = os.WriteFile(configFile, yamlData, 0644)
	if err != nil {
		t.Fatalf("Unable to write config file %s", err)
	}
	log.Printf("Config file written to %s", configFile)
}

func TestReadMiddleware(t *testing.T) {
	TestMiddleware(t)
	middlewares := getMiddlewares(t)
	middleware, err := searchMiddleware(rules, middlewares)
	if err != nil {
		t.Fatalf("Error searching middleware %s", err.Error())
	}
	switch middleware.Type {
	case "basic":
		log.Println("Basic auth")
		basicAuth, err := ToBasicAuth(middleware.Rule)
		if err != nil {
			log.Fatalln("error:", err)
		}
		log.Printf("Username: %s and password: %s\n", basicAuth.Username, basicAuth.Password)
	case "jwt":
		log.Println("JWT auth")
		jwt, err := ToJWTRuler(middleware.Rule)
		if err != nil {
			log.Fatalln("error:", err)
		}
		log.Printf("JWT authentification URL is %s\n", jwt.URL)
	default:
		t.Errorf("Unknown middleware type %s", middleware.Type)

	}

}

func TestFoundMiddleware(t *testing.T) {
	middlewares := getMiddlewares(t)
	middleware, err := searchMiddleware(rules, middlewares)
	if err != nil {
		t.Errorf("Error getting middleware %v", err)
	}
	fmt.Println(middleware.Type)
}

func getMiddlewares(t *testing.T) []Middleware {
	buf, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("Unable to read config file %s", configFile)
	}
	c := &[]Middleware{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		t.Fatalf("Unable to parse config file %s", configFile)
	}
	return *c
}
