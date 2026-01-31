package api

import "fmt"

type Domain struct {
	Domain       string `json:"domain"`
	Status       string `json:"status"`
	TLD          string `json:"tld"`
	CreateDate   string `json:"createDate"`
	ExpireDate   string `json:"expireDate"`
	SecurityLock string `json:"securityLock"`
	WhoisPrivacy string `json:"whoisPrivacy"`
	AutoRenew    string `json:"autoRenew"`
	NotLocal     int    `json:"notLocal"`
}

type domainListResponse struct {
	Response
	Domains []Domain `json:"domains"`
}

type domainGetResponse struct {
	Status       string `json:"status"`
	Message      string `json:"message,omitempty"`
	Domain       string `json:"domain"`
	DomainStatus string `json:"domainStatus"`
	TLD          string `json:"tld"`
	CreateDate   string `json:"createDate"`
	ExpireDate   string `json:"expireDate"`
	SecurityLock string `json:"securityLock"`
	WhoisPrivacy string `json:"whoisPrivacy"`
	AutoRenew    string `json:"autoRenew"`
	NotLocal     int    `json:"notLocal"`
}

func (c *Client) DomainList(start int) ([]Domain, error) {
	body := c.authBodyWith(map[string]any{})
	if start > 0 {
		body["start"] = start
	}

	var resp domainListResponse
	err := c.post("/domain/listAll", body, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Domains, nil
}

func (c *Client) DomainGet(domain string) (*Domain, error) {
	var resp domainGetResponse
	err := c.post(fmt.Sprintf("/domain/getDomain/%s", domain), c.authBody(), &resp)
	if err != nil {
		return nil, err
	}
	return &Domain{
		Domain:       resp.Domain,
		Status:       resp.DomainStatus,
		TLD:          resp.TLD,
		CreateDate:   resp.CreateDate,
		ExpireDate:   resp.ExpireDate,
		SecurityLock: resp.SecurityLock,
		WhoisPrivacy: resp.WhoisPrivacy,
		AutoRenew:    resp.AutoRenew,
		NotLocal:     resp.NotLocal,
	}, nil
}

func (c *Client) DomainUpdateNameservers(domain string, nameservers []string) error {
	body := c.authBodyWith(map[string]any{
		"ns": nameservers,
	})
	var resp Response
	return c.post(fmt.Sprintf("/domain/updateNs/%s", domain), body, &resp)
}

func (c *Client) DomainGetNameservers(domain string) ([]string, error) {
	var resp struct {
		Response
		NS []string `json:"ns"`
	}
	err := c.post(fmt.Sprintf("/domain/getNs/%s", domain), c.authBody(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.NS, nil
}

func (c *Client) DomainAddForward(domain, location string, opts ForwardOpts) error {
	body := c.authBodyWith(map[string]any{
		"location": location,
	})
	if opts.Type != "" {
		body["type"] = opts.Type
	}
	if opts.IncludePath {
		body["includePath"] = "yes"
	}
	if opts.Wildcard {
		body["wildcard"] = "yes"
	}
	if opts.Subdomain != "" {
		body["subdomain"] = opts.Subdomain
	}
	var resp Response
	return c.post(fmt.Sprintf("/domain/addUrlForward/%s", domain), body, &resp)
}

type ForwardOpts struct {
	Type        string // temporary or permanent
	IncludePath bool
	Wildcard    bool
	Subdomain   string
}

func (c *Client) DomainGetForwards(domain string) ([]URLForward, error) {
	var resp struct {
		Response
		Forwards []URLForward `json:"forwards"`
	}
	err := c.post(fmt.Sprintf("/domain/getUrlForwarding/%s", domain), c.authBody(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.Forwards, nil
}

type URLForward struct {
	ID          string `json:"id"`
	Subdomain   string `json:"subdomain"`
	Location    string `json:"location"`
	Type        string `json:"type"`
	IncludePath string `json:"includePath"`
	Wildcard    string `json:"wildcard"`
}

func (c *Client) DomainDeleteForward(domain, forwardID string) error {
	var resp Response
	return c.post(fmt.Sprintf("/domain/deleteUrlForward/%s/%s", domain, forwardID), c.authBody(), &resp)
}
