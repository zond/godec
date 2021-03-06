#!/bin/sh

go test -c
./godec.test -test.run=asdfadsf -test.benchtime=10s -test.bench=$1 -test.cpuprofile=$1.cpuprof -test.memprofile=$1.memprof
go tool pprof --svg --ignore=runtime.mach_semaphore_signal --nodefraction=0.05 godec.test $1.cpuprof > $1.cpuprof.svg
go tool pprof --svg --ignore=runtime.mach_semaphore_signal --nodefraction=0.05 godec.test $1.memprof > $1.memprof.svg
