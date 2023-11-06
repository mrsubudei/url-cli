package service

import (
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
	RequestTimeOut = time.Second * 5

	ErrNoResponse = errors.New("no response from url")
)

type UrlOut struct {
	Sequence      int
	Url           string
	ErrorMsg      string
	ContentLength int64
	ProcessedIn   float64
}

func Handle(urls []string, sequentially bool, outFileName string) error {
	switch sequentially {
	case true:
		out := KeepSequence(urls)
		if outFileName != "" {
			return writeToFile(out, outFileName)
		} else {
			printData(out)
		}
	case false:
		FastHandle(urls)
	}

	return nil
}

func KeepSequence(urls []string) []UrlOut {
	c := make(chan UrlOut)

	wg := &sync.WaitGroup{}
	wg.Add(len(urls))

	go func() {
		for i, val := range urls {
			go func(idx int, url string) {
				defer wg.Done()

				now := time.Now()
				lenght, err := get(url)
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

func FastHandle(urls []string) {
	wg := &sync.WaitGroup{}
	wg.Add(len(urls))
	for _, val := range urls {
		go func(url string) {
			defer wg.Done()

			now := time.Now()
			lenght, err := get(url)
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
		}(val)
	}

	wg.Wait()
}

func get(url string) (int64, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(RequestTimeOut))
	defer cancel()
	req = req.WithContext(ctx)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "context deadline exceeded") {
			return 0, ErrNoResponse
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
		str = fmt.Sprintf("%v - status: failed, error: %v, processed in: %v\n",
			out.Url, out.ErrorMsg, out.ProcessedIn)
	default:
		str = fmt.Sprintf("%v - status: succeed, content length: %v, processed in: %v\n",
			out.Url, out.ContentLength, out.ProcessedIn)
	}

	if ordered {
		str = fmt.Sprintf("%d. %v", out.Sequence+1, str)
	}

	return str
}
