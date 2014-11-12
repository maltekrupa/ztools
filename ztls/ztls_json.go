package ztls

import (
	"encoding/json"
	"errors"
)

type encodedCertificates struct {
	Certificates    [][]byte `json:"certificates"`
	Valid           bool     `json:"valid"`
	ValidationError *string  `json:"validation_error"`
	CommonName      *string  `json:"common_name"`
	AltNames        []string `json:"alt_names"`
	Issuer          *string  `json:"issuer"`
}

func (ec *encodedCertificates) FromZTLS(c *Certificates) *encodedCertificates {
	ec.Certificates = c.Certificates
	ec.Valid = c.Valid
	if c.ValidationError != nil {
		s := c.ValidationError.Error()
		ec.ValidationError = &s
	}
	if c.CommonName != "" {
		ec.CommonName = &c.CommonName
	}
	ec.AltNames = c.AltNames
	if c.Issuer != "" {
		ec.Issuer = &c.Issuer
	}
	return ec
}

func (c *Certificates) FromEncoded(ec *encodedCertificates) *Certificates {
	c.Certificates = ec.Certificates
	c.Valid = ec.Valid
	if ec.ValidationError != nil {
		c.ValidationError = errors.New(*ec.ValidationError)
	}
	if ec.CommonName != nil {
		c.CommonName = *ec.CommonName
	}
	c.AltNames = ec.AltNames
	if ec.Issuer != nil {
		c.Issuer = *ec.Issuer
	}
	return c
}

func (c *Certificates) MarshalJSON() ([]byte, error) {
	ec := new(encodedCertificates).FromZTLS(c)
	return json.Marshal(ec)
}

func (c *Certificates) UnmarshalJSON(b []byte) error {
	ec := new(encodedCertificates)
	if err := json.Unmarshal(b, ec); err != nil {
		return err
	}
	c.FromEncoded(ec)
	return nil
}
