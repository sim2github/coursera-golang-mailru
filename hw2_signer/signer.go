package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

func SingleHash(in, out chan interface{}) {
	start := time.Now()
	var (
		wg    = &sync.WaitGroup{}
		quota = make(chan struct{}, 1)
	)

	for v := range in {
		defer fmt.Printf("SingleHash = %+v\n", time.Since(start))
		var data = fmt.Sprintf("%v", v)

		wg.Add(1)
		go func() {
			c := make(chan string)
			m := make(chan string)

			go func() {
				quota <- struct{}{}
				m <- DataSignerMd5(data)
				<-quota
			}()
			go func() {
				c <- DataSignerCrc32(<-m)
			}()
			out <- DataSignerCrc32(data) + "~" + <-c
			wg.Done()
		}()
	}
	wg.Wait()
}

func MultiHash(in, out chan interface{}) {
	start := time.Now()
	wg := &sync.WaitGroup{}

	for i := range in {
		wg.Add(1)
		go func(data string, out chan interface{}, wg *sync.WaitGroup) {
			const N int = 6
			var (
				mu    = &sync.Mutex{}
				jobWg = &sync.WaitGroup{}
				res   = make([]string, N)
			)

			for i := 0; i < N; i++ {
				jobWg.Add(1)
				data := strconv.Itoa(i) + data

				go func(res []string, i int, data string, jobWg *sync.WaitGroup, mu *sync.Mutex) {
					defer fmt.Printf("MultiHash = %+v\n", time.Since(start))
					defer jobWg.Done()
					data = DataSignerCrc32(data)

					mu.Lock()
					res[i] = data
					mu.Unlock()
				}(res, i, data, jobWg, mu)
			}

			jobWg.Wait()
			out <- strings.Join(res, "")
			wg.Done()
		}(i.(string), out, wg)

	}
	wg.Wait()
}

func CombineResults(in, out chan interface{}) {
	start := time.Now()
	defer fmt.Printf("CombineResults = %+v\n", time.Since(start))

	var res []string
	for v := range in {
		res = append(res, v.(string))
	}

	sort.Strings(res)
	out <- strings.Join(res, "_")
}

func ExecutePipeline(fn ...job) {
	start := time.Now()
	wg := &sync.WaitGroup{}
	in := make(chan interface{})
	out := in

	for _, f := range fn {
		out = make(chan interface{})
		wg.Add(1)
		go func(fn job, in, out chan interface{}, wg *sync.WaitGroup) {
			defer fmt.Printf("ExecutePipeline = %+v\n", time.Since(start))
			defer wg.Done()
			defer close(out)
			fn(in, out)
		}(f, in, out, wg)
		in = out
	}
	wg.Wait()
}

func main() {
	inputData := []int{0, 1, 1, 2, 3, 5, 8}
	// inputData := []int{0, 1}
	testResult := "NOT_SET"

	hashSignJobs := []job{
		job(func(in, out chan interface{}) {
			for _, fibNum := range inputData {
				out <- fibNum
			}
		}),
		job(SingleHash),
		job(MultiHash),
		job(CombineResults),
		job(func(in, out chan interface{}) {
			dataRaw := <-in
			data, ok := dataRaw.(string)
			if !ok {
				panic("cant convert result data to string")
			}
			testResult = data
		}),
	}

	start := time.Now()
	ExecutePipeline(hashSignJobs...)
	fmt.Printf("Result = %+v\n", testResult)
	fmt.Printf("time.Since(start) = %+v\n", time.Since(start))
}
