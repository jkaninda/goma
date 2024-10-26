# Goma - simple lightweight API Gateway and Reverse Proxy.

```
   _____                       
  / ____|                      
 | |  __  ___  _ __ ___   __ _ 
 | | |_ |/ _ \| '_ ` _ \ / _` |
 | |__| | (_) | | | | | | (_| |
  \_____|\___/|_| |_| |_|\__,_|
                               
```
Goma is a lightweight API Gateway and Reverse Proxy.

[![Build](https://github.com/jkaninda/goma/actions/workflows/release.yml/badge.svg)](https://github.com/jkaninda/goma/actions/workflows/release.yml)
[![Go Report](https://goreportcard.com/badge/github.com/jkaninda/mysql-bkup)](https://goreportcard.com/report/github.com/jkaninda/goma)
[![Go Reference](https://pkg.go.dev/badge/github.com/jkaninda/goma.svg)](https://pkg.go.dev/github.com/jkaninda/goma)
![Docker Image Size (latest by date)](https://img.shields.io/docker/image-size/jkaninda/goma?style=flat-square)

## Links:

- [Docker Hub](https://hub.docker.com/r/jkaninda/goma)
- [Github](https://github.com/jkaninda/goma)

### Feature

- [x] Reverse proxy
- [x] API Gateway
- [x] Cors
- [ ] Add Load balancing feature
- [ ] Support TLS
- [x] Authentication middleware
  - [x] JWT `HTTP Bearer Token`
  - [x] Basic-Auth
  - [ ] OAuth2
- [x] Implement rate limiting
  - [x] In-Memory Token Bucket based
  - [x] In-Memory client IP based
  - [ ] Distributed Rate Limiting for Token based across multiple instances using Redis
  - [ ] Distributed Rate Limiting for In-Memory client IP based across multiple instances  using Redis

## Usage

### 1. Initialize configuration

```shell
docker run --rm  --name goma \
 -v "${PWD}/config:/config" \
 jkaninda/goma config init --output /config/goma.yml
```
### 2. Run server

```shell
docker run --rm --name goma \
 -v "${PWD}/config:/config" \
 -p 80:80 \
 jkaninda/goma server
```

### 3. Start server with a custom config
```shell
docker run --rm --name goma \
 -v "${PWD}/config:/config" \
 -p 80:80 \
 jkaninda/goma server --config /config/config.yml
```
### 4. Healthcheck

[http://localhost/health](http://localhost/health)

> Healthcheck response body

```json
{
	"status": "healthy",
	"routes": [
		{
			"name": "Store",
			"status": "healthy",
			"error": ""
		},
		{
			"name": "Authentication service",
			"status": "unhealthy",
          "error": "error performing HealthCheck request: Get \"http://notification-service:8080/internal/health/ready\": dial tcp: lookup notification-service on 127.0.0.11:53: no such host "
          
		},
		{
			"name": "Notification",
			"status": "undefined",
			"error": ""
		}
	]
}
```


Create a config file in this format
## Customize configuration file

Example of configuration file
```yaml
## Goma - simple lightweight API Gateway and Reverse Proxy.
# Goma Gateway configurations
gateway:
  ########## Global settings
  listenAddr: 0.0.0.0:80
  # Proxy write timeout
  writeTimeout: 15
  # Proxy read timeout
  readTimeout: 15
  # Proxy idle timeout
  idleTimeout: 60
  # Proxy rate limit, it's In-Memory Token Bucket
  # Distributed Rate Limiting for Token based across multiple instances is not yet integrated
  rateLimiter: 0
  accessLog:    "/dev/Stdout"
  errorLog:     "/dev/stderr"
  ## Returns backend route healthcheck errors
  disableRouteHealthCheckError: false
  # Disable display routes on start
  disableDisplayRouteOnStart: false
  # Proxy Global HTTP Cors
  cors:
    # Cors origins are global for all routes
    origins:
      - https://example.com
      - https://dev.example.com
      - http://localhost:80
    # Allowed headers are global for all routes
    headers:
      Access-Control-Allow-Headers: 'Origin, Authorization, Accept, Content-Type, Access-Control-Allow-Headers, X-Client-Id, X-Session-Id'
      Access-Control-Allow-Credentials: 'true'
      Access-Control-Max-Age: 1728000
  ##### Define routes
  routes:
    # Example of a route | 1
    - name: Store
      path: /store
      ## Rewrite a request path
      # e.g rewrite: /store to /
      rewrite: /
      destination: 'http://store-service:8080'
      #DisableHeaderXForward Disable X-forwarded header.
      # [X-Forwarded-Host, X-Forwarded-For, Host, Scheme ]
      # It will not match the backend route, by default, it's disabled
      disableHeaderXForward: false
      # Internal health check
      healthCheck: /internal/health/ready
      # Proxy route HTTP Cors
      cors:
        headers:
          Access-Control-Allow-Methods: 'GET'
          Access-Control-Allow-Headers: 'Origin, Authorization, Accept, Content-Type, Access-Control-Allow-Headers, X-Client-Id, X-Session-Id'
          Access-Control-Allow-Credentials: 'true'
          Access-Control-Max-Age: 1728000
      #### Define route blocklist paths
      blocklist:
        - /swagger-ui/*
        - /v2/swagger-ui/*
        - /api-docs/*
        - /internal/*
        - /actuator/*
      ##### Define route middlewares from middlewares names
      ## The name must be unique
      ## List of middleware name
      middlewares:
        # path to protect
        - path: /user/account
          # Rules defines which specific middleware applies to a route path
          rules:
            - auth
        # path to protect
        - path: /cart
          # Rules defines which specific middleware applies to a route path
          rules:
            - google-auth
            - auth
        - path: /history
          http:
            url: http://security-service:8080/security/authUser
            headers:
              #Key from backend authentication header, and inject to the request with custom key name
              userId: X-Auth-UserId
              userCountryId: X-Auth-UserCountryId
            params:
              userCountryId: X-countryId
    # Example of a route | 2
    - name: Authentication service
      path: /auth
      rewrite: /
      destination: 'http://security-service:8080'
      healthCheck: /internal/health/ready
      cors: {}
      blocklist: []
      middlewares: []
    # Example of a route | 3
    - name: Notification
      path: /notification
      rewrite: /
      destination: 'http://notification-service:8080'
      healthCheck:
      cors: {}
      blocklist: []
      middlewares: []

#Defines proxy middlewares
middlewares:
  # Enable Basic auth authorization based
  - name: local-auth-basic
    # Authentication types | jwt, basic, auth0
    type: basic
    rule:
      username: admin
      password: admin
  #Enables JWT authorization based on the result of a request and continues the request.
  - name: google-auth
    # Authentication types | jwt, basic, auth0
    type: jwt
    rule:
      url: https://www.googleapis.com/auth/userinfo.email
      # Required headers, if not present in the request, the proxy will return 403
      requiredHeaders:
        - Authorization
      #Sets the request variable to the given value after the authorization request completes.
      #
      # Add header to the next request from AuthRequest header, depending on your requirements
      # Key is AuthRequest's response header Key, and value is Request's header Key
      # In case you want to get headers from the Authentication service and inject them into the next request's headers
      #Sets the request variable to the given value after the authorization request completes.
      #
      # Add header to the next request from AuthRequest header, depending on your requirements
      # Key is AuthRequest's response header Key, and value is Request's header Key
      # In case you want to get headers from the Authentication service and inject them into the next request's headers
      headers:
        userId: X-Auth-UserId
        userCountryId: X-Auth-UserCountryId
      # In case you want to get headers from the Authentication service and inject them to the next request's params
      params:
        auth_userCountryId: countryId
```

## Requirement

- Docker
