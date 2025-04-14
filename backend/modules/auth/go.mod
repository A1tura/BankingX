module auth

go 1.23.4

toolchain go1.23.8

replace error => ../../modules/error/

require (
	error v0.0.0-00010101000000-000000000000
	middlewares v0.0.0-00010101000000-000000000000
)

require (
	db v0.0.0-00010101000000-000000000000 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/rabbitmq/amqp091-go v1.10.0 // indirect
)

replace middlewares => ../../modules/middlewares/

replace db => ../../modules/db/
