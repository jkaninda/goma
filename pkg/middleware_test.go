package pkg

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"testing"
)

func TestMiddleware(t *testing.T) {
	middlewares := []Middleware{
		{
			Name: "basic-auth",
			Type: "basic",
			Rule: BasicRule{
				Username: "goma",
				Password: "goma",
			},
		}, {
			Name: "google-jwt",
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
	buf, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("Unable to read config file %s", configFile)
	}
	c := &[]Middleware{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		t.Fatalf("Unable to parse config file %s", configFile)
	}
	for _, m := range *c {
		switch m.Type {
		case "basic":
			log.Println("Basic auth")
			basicAuth, err := ToBasicAuth(m.Rule)
			if err != nil {
				log.Fatalln("error:", err)
			}
			log.Printf("Username: %s and password: %s\n", basicAuth.Username, basicAuth.Password)
		case "jwt":
			log.Println("JWT auth")
			jwt, err := ToJWTRuler(m.Rule)
			if err != nil {
				log.Fatalln("error:", err)
			}
			log.Printf("JWT authentification URL is %s\n", jwt.URL)
		default:
			t.Errorf("Unknown middleware type %s", m.Type)

		}
	}

}
