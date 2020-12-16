package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/sparfenov/httpmux/internal/workers"
	"github.com/sparfenov/httpmux/pkg/logger"
	"github.com/sparfenov/httpmux/pkg/workerpool"
	"net/http"
	"strconv"
	"time"
)

type URLMuxHandler struct {
	// maximum URL count allowed in incoming request
	MaxURLCountToProcess int

	// maximum outbound request count
	ExternalRequestLimit int

	// outbound request timeout
	ExternalRequestTimeout time.Duration

	Logger logger.Interface
}

// incoming request body struct
type RequestBody struct {
	URLs []string `json:"urls"`
}

// result struct for processed url
type ResultItem struct {
	URL  string `json:"url"`
	Data string `json:"data"`
}

// overall result for request
type Result struct {
	Items []ResultItem `json:"items"`
	Error string       `json:"error"`
	Code  int          `json:"-"`
}

func (h URLMuxHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	req.Response = &http.Response{Request: req}

	if req.Method != http.MethodPost {
		req.Response.StatusCode = http.StatusMethodNotAllowed
		http.Error(rw, http.StatusText(req.Response.StatusCode), req.Response.StatusCode)

		return
	}

	opResult := make(chan Result, 1)

	ctx, cancel := context.WithCancel(req.Context())
	defer cancel()

	go func() {
		reqData := RequestBody{}

		err := h.extractRequestData(*req, &reqData)
		if err != nil {
			opResult <- Result{Error: err.Error(), Code: http.StatusBadRequest}

			return
		}

		if len(reqData.URLs) == 0 || len(reqData.URLs) > h.MaxURLCountToProcess {
			opResult <- Result{
				Error: "URL count must be greater than: 0 and lower or equals: " + strconv.Itoa(h.MaxURLCountToProcess),
				Code:  http.StatusBadRequest,
			}

			return
		}

		result := Result{Code: http.StatusOK}

		result.Items, err = h.fetchURLs(ctx, reqData.URLs)
		if err != nil {
			opResult <- Result{Error: err.Error(), Code: http.StatusBadGateway}

			return
		}

		opResult <- result
	}()

	select {
	case <-req.Context().Done():
		h.Logger.Infof("request has been canceled")
		cancel()

		return
	case r := <-opResult:
		responseBody, err := json.Marshal(r)
		if err != nil {
			h.Logger.Errorf("failed to marshal response json: %s", err)

			req.Response.StatusCode = http.StatusInternalServerError
			http.Error(rw, http.StatusText(req.Response.StatusCode), req.Response.StatusCode)

			return
		}

		req.Response.StatusCode = r.Code
		rw.WriteHeader(r.Code)
		_, _ = rw.Write(responseBody)
	}
}

func (h URLMuxHandler) extractRequestData(req http.Request, reqBody *RequestBody) error {
	if req.ContentLength == 0 {
		return errors.New("body should be non empty JSON")
	}

	defer req.Body.Close()

	p := make([]byte, 0, req.ContentLength)
	buf := bytes.NewBuffer(p)

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		h.Logger.Errorf("failed to read request body: %s", err)

		return errors.New("failed to read body")
	}

	err = json.Unmarshal(buf.Bytes(), &reqBody)
	if err != nil {
		h.Logger.Errorf("failed to unmarshal request body to json: %s", err)

		return errors.New("bad json input")
	}

	return nil
}

func (h URLMuxHandler) fetchURLs(ctx context.Context, urls []string) ([]ResultItem, error) {
	cp := workerpool.NewWorkerPool(h.ExternalRequestLimit)

	jobs := make([]workerpool.Job, 0, len(urls))

	for _, url := range urls {
		httpClient := http.Client{Timeout: h.ExternalRequestTimeout}
		jobs = append(jobs, workers.NewURLFetcher(url, &httpClient, h.Logger))
	}

	finishedJobs := cp.Run(ctx, jobs)

	results := make([]ResultItem, 0, len(finishedJobs))

	for _, job := range finishedJobs {
		j := job.(*workers.URLFetcher)
		if j.Result.Error != nil {
			return results, j.Result.Error
		}

		results = append(results, ResultItem{
			URL:  j.URL,
			Data: j.Result.Data,
		})
	}

	return results, nil
}
