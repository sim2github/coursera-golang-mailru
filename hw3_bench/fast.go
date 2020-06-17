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
	SlowSearch(os.Stdout)
}

// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	const N = 1000
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var (
		scanner      = bufio.NewScanner(file)
		seenBrowsers = map[string]struct{}{}
		userPool     = sync.Pool{
			New: func() interface{} {
				return new(user.User)
			}}
		dataPool = sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			}}
	)

	fmt.Fprintln(out, "found users:")
	var i = 0
	for scanner.Scan() {
		user := userPool.Get().(*user.User)
		defer userPool.Put(user)

		err := user.UnmarshalJSON(scanner.Bytes())
		if err != nil {
			log.Fatal(err)
		}
		user.Email = strings.ReplaceAll(user.Email, "@", " [at] ")

		for _, browser := range user.Browsers {
			var (
				isAndroid = strings.Contains(browser, "Android")
				isMSIE    = strings.Contains(browser, "MSIE")
			)

			if isAndroid {
				user.IsAndroid = isAndroid
			}

			if isMSIE {
				user.IsMSIE = isMSIE
			}

			if isAndroid || isMSIE {
				if _, ok := seenBrowsers[browser]; !ok {
					seenBrowsers[browser] = struct{}{}
				}
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
