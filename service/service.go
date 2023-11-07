package service

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	ErrTimeoutExceeded = errors.New("timeout exceeded")
)

type UrlIn struct {
	IsSequential bool
	OutFileName  string
	ReqTimeout   int
}

type UrlOut struct {
	Sequence      int
	Url           string
	ErrorMsg      string
	ContentLength int64
	ProcessedIn   float64
}

func Handle(in UrlIn, file *os.File) error {
	switch in.IsSequential {
	case true:
		fileScanner := bufio.NewScanner(file)
		urls := []string{}
		for fileScanner.Scan() {
			url := strings.TrimSpace(fileScanner.Text())
			if url != "" {
				urls = append(urls, url)
			}
		}

		out := KeepSequence(urls, in.ReqTimeout)
		if in.OutFileName != "" {
			return writeToFile(out, in.OutFileName)
		} else {
			printData(out)
		}
	case false:
		FastHandle(file, in.ReqTimeout)
	}

	return nil
}

func KeepSequence(urls []string, reqTimeout int) []UrlOut {
	c := make(chan UrlOut)

	wg := &sync.WaitGroup{}
	wg.Add(len(urls))

	go func() {
		for i, val := range urls {
			go func(idx int, url string) {
				defer wg.Done()

				now := time.Now()
				lenght, err := get(url, reqTimeout)
				procT := time.Since(now).Seconds()
				ans := UrlOut{
					Sequence:      idx,
					Url:           url,
					ContentLength: lenght,
					ProcessedIn:   procT,
				}
				if err != nil {
					ans.ErrorMsg = err.Error()
				}
				c <- ans
			}(i, val)
		}
	}()

	go func() {
		wg.Wait()
		close(c)
	}()

	out := make([]UrlOut, len(urls))
	for val := range c {
		out[val.Sequence] = val
	}

	return out
}

func FastHandle(file *os.File, reqTimeout int) {
	wg := &sync.WaitGroup{}
	fileScanner := bufio.NewScanner(file)

	for fileScanner.Scan() {
		line := strings.TrimSpace(fileScanner.Text())
		if line != "" {
			wg.Add(1)
			go func(url string) {
				defer wg.Done()

				now := time.Now()
				lenght, err := get(url, reqTimeout)
				procT := time.Since(now).Seconds()
				out := UrlOut{
					Url:           url,
					ContentLength: lenght,
					ProcessedIn:   procT,
				}
				if err != nil {
					out.ErrorMsg = err.Error()
				}
				fmt.Print(createStrFromUrlOut(out, false))
			}(line)
		}
	}

	wg.Wait()
}

func get(url string, timeout int) (int64, error) {
	factor := 5
	if timeout != 0 {
		factor = timeout
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(time.Second*time.Duration(factor)))
	defer cancel()
	req = req.WithContext(ctx)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "context deadline exceeded") {
			return 0, ErrTimeoutExceeded
		}
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf(fmt.Sprintf("status code %d is not expected", resp.StatusCode))
	}

	for key, values := range resp.Header {
		if key == "Content-Length" {
			if len(values) > 0 {
				ln, err := strconv.ParseInt(values[0], 10, 64)
				if err != nil {
					return 0, err
				}

				return ln, nil
			}
		}
	}

	// if headers does not contain info about content length
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	return int64(len(body)), nil
}

func writeToFile(out []UrlOut, outFileName string) error {
	var sb strings.Builder
	for _, val := range out {
		str := createStrFromUrlOut(val, true)
		sb.WriteString(str)
	}

	f, err := os.Create(outFileName)
	if err != nil {
		return err
	}
	defer f.Close()
	f.Write([]byte(sb.String()))

	return nil
}

func printData(out []UrlOut) {
	for _, val := range out {
		fmt.Print(createStrFromUrlOut(val, true))
	}
}

func createStrFromUrlOut(out UrlOut, ordered bool) string {
	str := ""
	switch {
	case out.ErrorMsg != "":
		str = fmt.Sprintf("%v - status: failed, error: %v, processed in: %.3f sec\n",
			out.Url, out.ErrorMsg, out.ProcessedIn)
	default:
		str = fmt.Sprintf("%v - status: succeed, content length: %v, processed in: %.3f sec\n",
			out.Url, out.ContentLength, out.ProcessedIn)
	}

	if ordered {
		str = fmt.Sprintf("%d. %v", out.Sequence+1, str)
	}

	return str
}
