package api

import "fmt"

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
	err := c.post("/pricing/get", map[string]string{}, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Pricing, nil
}

func (c *Client) DomainCheck(domain string) (bool, float64, error) {
	var resp struct {
		Response
		Domain struct {
			Available string `json:"avail"`
			Price     string `json:"price"`
		} `json:"response"`
	}
	err := c.post("/domain/checkDomain/"+domain, c.authBody(), &resp)
	if err != nil {
		return false, 0, err
	}

	available := resp.Domain.Available == "yes"
	var price float64
	if resp.Domain.Price != "" {
		// Price is a string, parse it
		_, _ = fmt.Sscanf(resp.Domain.Price, "%f", &price)
	}
	return available, price, nil
}
