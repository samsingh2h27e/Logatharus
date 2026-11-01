package main

import (
	"context"

	Logger "github.com/samsingh2h27/Logatharus/Logger/LoggerClass"
)

func main() {
	logger := Logger.NewLogger("LogatharusService")
	context := context.Background()
	logger.Info(context,"Logger initialized successfully")
}
