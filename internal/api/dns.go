package api

import "fmt"

// dnsByNameTypeEndpoint builds the endpoint for the *ByNameType DNS
// endpoints. The subdomain segment is omitted for root-domain records,
// since the API does not accept a trailing slash there.
func dnsByNameTypeEndpoint(action, domain, recordType, subdomain string) string {
	endpoint := fmt.Sprintf("/dns/%s/%s/%s", action, domain, recordType)
	if subdomain != "" {
		endpoint += "/" + subdomain
	}
	return endpoint
}

type DNSRecord struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Content string `json:"content"`
	TTL     string `json:"ttl"`
	Prio    string `json:"prio,omitempty"`
	Notes   string `json:"notes,omitempty"`
}

type dnsListResponse struct {
	Response
	Records []DNSRecord `json:"records"`
}

type dnsCreateResponse struct {
	Response
	ID int64 `json:"id"`
}

func (c *Client) DNSList(domain string) ([]DNSRecord, error) {
	var resp dnsListResponse
	err := c.post(fmt.Sprintf("/dns/retrieve/%s", domain), c.authBody(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.Records, nil
}

func (c *Client) DNSListByType(domain, recordType string) ([]DNSRecord, error) {
	var resp dnsListResponse
	err := c.post(fmt.Sprintf("/dns/retrieveByNameType/%s/%s", domain, recordType), c.authBody(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.Records, nil
}

func (c *Client) DNSListByTypeAndSubdomain(domain, recordType, subdomain string) ([]DNSRecord, error) {
	var resp dnsListResponse
	err := c.post(dnsByNameTypeEndpoint("retrieveByNameType", domain, recordType, subdomain), c.authBody(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.Records, nil
}

func (c *Client) DNSCreate(domain, recordType, content string, opts DNSCreateOpts) (int64, error) {
	body := c.authBodyWith(map[string]any{
		"type":    recordType,
		"content": content,
	})
	if opts.Name != "" {
		body["name"] = opts.Name
	}
	if opts.TTL != "" {
		body["ttl"] = opts.TTL
	}
	if opts.Prio != "" {
		body["prio"] = opts.Prio
	}

	var resp dnsCreateResponse
	err := c.post(fmt.Sprintf("/dns/create/%s", domain), body, &resp)
	if err != nil {
		return 0, err
	}
	return resp.ID, nil
}

type DNSCreateOpts struct {
	Name string
	TTL  string
	Prio string
}

func (c *Client) DNSUpdate(domain, recordID, recordType, content string, opts DNSCreateOpts) error {
	body := c.authBodyWith(map[string]any{
		"type":    recordType,
		"content": content,
	})
	if opts.Name != "" {
		body["name"] = opts.Name
	}
	if opts.TTL != "" {
		body["ttl"] = opts.TTL
	}
	if opts.Prio != "" {
		body["prio"] = opts.Prio
	}

	var resp Response
	return c.post(fmt.Sprintf("/dns/edit/%s/%s", domain, recordID), body, &resp)
}

func (c *Client) DNSUpdateByTypeAndSubdomain(domain, recordType, subdomain, content string, opts DNSCreateOpts) error {
	body := c.authBodyWith(map[string]any{
		"type":    recordType,
		"content": content,
	})
	if opts.Name != "" {
		body["name"] = opts.Name
	}
	if opts.TTL != "" {
		body["ttl"] = opts.TTL
	}
	if opts.Prio != "" {
		body["prio"] = opts.Prio
	}

	var resp Response
	return c.post(dnsByNameTypeEndpoint("editByNameType", domain, recordType, subdomain), body, &resp)
}

func (c *Client) DNSDelete(domain, recordID string) error {
	var resp Response
	return c.post(fmt.Sprintf("/dns/delete/%s/%s", domain, recordID), c.authBody(), &resp)
}

func (c *Client) DNSDeleteByTypeAndSubdomain(domain, recordType, subdomain string) error {
	var resp Response
	return c.post(dnsByNameTypeEndpoint("deleteByNameType", domain, recordType, subdomain), c.authBody(), &resp)
}
