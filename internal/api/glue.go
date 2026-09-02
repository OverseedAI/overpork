package api

import (
	"encoding/json"
	"fmt"
	"strings"
)

type GlueRecord struct {
	Subdomain string   `json:"subdomain"`
	IPs       []string `json:"ips"`
}

type glueListResponse struct {
	Response
	Hosts []glueHost `json:"hosts"`
}

type glueHost struct {
	Hostname string
	V4       []string
	V6       []string
}

func (h *glueHost) UnmarshalJSON(data []byte) error {
	var tuple []json.RawMessage
	if err := json.Unmarshal(data, &tuple); err != nil {
		return err
	}
	if len(tuple) != 2 {
		return fmt.Errorf("expected host tuple with 2 items, got %d", len(tuple))
	}
	if err := json.Unmarshal(tuple[0], &h.Hostname); err != nil {
		return err
	}
	var addresses struct {
		V4 []string `json:"v4"`
		V6 []string `json:"v6"`
	}
	if err := json.Unmarshal(tuple[1], &addresses); err != nil {
		return err
	}
	h.V4 = addresses.V4
	h.V6 = addresses.V6
	return nil
}

func (c *Client) GlueList(domain string) ([]GlueRecord, error) {
	var resp glueListResponse
	err := c.post(fmt.Sprintf("/domain/getGlue/%s", domain), c.authBody(), &resp)
	if err != nil {
		return nil, err
	}
	records := make([]GlueRecord, len(resp.Hosts))
	for i, host := range resp.Hosts {
		subdomain := host.Hostname
		suffix := "." + domain
		if len(subdomain) > len(suffix) && strings.EqualFold(subdomain[len(subdomain)-len(suffix):], suffix) {
			subdomain = subdomain[:len(subdomain)-len(suffix)]
		}
		records[i] = GlueRecord{
			Subdomain: subdomain,
			IPs:       append(host.V4, host.V6...),
		}
	}
	return records, nil
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
