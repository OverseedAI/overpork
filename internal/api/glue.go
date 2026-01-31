package api

import "fmt"

type GlueRecord struct {
	Subdomain string   `json:"subdomain"`
	IPs       []string `json:"ips"`
}

type glueListResponse struct {
	Response
	Records []GlueRecord `json:"records"`
}

func (c *Client) GlueList(domain string) ([]GlueRecord, error) {
	var resp glueListResponse
	err := c.post(fmt.Sprintf("/domain/getGlue/%s", domain), c.authBody(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.Records, nil
}

func (c *Client) GlueCreate(domain, subdomain string, ips []string) error {
	body := c.authBodyWith(map[string]any{
		"ip": ips,
	})
	var resp Response
	return c.post(fmt.Sprintf("/domain/createGlue/%s/%s", domain, subdomain), body, &resp)
}

func (c *Client) GlueUpdate(domain, subdomain string, ips []string) error {
	body := c.authBodyWith(map[string]any{
		"ip": ips,
	})
	var resp Response
	return c.post(fmt.Sprintf("/domain/updateGlue/%s/%s", domain, subdomain), body, &resp)
}

func (c *Client) GlueDelete(domain, subdomain string) error {
	var resp Response
	return c.post(fmt.Sprintf("/domain/deleteGlue/%s/%s", domain, subdomain), c.authBody(), &resp)
}
