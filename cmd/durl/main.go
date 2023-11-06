package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mrsubudei/url-cli/service"
)

func main() {
	fileName := flag.String("f", "", "a string")
	sequentially := flag.Bool("s", false, "a bool")
	outFile := flag.String("o", "", "a string")
	reqTimeout := flag.Int("t", 0, "an int")

	flag.Parse()

	if *fileName == "" || len(flag.Args()) != 0 {
		fmt.Println(`Usage: durl -f="{file name}" [OPTION]...

	-f		input file name (required argument)
	-s		keep the original order of urls
	-o 		outout file name (available only with -s flag)
	-t		custom timeout for http requests

Examples: 
durl -f="data.txt"			read data from file data.txt and print out the answer
durl -f="data.txt" -s -o="output.txt"	save the answer to output.txt
durl -f="data.txt" -t=10		set request timeout to 10 sec`)
		return
	}

	if !*sequentially && *outFile != "" {
		fmt.Println("flag -o only available with flag -s")
		return
	}

	if *reqTimeout < 1 {
		fmt.Println("request timeout should be integer greater or equal to 1")
		return
	}

	readFile, err := os.Open(*fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)
	urls := []string{}
	for fileScanner.Scan() {
		url := strings.TrimSpace(fileScanner.Text())
		if url != "" {
			urls = append(urls, url)
		}
	}

	err = service.Handle(service.UrlIn{
		Urls:         urls,
		IsSequential: *sequentially,
		OutFileName:  *outFile,
		ReqTimeout:   *reqTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}
}
