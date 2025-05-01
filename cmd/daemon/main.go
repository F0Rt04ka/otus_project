package main

import (
	"time"

	"github.com/F0Rt04ka/otus_project/config"
	"github.com/F0Rt04ka/otus_project/internal/daemon"
)

func main() {
	config.Load()
	daemon.Run()

	time.Sleep(10000 * time.Second)

}
