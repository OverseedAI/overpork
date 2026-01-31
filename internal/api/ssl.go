package api

import "fmt"

type SSLBundle struct {
	IntermediateCertificate string `json:"intermediatecertificate"`
	CertificateChain        string `json:"certificatechain"`
	PrivateKey              string `json:"privatekey"`
	PublicKey               string `json:"publickey"`
}

type sslResponse struct {
	Response
	SSLBundle
}

func (c *Client) SSLRetrieve(domain string) (*SSLBundle, error) {
	var resp sslResponse
	err := c.post(fmt.Sprintf("/ssl/retrieve/%s", domain), c.authBody(), &resp)
	if err != nil {
		return nil, err
	}
	return &resp.SSLBundle, nil
}
