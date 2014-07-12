package main

import (
	"fmt"
	"os"
	"runtime/pprof"

	"github.com/zond/godec"
)

func main() {
	var slice []string
	for i := 0; i < 1000; i++ {
		slice = append(slice, fmt.Sprintf("String nr %v", i))
	}
	f, err := os.Create("cpuprofile")
	if err != nil {
		panic(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	for i := 0; i < 10000; i++ {
		b, err := godec.Marshal(slice)
		if err != nil {
			panic(err)
		}
		var newSlice []string
		err = godec.Unmarshal(b, &newSlice)
		if err != nil {
			panic(err)
		}
	}
}
