package api

import "fmt"

type DNSSECRecord struct {
	KeyTag     string `json:"keyTag"`
	Algorithm  string `json:"algorithm"`
	DigestType string `json:"digestType"`
	Digest     string `json:"digest"`
	PublicKey  string `json:"publicKey,omitempty"`
	Flags      string `json:"flags,omitempty"`
}

type dnssecListResponse struct {
	Response
	Records []DNSSECRecord `json:"records"`
}

func (c *Client) DNSSECList(domain string) ([]DNSSECRecord, error) {
	var resp dnssecListResponse
	err := c.post(fmt.Sprintf("/dns/getDnssecRecords/%s", domain), c.authBody(), &resp)
	if err != nil {
		return nil, err
	}
	return resp.Records, nil
}

func (c *Client) DNSSECCreate(domain string, record DNSSECRecord) error {
	body := c.authBodyWith(map[string]any{
		"keyTag":     record.KeyTag,
		"algorithm":  record.Algorithm,
		"digestType": record.DigestType,
		"digest":     record.Digest,
	})
	if record.PublicKey != "" {
		body["publicKey"] = record.PublicKey
	}
	if record.Flags != "" {
		body["flags"] = record.Flags
	}

	var resp Response
	return c.post(fmt.Sprintf("/dns/createDnssecRecord/%s", domain), body, &resp)
}

func (c *Client) DNSSECDelete(domain, keyTag string) error {
	var resp Response
	return c.post(fmt.Sprintf("/dns/deleteDnssecRecord/%s/%s", domain, keyTag), c.authBody(), &resp)
}
