module user

go 1.23.4

require middlewares v0.0.0-00010101000000-000000000000

require (
	db v0.0.0-00010101000000-000000000000 // indirect
	github.com/lib/pq v1.10.9 // indirect
)

replace middlewares => ../middlewares

replace db => ../db
