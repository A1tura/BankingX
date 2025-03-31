module user

go 1.23.4

require middlewares v0.0.0-00010101000000-000000000000

require (
	db v0.0.0-00010101000000-000000000000 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/lib/pq v1.10.9 // indirect
)

replace middlewares => ../middlewares

replace db => ../db
