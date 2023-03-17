package neverbounce

import (
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"net/http"
	"net/url"
)

type CheckEmailConfig struct {
	Email       string
	AddressInfo bool
	CreditsInfo bool
	Timeout     *int64
}

type CheckEmailResult struct {
	Status              string   `json:"status"`
	Result              string   `json:"result"`
	Message             string   `json:"message"`
	Flags               []string `json:"flags"`
	SuggestedCorrection string   `json:"suggested_correction"`
	AddressInfo         *struct {
		OriginalEmail   string `json:"original_email"`
		NormalizedEmail string `json:"normalized_email"`
		Addr            string `json:"addr"`
		Alias           string `json:"alias"`
		Host            string `json:"host"`
		Fqdn            string `json:"fqdn"`
		Domain          string `json:"domain"`
		Subdomain       string `json:"subdomain"`
		Tld             string `json:"tld"`
	} `json:"address_info"`
	CreditsInfo *struct {
		PaidCreditsUsed      int64 `json:"paid_credits_used"`
		FreeCreditsUsed      int64 `json:"free_credits_used"`
		PaidCreditsRemaining int64 `json:"paid_credits_remaining"`
		FreeCreditsRemaining int64 `json:"free_credits_remaining"`
	} `json:"credits_info"`
	ExecutionTime int64 `json:"execution_time"`
}

func (service *Service) CheckEmail(cfg *CheckEmailConfig) (*CheckEmailResult, *errortools.Error) {
	if cfg == nil {
		return nil, errortools.ErrorMessage("CheckEmailConfig must not be nil")
	}

	var values = url.Values{}

	values.Set("email", cfg.Email)
	if cfg.AddressInfo {
		values.Set("address_info", "1")
	}
	if cfg.CreditsInfo {
		values.Set("credits_info", "1")
	}
	if cfg.Timeout != nil {
		values.Set("timeout", fmt.Sprintf("%v", *cfg.Timeout))
	}

	var result CheckEmailResult

	var requestConfig = go_http.RequestConfig{
		Method:        http.MethodGet,
		Url:           service.url(fmt.Sprintf("single/check?%s", values.Encode())),
		ResponseModel: &result,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	if result.Status != "success" {
		return nil, errortools.ErrorMessagef("%s: %s", result.Status, result.Message)
	}

	return &result, nil
}
