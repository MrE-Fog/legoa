package dns01

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

const (
	dnsTemplate = `%s %d IN TXT %q`
)

// DNSProviderManual is an implementation of the ChallengeProvider interface.
type DNSProviderManual struct{}

// NewDNSProviderManual returns a DNSProviderManual instance.
func NewDNSProviderManual() (*DNSProviderManual, error) {
	return &DNSProviderManual{}, nil
}

// Present prints instructions for manually creating the TXT record.
func (*DNSProviderManual) Present(domain, token, keyAuth string) error {
	info := GetChallengeInfo(domain, keyAuth)

	authZone, err := FindZoneByFqdn(info.EffectiveFQDN)
	if err != nil {
		return err
	}

	fmt.Printf("lego: Please create the following TXT record in your %s zone:\n", authZone)
	fmt.Printf(dnsTemplate+"\n", info.EffectiveFQDN, DefaultTTL, info.Value)
	fmt.Printf("lego: Press 'Enter' when you are done\n")

	_, err = bufio.NewReader(os.Stdin).ReadBytes('\n')

	return err
}

// CleanUp prints instructions for manually removing the TXT record.
func (*DNSProviderManual) CleanUp(domain, token, keyAuth string) error {
	info := GetChallengeInfo(domain, keyAuth)

	authZone, err := FindZoneByFqdn(info.EffectiveFQDN)
	if err != nil {
		return err
	}

	fmt.Printf("lego: You can now remove this TXT record from your %s zone:\n", authZone)
	fmt.Printf(dnsTemplate+"\n", info.EffectiveFQDN, DefaultTTL, "...")

	return nil
}

// Sequential All DNS challenges for this provider will be resolved sequentially.
// Returns the interval between each iteration.
func (d *DNSProviderManual) Sequential() time.Duration {
	return DefaultPropagationTimeout
}
