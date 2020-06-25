package azure

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

// Parameters for azure
type Parameters struct {
	ClientID       string `json:"clientId"`
	ClientSecret   string `json:"clientSecret"`
	TenantID       string `json:"tenantId"`
	SubscriptionID string `json:"subscriptionId"`
	ResourceGroup  string `json:"resourceGroup"`
	DNSName        string `json:"dnsName"`
	RecordSetName  string `json:"recordSetName"`
	Environment    string `json:"environment"`
	RecordType     string
}

// Validate parameters validate
func (p *Parameters) Validate() bool {
	if p.ClientID == "" || p.ClientSecret == "" || p.TenantID == "" || p.SubscriptionID == "" || p.ResourceGroup == "" || p.DNSName == "" || p.RecordSetName == "" {
		return false
	}
	if p.Environment == "" {
		p.Environment = "https://management.azure.com/"
	}
	if p.RecordType == "" {
		p.RecordType = "AAAA"
	}
	return true
}

var (
	client = &http.Client{
		Timeout: 20 * time.Second,
	}
)

// ParseCredentialFromEnv parse client secret credential
func ParseCredentialFromEnv(p *Parameters) {
	authFile := os.Getenv("AZURE_AUTH_LOCATION")
	if authFile != "" {
		auth, err := ioutil.ReadFile(authFile)
		if err == nil {
			json.Unmarshal(auth, p)
		}
	} else {
		p.ClientID = os.Getenv("AZURE_CLIENT_ID")
		p.ClientSecret = os.Getenv("AZURE_CLIENT_SECRET")
		p.TenantID = os.Getenv("AZURE_TENANT_ID")
		p.SubscriptionID = os.Getenv("AZURE_SUBSCRIPTION_ID")
	}
}

func login(p *Parameters) string {
	uri := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/token", p.TenantID)
	resp, err := client.PostForm(uri, url.Values{
		"resource":      []string{p.Environment},
		"client_id":     []string{p.ClientID},
		"client_secret": []string{p.ClientSecret},
		"client_info":   []string{"1"},
		"grant_type":    []string{"client_credentials"},
	})
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var tokens map[string]interface{}
	err = json.Unmarshal(respBody, &tokens)
	if err != nil {
		panic(err)
	}
	return tokens["access_token"].(string)
}

// UpdateDNS Update DNS
func UpdateDNS(p *Parameters, records map[string][]map[string]string) {
	if !p.Validate() {
		ParseCredentialFromEnv(p)
		if !p.Validate() {
			panic(fmt.Errorf("Parameter validation error: %+v", p))
		}
	}

	token := login(p)
	uri := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/dnsZones/%s/%s/%s?api-version=2018-05-01", p.Environment, p.SubscriptionID, p.ResourceGroup, p.DNSName, p.RecordType, p.RecordSetName)

	properties := map[string]interface{}{
		"TTL": 600,
	}
	for key, value := range records {
		properties[key] = value
	}

	body := map[string]interface{}{
		"properties": properties,
	}
	bodyBytes, _ := json.Marshal(body)

	req, _ := http.NewRequest("PUT", uri, bytes.NewBuffer(bodyBytes))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	c, _ := ioutil.ReadAll(resp.Body)
	var m map[string]interface{}
	json.Unmarshal(c, &m)
	fmt.Println(m)
}
