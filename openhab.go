package openhab

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

// OpenHAB is an instance of the OpenHAB client
type OpenHAB struct {
	baseURL    string
	httpClient *http.Client
}

// Item is an item in OpenHAB
type Item struct {
	Link       string
	State      string
	Editable   bool
	Type       string
	Name       string
	Label      string
	Tags       []string
	GroupNames []string
	openHAB    *OpenHAB
}

// StateDescription is a description of an Item.State
type StateDescription struct {
	Pattern  string
	ReadOnly bool
	Options  []StateOption
}

// StateOption is an option for an enumerated item type
type StateOption struct {
	Value string
	Label string
}

// New returns a new OpenHAB
func New(baseURL string, httpClient *http.Client) *OpenHAB {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: time.Second * 5,
		}
	}
	return &OpenHAB{
		baseURL:    strings.TrimRight(baseURL, "/ "),
		httpClient: httpClient,
	}
}

// Items returns all items
func (oh *OpenHAB) Items(httpClient *http.Client) ([]Item, error) {
	c := httpClient
	if c == nil {
		c = oh.httpClient
	}

	var items []Item

	r, err := c.Get(oh.baseURL + "/items")
	if err != nil {
		return items, err
	}

	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(&items)
	if err != nil {
		return items, err
	}

	for i := range items {
		items[i].openHAB = oh
	}

	return items, nil
}

// Set sets a value on an Item
func (i Item) Set(value string, httpClient *http.Client) error {
	c := httpClient
	if c == nil {
		c = i.openHAB.httpClient
	}

	r, err := c.Post(i.Link, "text/plain", bytes.NewBuffer([]byte(value)))
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return nil
}
