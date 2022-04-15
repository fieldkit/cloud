package partners

import (
	"context"
	"net/http"
	"strings"
)

var partners = []*PartnerConfiguration{
	{
		HostPattern: "floodnet",
		Sharing: PartnerSharingConfiguration{
			Title: "FloodNet Chart", // TODO i18n
		},
	},
	{
		HostPattern: "",
		Sharing: PartnerSharingConfiguration{
			Title: "FieldKit Chart", // TODO i18n
		},
	},
}

type PartnerSharingConfiguration struct {
	Title string `json:"title"`
}

// TODO Consider embedding same JSON used by partners.ts so we can staticly link
// to the JS/HTML and also use server side.
type PartnerConfiguration struct {
	HostPattern string                      `json:"host_pattern"`
	Sharing     PartnerSharingConfiguration `json:"sharing"`
}

func GetPartnerForRequest(ctx context.Context, req *http.Request) (*PartnerConfiguration, error) {
	for _, partner := range partners {
		if partner.HostPattern == "" || strings.Contains(req.Host, partner.HostPattern) {
			return partner, nil
		}
	}

	return partners[len(partners)-1], nil
}
