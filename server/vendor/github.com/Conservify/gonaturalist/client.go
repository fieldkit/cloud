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

type pageHeaders struct {
	totalEntries int
	perPage      int
	page         int
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

func (c *Client) get(url string, result interface{}) (paging *pageHeaders, err error) {
	for {
		resp, err := c.http.Get(url)
		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()

		total := resp.Header["X-Total-Entries"]
		if len(total) > 0 {
			perPage := resp.Header["X-Per-Page"][0]
			page := resp.Header["X-Page"][0]
			t, _ := strconv.Atoi(string(total[0]))
			p, _ := strconv.Atoi(string(page[0]))
			pp, _ := strconv.Atoi(string(perPage[0]))
			paging = &pageHeaders{
				totalEntries: t,
				page:         p,
				perPage:      pp,
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

type NaturalistTime struct {
	time.Time
}

var AcceptableFormats = []string{
	time.RFC3339,
	"2006-01-02",
}

func (t *NaturalistTime) UnmarshalJSON(b []byte) (err error) {
	str := strings.Trim(string(b), "\"")
	if str == "null" {
		t.Time = time.Time{}
		return
	}
	for _, l := range AcceptableFormats {
		t.Time, err = time.Parse(l, str)
		if err == nil {
			break
		}
	}
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
