package planetscale

import "fmt"

type planetscaleError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (p *planetscaleError) Error() string {
	return fmt.Sprintf("%s: %s", p.Code, p.Message)
}
