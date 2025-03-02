package tiime_client

import (
	"context"
	"strconv"

	"github.com/auth0/go-auth0/authentication"
	"github.com/auth0/go-auth0/authentication/oauth"
	req "github.com/imroc/req/v3"
)

type ClientConfig struct {
	Email     string
	Password  string
	CompanyID int
}

type Client struct {
	*req.Client
	config  ClientConfig
	token   *oauth.TokenSet
	authAPI *authentication.Authentication
}

type Invoice struct {
	ID                  int     `json:"id"`
	ClientID            int     `json:"client_id"`
	CompiledNumber      string  `json:"compiled_number"`
	Number              int     `json:"number"`
	EmissionDate        string  `json:"emission_date"`
	Template            string  `json:"template"`
	Color               string  `json:"color"`
	ClientName          string  `json:"client_name"`
	TotalExcludingTaxes float32 `json:"total_excluding_taxes"`
	TotalIncludingTaxes float32 `json:"total_including_taxes"`
	Comment             string  `json:"comment"`
}

func New(ctx context.Context, config ClientConfig) (*Client, error) {
	domain := "auth0.tiime.fr"
	clientID := "iEbsbe3o66gcTBfGRa012kj1Rb6vjAND"

	authAPI, err := authentication.New(
		ctx,
		domain,
		authentication.WithClientID(clientID),
	)
	if err != nil {
		return nil, err
	}

	tokenSet, err := authAPI.OAuth.LoginWithPassword(ctx, oauth.LoginWithPasswordRequest{
		Username: config.Email,
		Password: config.Password,
		Scope:    "openid email",
		Audience: "https://chronos/",
		Realm:    "Chronos-prod-db",
	}, oauth.IDTokenValidationOptions{})

	if err != nil {
		return nil, err
	}

	c := &Client{
		Client:  req.C(),
		config:  config,
		authAPI: authAPI,
		token:   tokenSet,
	}

	c.Client.
		SetBaseURL("https://chronos-api.tiime-apps.com/v1").
		SetCommonPathParam("company_id", strconv.Itoa(c.config.CompanyID))

	return c, nil
}

func (c *Client) GetInvoices(ctx context.Context) (invoices []Invoice, err error) {
	err = c.Get("/companies/{company_id}/invoices").
		SetBearerAuthToken(c.token.AccessToken).
		Do(ctx).
		Into(&invoices)
	return
}

func (c *Client) GetInvoice(ctx context.Context, id int64) (invoice Invoice, err error) {
	err = c.Get("/companies/{company_id}/invoices/{id}").
		SetPathParam("id", strconv.FormatInt(id, 10)).
		SetBearerAuthToken(c.token.AccessToken).
		Do(ctx).
		Into(invoice)
	return
}
