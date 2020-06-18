.PHONY: bench cpu mem json

bench:
	go test -bench . -benchmem -cpuprofile=cpu.out -memprofile=mem.out -memprofilerate=1 .

fast:
	go test -bench Fast -benchmem -cpuprofile=cpu.out -memprofile=mem.out -memprofilerate=1 -run Fast

cpu:
	go tool pprof hw3_bench.test cpu.out

mem:
	go tool pprof hw3_bench.test mem.out

json:
	easyjson -all user/user.go
