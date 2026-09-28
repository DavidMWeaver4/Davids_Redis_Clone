package main

import (
	"time"

	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/persistence"
)

const (
	serverAddr   = ":6379"
	aofFilePath  = "appendonly.aof"
	fsyncPolicy  = persistence.FsyncEverySec
	shutdownTime = 10 * time.Second
)
