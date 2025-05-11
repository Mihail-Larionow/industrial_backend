package main

import (
	"github.com/Mihail-Larionow/industrial_backend/internal/app"

	_ "github.com/Mihail-Larionow/industrial_backend/docs"
)

// @title My API
// @version 1.0
// @description Test API
// @host localhost:8080
// @BasePath /
func main() {
	app.Run()
}
