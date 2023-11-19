package zabbix

import (
	"fmt"
	"testing"
)

func TestHosts(t *testing.T) {
	session := GetTestSession(t)

	params := HostGetParams{}
	params.OutputFields = []string{"hostid"}

	hosts, err := session.GetHosts(params)
	if err != nil {
		t.Fatalf("Error getting Hosts: %v", err)
	}

	if len(hosts) == 0 {
		t.Fatal("No Hosts found")
	}

	for i, host := range hosts {
		if host.HostID == "" {
			t.Fatalf("Host %d returned in response body has no Host ID", i)
		}
		fmt.Println(host.HostID)
	}

	t.Logf("Validated %d Hosts", len(hosts))
}
