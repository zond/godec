#!/bin/sh
./godec.test -test.run=asdfadsf -test.bench=BenchmarkGodecStringSlice -test.cpuprofile=cpuprofile.godec -test.memprofile=memprofile.godec
./godec.test -test.run=asdfadsf -test.bench=BenchmarkBincStringSlice -test.cpuprofile=cpuprofile.binc -test.memprofile=memprofile.binc
go tool pprof --svg --nodefraction=0.05 godec.test cpuprofile.binc > ~/tmp/cpuprofile.binc.svg
go tool pprof --svg --nodefraction=0.05 godec.test cpuprofile.godec > ~/tmp/cpuprofile.godec.svg
go tool pprof --svg --nodefraction=0.05 godec.test memprofile.binc > ~/tmp/memprofile.binc.svg
go tool pprof --svg --nodefraction=0.05 godec.test memprofile.godec > ~/tmp/memprofile.godec.svg
