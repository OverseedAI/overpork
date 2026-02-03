package api

import "fmt"

const PricingURL = "https://porkbun.com/api/json/v3"

type Pricing struct {
	Registration    string `json:"registration"`
	Renewal         string `json:"renewal"`
	Transfer        string `json:"transfer"`
	Coupons         any    `json:"coupons"`
	SpecialType     string `json:"specialType,omitempty"`
	SpecialDiscount string `json:"specialDiscount,omitempty"`
}

type pricingResponse struct {
	Response
	Pricing map[string]Pricing `json:"pricing"`
}

func (c *Client) PricingList() (map[string]Pricing, error) {
	var resp pricingResponse
	// Pricing endpoint uses porkbun.com (not api.porkbun.com)
	err := c.postURL(PricingURL+"/pricing/get", map[string]string{}, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Pricing, nil
}

func (c *Client) DomainCheck(domain string) (bool, float64, error) {
	var resp struct {
		Response
		Available string `json:"avail"`
		Price     string `json:"price"`
	}
	// Auth not required for availability check
	err := c.post("/domain/checkDomain/"+domain, map[string]string{}, &resp)
	if err != nil {
		return false, 0, err
	}

	available := resp.Available == "yes"
	var price float64
	if resp.Price != "" {
		// Price is a string, parse it
		_, _ = fmt.Sscanf(resp.Price, "%f", &price)
	}
	return available, price, nil
}
