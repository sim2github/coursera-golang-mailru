package main

import (
	"bufio"
	"bytes"
	"coursera/hw3_bench/user"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

func main() {
	FastSearch(os.Stdout)
	// SlowSearch(os.Stdout)
}

// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var (
		scanner      = bufio.NewScanner(file)
		seenBrowsers = map[string]struct{}{}
		dataPool     = sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			}}
		user = &user.User{}
	)

	fmt.Fprintln(out, "found users:")
	var i = 0
	for scanner.Scan() {
		// TODO: in future email field may contain more than one email devided by separator.
		ub := bytes.Replace(scanner.Bytes(), []byte("@"), []byte(" [at] "), 1)
		user.IsAndroid = false
		user.IsMSIE = false
		// TODO: check data - unmarshal must overwrite user.Browsers and user.Name at list with default values.
		err := user.UnmarshalJSON(ub)
		if err != nil {
			log.Fatal(err)
		}

		for _, browser := range user.Browsers {
			if strings.Contains(browser, "Android") {
				user.IsAndroid = true
			} else if strings.Contains(browser, "MSIE") {
				user.IsMSIE = true
			} else {
				continue
			}

			if _, ok := seenBrowsers[browser]; !ok {
				seenBrowsers[browser] = struct{}{}
			}
		}

		if user.IsAndroid && user.IsMSIE {
			b := dataPool.Get().(*bytes.Buffer)
			b.Reset()
			b.WriteByte('[')
			b.WriteString(strconv.Itoa(i))
			b.WriteString("] ")
			b.WriteString(user.Name)
			b.WriteString(" <")
			b.WriteString(user.Email)
			b.WriteString(">\n")
			out.Write(b.Bytes())
			dataPool.Put(b)
		}
		i++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Fprintln(out, "\nTotal unique browsers", len(seenBrowsers))
}
