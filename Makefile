server:
	set APP_ENV=dev && nodemon --watch "./**/*.go" --signal SIGTERM --exec "go run application.go"