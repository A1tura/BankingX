module user

go 1.23.4

require (
	db v0.0.0-00010101000000-000000000000
	error v0.0.0-00010101000000-000000000000
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/joho/godotenv v1.5.1
	github.com/rabbitmq/amqp091-go v1.10.0
	middlewares v0.0.0-00010101000000-000000000000
)

require github.com/lib/pq v1.10.9 // indirect

replace middlewares => ../../modules/middlewares
replace db => ../../modules/db
replace error => ../../modules/error
