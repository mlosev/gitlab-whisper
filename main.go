package main

import (
	"github.com/mlosev/gitlab-whisper/cmd"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	cmd.Execute()
}
