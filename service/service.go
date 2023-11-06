package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	RequestTimeOut = time.Second * 3

	ErrNotExpectedStatusCode = errors.New("status code is not expected")
	ErrNoResponse            = errors.New("no response from url")
)

func Handle(urls []string, sequentially bool) {
	if sequentially {
		KeepSequence(urls)
	} else {
		FastHandle(urls)
	}
}

type UrlOut struct {
	Url           string
	ErrorMsg      string
	ContentLength int64
	ProcessedIn   float64
}

func KeepSequence(urls []string) {
	out := make([]UrlOut, len(urls))

	mu := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	wg.Add(len(urls))

	for i, val := range urls {
		go func(idx int, url string) {
			defer wg.Done()

			now := time.Now()
			lenght, err := Get(url)
			procT := time.Since(now).Seconds()

			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				out[idx].ErrorMsg = err.Error()
			}
			out[idx].ContentLength = lenght
			out[idx].Url = url
			out[idx].ProcessedIn = procT
		}(i, val)
	}

	wg.Wait()

	for i, val := range out {
		if val.ErrorMsg != "" {
			fmt.Printf("%d. %v - status: failed, error: %v, processed in: %v\n",
				i+1, val.Url, val.ErrorMsg, val.ProcessedIn)
		} else {
			fmt.Printf("%d. %v - status: succeed, content length: %v, processed in: %v\n",
				i+1, val.Url, val.ContentLength, val.ProcessedIn)
		}
	}
}

func FastHandle(urls []string) {
	wg := &sync.WaitGroup{}
	wg.Add(len(urls))
	for _, val := range urls {
		go func(url string) {
			defer wg.Done()

			now := time.Now()
			lenght, err := Get(url)
			procT := time.Since(now).Seconds()
			if err != nil {
				fmt.Printf("%v - status: failed, error: %v, processed in: %v\n", url, err.Error(), procT)
			} else {
				fmt.Printf("%v - status: succeed, content length: %v, processed in: %v\n", url, lenght, procT)
			}
		}(val)
	}

	wg.Wait()
}

func Get(url string) (int64, error) {
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
		return 0, ErrNotExpectedStatusCode
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
