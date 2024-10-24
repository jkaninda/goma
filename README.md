# Goma Gateway

Goma is a lightweight API Gateway, Reverse Proxy.

Example of configuration file
```yaml
##### Goma Gateway configurations
gateway:
  ########## Global settings
  listenAddr: 0.0.0.0:8080
  writeTimeout: 15
  readTimeout: 15
  idleTimeout: 60
  rateLimiter: 0
  cors:
    Access-Control-Allow-Origin: '*'
    Access-Control-Allow-Methods: 'GET, POST, PUT, DELETE, OPTIONS'
    Access-Control-Allow-Headers: 'Content-Type, Authorization'
##### Define routes
  routes:
    - name: Store
      path: /store
      rewrite: /
      destination: 'http://store-service:8080'
      # Internal healthCheck
      healthCheck: /internal/health/ready
      cors:
        Access-Control-Allow-Origin: '*'
        Access-Control-Allow-Methods: 'GET, POST, PUT, DELETE, OPTIONS'
        Access-Control-Allow-Headers: 'Content-Type, Authorization'
      #### Define route blocklist paths
      blocklist:
        - /swagger-ui/*
        - /v2/swagger-ui/*
        - /api-docs/*
        - /internal/*
        - /actuator/*
      ##### Define route middlewares
      middlewares:
        - path: /cart
          #Enables authorization based on the result of a subrequest and sets the URI to which the subrequest will be sent.
          http:
            url: http://security-service:8080/security/authUser
            # Required headers, if not present in the request, the proxy will block access
            requiredHeaders:
              - Authorization
            #Sets the request variable to the given value after the authorization request completes.
            #
            # Add header to the next request from AuthRequest header, depending on your requirements
            # Key is AuthRequest's response header Key, and value  is Request's header Key
            # In case you want to get headers from Authentication service and inject them to next request's headers
            headers:
              userId: X-Auth-UserId
              userCountryId: X-Auth-UserCountryId
            # In case you want to get headers from Authentication service and inject them to next request's params
            params:
              auth_userCountryId: countryId
        - path: /items
          #Enables basic authorization
          # Protect path with basic authentication
          basic:
            username: goma
            password: goma
        - path: /business
          http:
            url: http://security-service:8080/security/authUser
            headers:
              #Key from backend authentication is header, and inject to the request with custom key name
              userId: X-Auth-UserId
              userCountryId: X-Auth-UserCountryId
            params:
              userCountryId: X-countryId
    - name: Authentication service
      path: /auth
      rewrite: /
      destination: 'http://security-service:8080'
      healthCheck: /internal/health/ready
      cors: {}
      rateLimiter: 0
      blocklist: []
      middlewares: []
    - name: Business
      path: /cart
      rewrite: /
      destination: 'http://cart-service:8080'
      healthCheck: 
      cors: {}
      rateLimiter: 0
      blocklist: []
      middlewares: []
    
    - name: Notification
      path: /notification
      rewrite: /
      destination: 'http://notification-service:8080'
      healthCheck: 
      cors: {}
      rateLimiter: 0
      blocklist: []
      middlewares: []
```