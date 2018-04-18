package gonaturalist

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	rateLimitExceededStatusCode = 429
)

type Client struct {
	rootUrl       string
	http          *http.Client
	autoRetry     bool
	retryDuration time.Duration
}

type PageHeaders struct {
	TotalEntries int
	PerPage      int
	Page         int
}

func isFailure(code int, validCodes []int) bool {
	for _, item := range validCodes {
		if item == code {
			return false
		}
	}
	return true
}

func shouldRetry(status int) bool {
	return status == http.StatusAccepted || status == http.StatusTooManyRequests
}

func (c *Client) execute(req *http.Request, result interface{}, needsStatus ...int) error {
	for {
		req.Header.Set("Content-Type", "application/json")

		resp, err := c.http.Do(req)
		if err != nil {
			return err
		}

		defer resp.Body.Close()

		if c.autoRetry && shouldRetry(resp.StatusCode) {
			time.Sleep(c.retryDuration)
			continue
		} else if resp.StatusCode != http.StatusOK && isFailure(resp.StatusCode, needsStatus) {
			errorMessage := c.decodeError(resp)
			return errorMessage
		}

		if result != nil {
			if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
				return fmt.Errorf("Decoding body: %v", err)
			}
		}

		break
	}

	return nil
}

func (c *Client) get(url string, result interface{}) (paging *PageHeaders, err error) {
	for {
		resp, err := c.http.Get(url)
		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()

		total := resp.Header["X-Total-Entries"]
		if len(total) > 0 {
			perPage := resp.Header["X-Per-Page"]
			page := resp.Header["X-Page"]
			t, _ := strconv.Atoi(string(total[0]))
			p, _ := strconv.Atoi(string(page[0]))
			pp, _ := strconv.Atoi(string(perPage[0]))
			paging = &PageHeaders{
				TotalEntries: t,
				Page:         p,
				PerPage:      pp,
			}
		}
		if resp.StatusCode == rateLimitExceededStatusCode && c.autoRetry {
			time.Sleep(c.retryDuration)
			continue
		} else if resp.StatusCode != http.StatusOK {
			errorMessage := c.decodeError(resp)
			return nil, errorMessage
		}

		err = json.NewDecoder(resp.Body).Decode(result)
		if err != nil {
			return nil, fmt.Errorf("Decoding body: %v (%s)", err, url)
		}

		break
	}

	return paging, nil
}

func (c *Client) buildUrl(f string, args ...interface{}) string {
	return fmt.Sprintf(c.rootUrl+f, args...)
}

func (c *Client) decodeError(resp *http.Response) error {
	return nil
}

var AcceptableFormats = []string{
	time.RFC3339,
	"2006/01/02 3:04 PM MST",
	"2006-01-02",
}

func TryParseObservedOn(s string) (time.Time, error) {
	str := strings.Trim(s, "\"")
	if str == "null" {
		return time.Time{}, nil
	}
	for _, l := range AcceptableFormats {
		time, err := time.Parse(l, str)
		if err == nil {
			return time, nil
		}
	}
	return time.Time{}, fmt.Errorf("Unable to parse time: %s", s)
}

type NaturalistTime struct {
	time.Time
}

func (t *NaturalistTime) UnmarshalJSON(b []byte) (err error) {
	parsed, err := TryParseObservedOn(string(b))
	if err != nil {
		return err
	}
	t.Time = parsed
	return
}

func (t *NaturalistTime) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte("null"), nil
	}

	return []byte(fmt.Sprintf("\"%s\"", t.Time.Format("2006-01-02 15:04:05"))), nil
}

func (t *NaturalistTime) IsSet() bool {
	return !t.IsZero()
}
