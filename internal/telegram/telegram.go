package telegram

import (
	"encoding/json"
	"go-todobot/internal/domain/models"
	"go-todobot/lib/e"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

type Client struct {
	host     string
	basePath string
	client   *http.Client
}

const (
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
)

func New(host, token string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   &http.Client{},
	}
}

func (c *Client) Updates(offset, limit int) ([]models.Update, error) {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest(getUpdatesMethod, q)
	if err != nil {
		return nil, err
	}

	var r models.UpdateResponse
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, err
	}

	return r.Result, nil
}

func (c *Client) SendMessage(chatID int, text string) error {
	q := url.Values{}
	log.Println("sending message:", text, chatID)
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)

	_, err := c.doRequest(sendMessageMethod, q)
	if err != nil {
		return e.Wrap("cant send message", err)
	}

	return nil
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) doRequest(method string, query url.Values) (data []byte, err error) {
	defer func() { err = e.NilWrap("failed to do request", err) }()
	log.Println("method and query:", method, query)

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	r, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	r.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(r)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
