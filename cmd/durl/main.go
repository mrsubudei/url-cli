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
	flag.Parse()
	if *fileName == "" {
		fmt.Println(`USAGE: durl -f="file name" OPTION[-s]
EXAMPLE: durl -f="data.txt OR durl -f="data.txt -s"`)
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

	service.Handle(urls, *sequentially)
}
