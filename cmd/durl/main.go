package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mrsubudei/url-cli/service"
)

func main() {
	// parse flags
	fileName := flag.String("i", "", "a string")
	sequentially := flag.Bool("s", false, "a bool")
	outFile := flag.String("o", "", "a string")
	reqTimeout := flag.Int("t", 0, "an int")

	flag.Parse()

	if *fileName == "" || len(flag.Args()) != 0 {
		fmt.Println(`Usage: durl -i="{file name}" [OPTION]...

	-i		input file name (required argument)
	-s		keep the original order of urls
	-o 		output file name (available only with -s flag)
	-t		custom timeout for http requests in sec (should be greater than 0)

Examples: 
durl -i="data.txt"			read data from file data.txt and print out the answer
durl -i="data.txt" -s -o="output.txt"	save the answer to output.txt
durl -i="data.txt" -t=10		set request timeout to 10 sec`)
		return
	}

	if !*sequentially && *outFile != "" {
		fmt.Println("flag -o only available with flag -s")
		return
	}

	if *reqTimeout < 0 {
		fmt.Println("request timeout should be integer greater than 0")
		return
	}

	// open file
	readFile, err := os.Open(*fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer readFile.Close()

	// handle
	err = service.Handle(service.UrlIn{
		IsSequential: *sequentially,
		OutFileName:  *outFile,
		ReqTimeout:   *reqTimeout,
	}, readFile)
	if err != nil {
		log.Fatal(err)
	}
}
