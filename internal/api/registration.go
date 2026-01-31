package api

import "fmt"

type DomainCreateOpts struct {
	Years        int
	Coupon       string
	Nameservers  []string
	WhoisPrivacy bool
	AutoRenew    bool
}

type domainCreateResponse struct {
	Response
	Domain string `json:"domain"`
}

func (c *Client) DomainRegister(domain string, opts DomainCreateOpts) error {
	body := c.authBodyWith(map[string]any{})

	if opts.Years > 0 {
		body["years"] = opts.Years
	}
	if opts.Coupon != "" {
		body["coupon"] = opts.Coupon
	}
	if len(opts.Nameservers) > 0 {
		body["ns"] = opts.Nameservers
	}
	if opts.WhoisPrivacy {
		body["whoisPrivacy"] = "yes"
	}
	if opts.AutoRenew {
		body["autoRenew"] = "yes"
	}

	var resp domainCreateResponse
	return c.post(fmt.Sprintf("/domain/create/%s", domain), body, &resp)
}

func (c *Client) DomainSetAutoRenew(domain string, enabled bool) error {
	status := "disable"
	if enabled {
		status = "enable"
	}
	body := c.authBodyWith(map[string]any{
		"autoRenew": status,
	})
	var resp Response
	return c.post(fmt.Sprintf("/domain/updateAutoRenew/%s", domain), body, &resp)
}
