package main

import (
	"context"

	Logger "github.com/samsingh2h27/Logatharus/Logger/LoggerClass"
)

func main() {
	logger,_:= Logger.NewLogger("LogatharusService","your_token_here")
	context := context.Background()
	logger.Info(context,"Logger initialized successfully")
	logger.Debug(context,"This is a log message from LogatharusService")
	logger.Error(context,"This is an error log message from LogatharusService")
	logger.Warn(context,"This is a warning log message from LogatharusService")

}
