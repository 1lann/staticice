package staticice

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Client represents a client to access staticICE.
type Client struct {
	client *http.Client
}

// NewClient returns a new Client using the provided http client.
func NewClient(c *http.Client) *Client {
	return &Client{
		client: c,
	}
}

// Region represents a staticICE region.
type Region string

// Available staticICE regions.
const (
	RegionAU Region = "https://www.staticice.com.au"
	RegionNZ Region = "https://www.staticice.co.nz"
	RegionUK Region = "https://www.staticice.co.uk"
	RegionUS Region = "https://www.staticice.com"
)

// ItemEntry represents an entry from staticICE's search.
type ItemEntry struct {
	Description string
	Seller      string
	Link        string
	LastUpdated time.Time
	Price       float64
}

var datePattern = regexp.MustCompile(`\| updated: (.+)`)

// Search performs a search query on the provided region. The first 100
// results will be returned.
func (c *Client) Search(region Region,
	query *SearchQuery) ([]*ItemEntry, error) {
	req, err := http.NewRequest("GET", string(region)+"/cgi-bin/search.cgi?"+
		query.values.Encode(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("staticice: http status not OK: " + resp.Status)
	}

	rd := bufio.NewReader(resp.Body)

	var results []*ItemEntry

	for {
		data, err := rd.ReadBytes('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		if !bytes.HasPrefix(data,
			[]byte(`<tr valign="top"><td align="left"><a`)) {
			continue
		}

		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(
			append([]byte("<table>"),
				append(data, []byte("</tr></table>")...)...)))
		if err != nil {
			return nil, err
		}

		var entry ItemEntry

		entry.Price, err = strconv.ParseFloat(
			doc.Find(`td[align="left"] > a`).Text()[1:], 64)
		if err != nil {
			return nil, err
		}

		entry.Seller = strings.TrimSpace(
			doc.Find(`td[valign="bottom"] > font > a`).Text())

		rawLink, _ := doc.Find(`td[align="left"] > a`).Attr("href")
		u, err := url.ParseRequestURI(rawLink)
		if err != nil {
			return nil, err
		}

		entry.Link = u.Query()["newurl"][0]

		entry.Description = doc.Find(`td[valign="bottom"]`).Contents().
			First().Text()

		matches := datePattern.FindAllStringSubmatch(
			doc.Find(`td[valign="bottom"] > font`).Text(), 1)
		if len(matches) == 0 {
			return nil, errors.New("staticice: no date found")
		}

		entry.LastUpdated, err = time.Parse("02-01-2006",
			strings.TrimSpace(matches[0][1]))
		if err != nil {
			return nil, err
		}

		results = append(results, &entry)
	}

	return results, nil
}
