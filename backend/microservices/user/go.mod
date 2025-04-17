module user

go 1.23.4

require (
	db v0.0.0-00010101000000-000000000000
	error v0.0.0-00010101000000-000000000000
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/joho/godotenv v1.5.1
	github.com/rabbitmq/amqp091-go v1.10.0
	github.com/redis/go-redis/v9 v9.7.3
	github.com/rs/cors v1.11.1
	middlewares v0.0.0-00010101000000-000000000000
	sharedTypes v0.0.0-00010101000000-000000000000
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/lib/pq v1.10.9 // indirect
)

replace middlewares => ../../modules/middlewares

replace db => ../../modules/db

replace error => ../../modules/error

replace sharedTypes => ../../modules/sharedTypes/
