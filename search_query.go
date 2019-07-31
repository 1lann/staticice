package staticice

import (
	"net/url"
	"strconv"
)

// SearchQuery represents a query to search with.
type SearchQuery struct {
	values url.Values
}

// NewSearchQuery creates a new search query.
func NewSearchQuery() *SearchQuery {
	return &SearchQuery{
		values: url.Values{
			"links": []string{"100"},
		},
	}
}

// Manufacturer sets a comma seperated list of manufacturers to search for.
func (q *SearchQuery) Manufacturer(manufacturer string) *SearchQuery {
	q.values["manufacturer"] = []string{manufacturer}
	return q
}

// Model sets a comma seperated list of model names/numbers to search for.
func (q *SearchQuery) Model(model string) *SearchQuery {
	q.values["model"] = []string{model}
	return q
}

// Words sets a space seperated list of words that the listing must contain.
func (q *SearchQuery) Words(words string) *SearchQuery {
	q.values["words"] = []string{words}
	return q
}

// Phrase sets an exact phrase that the listing must contain.
func (q *SearchQuery) Phrase(phrase string) *SearchQuery {
	q.values["phrase"] = []string{phrase}
	return q
}

// ExcludeWords sets a space seperated list of words that the listing must
// not include.
func (q *SearchQuery) ExcludeWords(words string) *SearchQuery {
	q.values["excludewords"] = []string{words}
	return q
}

// Site sets a comma seperated list of strings to specify what the domain
// name must contain.
func (q *SearchQuery) Site(site string) *SearchQuery {
	q.values["site"] = []string{site}
	return q
}

// Query sets the normal search query to use.
func (q *SearchQuery) Query(query string) *SearchQuery {
	q.values["q"] = []string{query}
	return q
}

// MinPrice sets the minimum price for the search query in the local currency.
func (q *SearchQuery) MinPrice(min int) *SearchQuery {
	q.values["price-min"] = []string{strconv.Itoa(min)}
	return q
}

// MaxPrice sets the maximum price for the search query in the local currency.
func (q *SearchQuery) MaxPrice(max int) *SearchQuery {
	q.values["price-max"] = []string{strconv.Itoa(max)}
	return q
}
