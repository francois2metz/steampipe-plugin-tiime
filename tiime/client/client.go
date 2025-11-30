package tiime_client

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/auth0/go-auth0/v2/authentication"
	"github.com/auth0/go-auth0/v2/authentication/oauth"
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

type Line struct {
	ID                    int     `json:"id"`
	Description           string  `json:"description"`
	LineAmount            float32 `json:"line_amount"`
	InvoicingCategoryType string  `json:"invoicing_category_type"`
	UnitAmount            float32 `json:"unit_amount"`
	Quantity              float32 `json:"quantity"`
}

type ListQueryOpts struct {
	Status       string
	EmissionDate string
}

type Tag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Company struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Siret      string `json:"siret"`
	Street     string `json:"street"`
	PostalCode string `json:"postal_code"`
	City       string `json:"city"`
	Country    string `json:"country"`
	RcsCity    string `json:"rcs_city"`
	Logo       string `json:"logo"`
	LegalForm  string `json:"legal_form"`
}

type Invoice struct {
	ID                  int64   `json:"id"`
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
	Title               string  `json:"title"`
	Status              string  `json:"status"`
	Lines               []Line  `json:"lines"`
	Tags                []Tag   `json:"tags"`
}

type Quote struct {
	ID                  int64   `json:"id"`
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
	Title               string  `json:"title"`
	Status              string  `json:"status"`
	Lines               []Line  `json:"lines"`
	Tags                []Tag   `json:"tags"`
}

type Client2 struct {
	ID                    int     `json:"id"`
	Slug                  string  `json:"slug"`
	Name                  string  `json:"name"`
	Address               string  `json:"address"`
	PostalCode            string  `json:"postal_code"`
	City                  string  `json:"city"`
	Email                 string  `json:"email"`
	Phone                 string  `json:"phone"`
	PaymentStatus         string  `json:"payment_status"`
	BalanceIncludingTaxes float64 `json:"balance_including_taxes"`
	BilledIncludingTaxes  float32 `json:"billed_including_taxes"`
	BilledExcludingTaxes  float32 `json:"billed_excluding_taxes"`
}

type BankAccount struct {
	ID                  int     `json:"id"`
	SynchronizationDate string  `json:"synchronization_date"`
	LastPushDate        string  `json:"last_push_date"`
	ShortBankName       string  `json:"short_bank_name"`
	BankName            string  `json:"bank_name"`
	Disableable         bool    `json:"disableable"`
	AuthorizedBalance   float64 `json:"authorized_balance"`
	PendingBalance      float64 `json:"pending_balance"`
	Closurable          bool    `json:"closurable"`
	Iban                string  `json:"iban"`
	Bic                 string  `json:"bic"`
	Name                string  `json:"name"`
	Enabled             bool    `json:"enabled"`
	BalanceAmount       float64 `json:"balance_amount"`
	BalanceCurrency     string  `json:"balance_currency"`
	BalanceDate         string  `json:"balance_date"`
	Closed              bool    `json:"closed"`
	IsWallet            bool    `json:"is_wallet"`
}

type PaginationOpts struct {
	Start int
	End   int
}

type Pagination struct {
	CurrentStart int
	CurrentEnd   int
	Max          string
}

type APIError struct {
	ErrorStatus      int    `json:"error"`
	ErrorCode        string `json:"error_code"`
	ErrorDescription string `json:"error_description"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error: %d %s", e.ErrorStatus, e.ErrorCode)
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
		SetCommonPathParam("company_id", strconv.Itoa(c.config.CompanyID)).
		SetCommonErrorResult(&APIError{}).
		OnAfterResponse(func(client *req.Client, resp *req.Response) error {
			if resp.Err != nil { // There is an underlying error, e.g. network error or unmarshal error.
				return nil
			}
			if apiErr, ok := resp.ErrorResult().(*APIError); ok {
				// Server returns an error message, convert it to human-readable go error.
				resp.Err = apiErr
				return nil
			}
			// Corner case: neither an error state response nor a success state response,
			// dump content to help troubleshoot.
			if !resp.IsSuccessState() {
				return fmt.Errorf("bad response, raw dump:\n%s", resp.Dump())
			}
			return nil
		})

	return c, nil
}

func (c *Client) GetCompanies(ctx context.Context) (companies []Company, err error) {
	c.Get("/companies").
		SetBearerAuthToken(c.token.AccessToken).
		Do(ctx).
		Into(&companies)
	return
}

func (c *Client) GetInvoices(ctx context.Context, opts ListQueryOpts, paginationOpts PaginationOpts) (invoices []Invoice, pagination Pagination, err error) {
	res := c.Get("/companies/{company_id}/invoices").
		SetBearerAuthToken(c.token.AccessToken).
		SetHeader("Range", formatRange(paginationOpts)).
		SetQueryParams(getListQueryParams(opts)).
		Do(ctx)
	pagination, err = handlePagination(res)
	if err != nil {
		return
	}
	err = res.Into(&invoices)
	return
}

func (c *Client) GetQuotes(ctx context.Context, opts ListQueryOpts, paginationOpts PaginationOpts) (quotes []Quote, pagination Pagination, err error) {
	res := c.Get("/companies/{company_id}/quotations").
		SetBearerAuthToken(c.token.AccessToken).
		SetHeader("Range", formatRange(paginationOpts)).
		SetQueryParams(getListQueryParams(opts)).
		Do(ctx)
	pagination, err = handlePagination(res)
	if err != nil {
		return
	}
	err = res.Into(&quotes)
	return
}

func getListQueryParams(opts ListQueryOpts) map[string]string {
	var query = make(map[string]string)
	if opts.Status != "" {
		query["status"] = opts.Status
	}
	if opts.EmissionDate != "" {
		query["emission_date"] = opts.EmissionDate
	}
	return query
}

func (c *Client) GetInvoice(ctx context.Context, id int64) (invoice Invoice, err error) {
	err = c.Get("/companies/{company_id}/invoices/{id}").
		SetPathParam("id", strconv.FormatInt(id, 10)).
		SetBearerAuthToken(c.token.AccessToken).
		Do(ctx).
		Into(&invoice)
	return
}

func (c *Client) GetQuote(ctx context.Context, id int64) (quote Quote, err error) {
	err = c.Get("/companies/{company_id}/quotations/{id}").
		SetPathParam("id", strconv.FormatInt(id, 10)).
		SetBearerAuthToken(c.token.AccessToken).
		Do(ctx).
		Into(&quote)
	return
}

func (c *Client) GetClients(ctx context.Context, paginationOpts PaginationOpts) (clients []Client2, pagination Pagination, err error) {
	res := c.Get("/companies/{company_id}/clients").
		SetBearerAuthToken(c.token.AccessToken).
		SetHeader("Range", formatRange(paginationOpts)).
		Do(ctx)
	pagination, err = handlePagination(res)
	if err != nil {
		return
	}
	err = res.Into(&clients)
	return
}

func (c *Client) GetClient(ctx context.Context, id int64) (client Client2, err error) {
	err = c.Get("/companies/{company_id}/clients/{id}").
		SetPathParam("id", strconv.FormatInt(id, 10)).
		SetBearerAuthToken(c.token.AccessToken).
		Do(ctx).
		Into(&client)
	return
}

func formatRange(paginationOpts PaginationOpts) string {
	return fmt.Sprintf("items=%d-%d", paginationOpts.Start, paginationOpts.End)
}

func handlePagination(res *req.Response) (pagination Pagination, err error) {
	contentRange := strings.NewReader(res.GetHeader("Content-range"))
	_, err = fmt.Fscanf(contentRange, "items %d-%d/%s", &pagination.CurrentStart, &pagination.CurrentEnd, &pagination.Max)
	return
}

func (c *Client) GetBankAccounts(ctx context.Context) (bankAccounts []BankAccount, err error) {
	res := c.Get("/companies/{company_id}/bank_accounts").
		SetBearerAuthToken(c.token.AccessToken).
		Do(ctx)
	err = res.Into(&bankAccounts)

	for i, bankAccount := range bankAccounts {
		bankAccounts[i].SynchronizationDate = fixDate(bankAccount.SynchronizationDate)
		bankAccounts[i].LastPushDate = fixDate(bankAccount.LastPushDate)
		bankAccounts[i].BalanceDate = fixDate(bankAccount.BalanceDate)
	}

	return
}

func fixDate(s string) string {
	return strings.ReplaceAll(s, " ", "T")
}
