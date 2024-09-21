package process

import (
	"runtime"
	"runtime/debug"
)

func SetGCOptions() {
	// Optimizing GC
	debug.SetGCPercent(1000)

	// Limiting os threads to available cpu
	runtime.GOMAXPROCS(runtime.NumCPU())
}
