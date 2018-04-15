package common

import (
	"fmt"
	"github.com/bobesa/go-domain-util/domainutil"
	"time"
)

type Media struct {
	Id           uint64    `json:"id"`
	UserId       uint64    `json:"user_id,omitempty"`
	Title        string    `json:"title,omitempty"`
	Domain       string    `json:"domain,omitempty"`
	TopDomain    string    `json:"top_domain,omitempty"`
	Desc         string    `json:"desc,omitempty"`
	Verified     uint      `json:"verified"`
	OnlineStatus uint      `json:"online_status"`
	Identity     string    `json:"identity,omitempty"`
	VerifyDNS    string    `json:"verify_dns,omitempty"`
	VerifyURL    string    `json:"verify_url,omitempty"`
	DNSValue     string    `json:"dns_value,omitempty"`
	InsertedAt   time.Time `json:"inserted_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (this Media) Complete() Media {
	this.TopDomain = this.GetTopDomain()
	this.VerifyDNS = this.GetVerifyDNS()
	this.VerifyURL = this.GetVerifyURL()
	this.DNSValue = this.GetDNSValue()
	return this
}

func (this Media) GetTopDomain() string {
	return domainutil.Domain(this.Domain)
}

func (this Media) GetVerifyDNS() string {
	return fmt.Sprintf("tokenmama%s.%s", this.Identity, this.TopDomain)
}

func (this Media) GetDNSValue() string {
	return fmt.Sprintf("%s.dns.tokenmama.io", this.Identity)
}

func (this Media) GetVerifyURL() string {
	return fmt.Sprintf("%s/%s.txt", this.Domain, this.Identity)
}
