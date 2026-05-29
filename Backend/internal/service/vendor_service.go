package service

import (
	"context"
	"regexp"
)

type Vendor struct {
	Name     string
	GSTIN    string
	Verified bool
	Address  string
}

type VendorService struct{}

func NewVendorService() *VendorService {
	return &VendorService{}
}

func (v *VendorService) VerifyGST(gstin string) bool {
	ok, _ := v.VerifyGSTWithFallback(context.Background(), gstin)
	return ok
}

func (v *VendorService) VerifyGSTWithFallback(_ context.Context, gstin string) (bool, string) {
	// Real API integration can replace this call; we keep source explicit for audit transparency.
	if !isValidGSTFormat(gstin) {
		return false, "FORMAT"
	}
	return true, "FORMAT_FALLBACK"
}

func isValidGSTFormat(gstin string) bool {
	pattern := `^[0-9]{2}[A-Z]{5}[0-9]{4}[A-Z]{1}[A-Z0-9]{1}Z[A-Z0-9]{1}$`
	return regexp.MustCompile(pattern).MatchString(gstin)
}
