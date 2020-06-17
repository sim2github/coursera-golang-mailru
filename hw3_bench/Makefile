.PHONY: all bench pprof

all: bench pprof

bench:
	go test -bench . -benchmem -cpuprofile=cpu.out -memprofile=mem.out -memprofilerate=1 .

pprof:
	go tool pprof hw3_bench.test cpu.out

json:
	easyjson -all user/user.go
