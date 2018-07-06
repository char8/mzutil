package monzo

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/char8/mzutil/auth"
)

type BalanceResponse struct {
	Balance    int64  `json:'balance'`
	Currency   string `json:'currency'`
	SpendToday int64  `json:'spend_today'`
}

type AccountResponse struct {
	Id      string    `json:'id'`
	Desc    string    `json:'description'`
	Created time.Time `json:'created'`
}

type AccountsResponse struct {
	Accounts []AccountResponse `json:'accounts'`
}

type WhoAmIResponse struct {
	Authenticated bool   `json:'authenticated'`
	ClientId      string `json:'client_id'`
	UserId        string `json:'user_id'`
}

type Client struct {
	httpClient *http.Client
}

func NewClient(ctx context.Context, a auth.Authenticator) *Client {
	return &Client{
		httpClient: a.NewHttpClient(ctx),
	}
}

func (c *Client) Balance(accountId string) (b BalanceResponse, err error) {
	url := "https://api.monzo.com/balance?account_id=" + accountId
	resp, err := c.httpClient.Get(url)

	if err != nil {
		return
	}

	jd := json.NewDecoder(resp.Body)
	err = jd.Decode(&b)

	return
}

func (c *Client) WhoAmI() (w WhoAmIResponse, err error) {
	url := "https://api.monzo.com/ping/whoami"
	resp, err := c.httpClient.Get(url)

	if err != nil {
		return
	}

	jd := json.NewDecoder(resp.Body)
	err = jd.Decode(&w)

	return
}

func (c *Client) Accounts() (a []AccountResponse, err error) {
	url := "https://api.monzo.com/accounts"
	resp, err := c.httpClient.Get(url)

	if err != nil {
		return
	}

	accs := AccountsResponse{}

	jd := json.NewDecoder(resp.Body)
	err = jd.Decode(&accs)

	if err != nil {
		return
	}

	return accs.Accounts, nil
}
