package main

import (
	"flag"
	"syscall"

	"github.com/chentanyi/ddns/azure"
	"github.com/chentanyi/go-utils/realip"
)

var (
	p    = &azure.Parameters{}
	ipv4 bool
)

func init() {
	// rootCmd := &cobra.Command{}
	// rootCmd.PersistentFlags().StringVarP(&p.ClientID, "username", "u", "", "Client Id")
	// rootCmd.PersistentFlags().StringVarP(&p.ClientSecret, "password", "p", "", "Client Secret")
	// rootCmd.PersistentFlags().StringVarP(&p.TenantID, "tenant", "t", "", "Tenant Id")
	// rootCmd.PersistentFlags().StringVarP(&p.SubscriptionID, "subscription", "s", "", "Subscription Id")
	// rootCmd.PersistentFlags().StringVarP(&p.ResourceGroup, "group", "g", "", "Resource Group")
	// rootCmd.PersistentFlags().StringVarP(&p.DNSName, "name", "n", "", "Dns Name")
	// rootCmd.PersistentFlags().StringVarP(&p.RecordSetName, "record", "r", "", "Record Sets Name")
	// rootCmd.PersistentFlags().StringVarP(&p.Environment, "environment", "e", "", "Scheme://Host/ for management")
	// rootCmd.PersistentFlags().BoolVarP(&ipv4, "ipv4", "4", false, "Use ipv4 or not")
	// if err := rootCmd.Execute(); err != nil {
	// 	panic(err)
	// }
	flag.StringVar(&p.ClientID, "u", "", "Client Id")
	flag.StringVar(&p.ClientSecret, "p", "", "Client Secret")
	flag.StringVar(&p.TenantID, "t", "", "Tenant Id")
	flag.StringVar(&p.SubscriptionID, "s", "", "Subscription Id")
	flag.StringVar(&p.ResourceGroup, "g", "", "Resource Group")
	flag.StringVar(&p.DNSName, "n", "", "Dns Name")
	flag.StringVar(&p.RecordSetName, "r", "", "Record Sets Name")
	flag.StringVar(&p.Environment, "e", "", "Scheme://Host/ for management")
	flag.BoolVar(&ipv4, "4", false, "Use ipv4 or not")
	flag.Parse()
}

func main() {
	records := map[string][]map[string]string{}
	if ipv4 {
		p.RecordType = "A"
		ips, err := realip.GetRealIP(syscall.AF_INET)
		if err != nil {
			panic(err)
		}
		for _, ip := range ips {
			records["ARecords"] = append(records["ARecords"], map[string]string{"ipv4Address": ip})
		}
	} else {
		p.RecordType = "AAAA"
		ips, err := realip.GetRealIP(syscall.AF_INET6)
		if err != nil {
			panic(err)
		}
		for _, ip := range ips {
			records["AAAARecords"] = append(records["AAAARecords"], map[string]string{"ipv6Address": ip})
		}
	}
	azure.UpdateDNS(p, records)
}
