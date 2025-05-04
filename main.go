package main

import (
	"plusz-backend/api"
	"plusz-backend/batch"
	"plusz-backend/env"
)

func main() {
	env.Load()
	go batch.CheckForNewSchedules()
	api.Init()
}
