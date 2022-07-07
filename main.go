package main

import (
	"jwt/src/routes"
	"os"
)

func main() {
	r := routes.SetupRouter()

	r.Run(":" + os.Getenv("APP_PORT"))
}
