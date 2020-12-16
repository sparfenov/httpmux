package workers

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
)

type loggerInterface interface {
	Debugf(msg string, v ...interface{})
	Infof(msg string, v ...interface{})
	Errorf(msg string, v ...interface{})
}

type Result struct {
	Data  string
	Error error
}

type URLFetcher struct {
	URL        string
	Result     Result
	httpClient *http.Client
	logger     loggerInterface
}

func NewURLFetcher(url string, httpClient *http.Client, l loggerInterface) *URLFetcher {
	return &URLFetcher{
		URL:        url,
		httpClient: httpClient,
		logger:     l,
	}
}

func (w *URLFetcher) Process(ctx context.Context) {
	w.logger.Debugf("[STARTED] process for url: %s", w.URL)

	req, err := http.NewRequestWithContext(ctx, "GET", w.URL, nil)
	if err != nil {
		w.Result = Result{
			Error: fmt.Errorf("error while creating request for url '%s': %s", w.URL, err),
		}

		return
	}

	resp, err := w.httpClient.Do(req)
	if err != nil {
		w.Result = Result{
			Error: fmt.Errorf("error while requesting url '%s': %s", w.URL, err),
		}

		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		w.Result = Result{
			Error: fmt.Errorf("failed to read url's response body %s: %s", w.URL, err),
		}

		return
	}

	w.Result = Result{
		Data: string(body),
	}

	w.logger.Debugf("[FINISHED] process for url: %s", w.URL)
}
