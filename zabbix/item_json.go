package zabbix

import (
	"fmt"
	"strconv"
)

// jItem is a private map for the Zabbix API Host object.
// See: https://www.zabbix.com/documentation/4.0/manual/api/reference/item/get
type jItem struct {
	ItemID               string `json:"itemid"`
	ItemType             string `json:"type,omitempty"`
	SNMP_OID             string `json:"snmp_oid,omitempty"`
	ItemName             string `json:"name,omitempty"`
	Key                  string `json:"key_"`
	HostID               string `json:"hostid,omitempty"`
	Delay                string `json:"delay,omitempty"`
	History              string `json:"history,omitempty"`
	Trends               string `json:"trends,omitempty"`
	Status               string `json:"status,omitempty"`
	TrapperHosts         string `json:"trapper_hosts,omitempty"`
	Units                string `json:"units,omitempty"`
	SNMPv3SecurityName   string `json:"snmpv3_securityname,omitempty"`
	SNMPv3SecurityLevel  string `json:"snmpv3_securitylevel,omitempty"`
	SNMPv3AuthProtocol   string `json:"snmpv3_authprotocol,omitempty"`
	SNMPv3AuthPassphrase string `json:"snmpv3_authpassphrase,omitempty"`
	SNMPv3PrivProtocol   string `json:"snmpv3_privprotocol,omitempty"`
	SNMPv3PrivPassphrase string `json:"snmpv3_privpassphrase,omitempty"`
	Formula              string `json:"formula,omitempty"`
	Error                string `json:"error,omitempty"`
	LastError            string `json:"lasterror,omitempty"`
	LastLogSize          string `json:"lastlogsize,omitempty"`
	LogTimeFmt           string `json:"logtimefmt,omitempty"`
	TemplateID           string `json:"templateid,omitempty"`
	ValuemapID           string `json:"valuemapid,omitempty"`
	Params               string `json:"params,omitempty"`
	IPMIPath             string `json:"ipmi_sensor,omitempty"`
	Authtype             string `json:"authtype,omitempty"`
	Username             string `json:"username,omitempty"`
	Password             string `json:"password,omitempty"`
	PublicKey            string `json:"publickey,omitempty"`
	PrivateKey           string `json:"privatekey,omitempty"`
	Mtime                string `json:"mtime,omitempty"`
	Flags                string `json:"flags,omitempty"`
	InterfaceID          string `json:"interfaceid,omitempty"`
	Port                 string `json:"port,omitempty"`
	ItemDescr            string `json:"description,omitempty"`
	InventoryLink        string `json:"inventory_link,omitempty"`
	Lifetime             string `json:"lifetime,omitempty"`
	SNMPv3ContextName    string `json:"snmpv3_contextname,omitempty"`
	JmxEndpoint          string `json:"jmx_endpoint,omitempty"`
	MasterItemID         string `json:"master_itemid,omitempty"`
	Timeout              string `json:"timeout,omitempty"`
	URL                  string `json:"url,omitempty"`
	// QueryFields          []string `json:"query_fields,omitempty"`
	Posts           string `json:"posts,omitempty"`
	StatusCodes     string `json:"status_codes,omitempty"`
	FollowRedirects string `json:"follow_redirects,omitempty"`
	PostType        string `json:"post_type,omitempty"`
	HttpProxy       string `json:"http_proxy,omitempty"`
	// Headers              []string `json:"headers,omitempty"`
	RetrieveMode   string `json:"retrieve_mode,omitempty"`
	RequestMethod  string `json:"request_method,omitempty"`
	OutputFormat   string `json:"output_format,omitempty"`
	SslCertFile    string `json:"ssl_cert_file,omitempty"`
	SslKeyFile     string `json:"ssl_key_file,omitempty"`
	SslKeyPassword string `json:"ssl_key_password,omitempty"`
	VerifyPeer     string `json:"verify_peer,omitempty"`
	VerifyHost     string `json:"verify_host,omitempty"`
	AllowTraps     string `json:"allow_traps,omitempty"`
	Discover       string `json:"discover,omitempty"`
	LastClock      string `json:"lastclock,omitempty"`
	LastValue      string `json:"lastvalue,omitempty"`
	PrevValue      string `json:"prevvalue,omitempty"`
	LastValueType  string `json:"value_type,omitempty"`
}

// Item returns a native Go Item struct mapped from the given JSON Item data.
func (c *jItem) Item() (*Item, error) {
	var err error
	item := &Item{}
	// item.HostID, err = strconv.Atoi(c.HostID)
	// if err != nil {
	// 	return nil, fmt.Errorf("Error parsing Host ID: %v", err)
	// }
	item.HostID = c.HostID
	// item.ItemID, err = strconv.Atoi(c.ItemID)
	// if err != nil {
	// 	return nil, fmt.Errorf("Error parsing Item ID: %v", err)
	// }
	item.ItemID = c.ItemID
	item.ItemName = c.ItemName
	item.ItemDescr = c.ItemDescr

	item.LastClock, _ = strconv.ParseInt(c.LastClock, 10, 64)
	// item.LastClock, err = strconv.Atoi(c.LastClock)
	// if err != nil {
	// 	return nil, fmt.Errorf("Error parsing Item LastClock: %v", err)
	// }
	item.LastValue = c.LastValue
	item.PrevValue = c.PrevValue

	item.LastValueType, err = strconv.Atoi(c.LastValueType)
	if err != nil {
		return nil, fmt.Errorf("Error parsing Item LastValueType: %v", err)
	}

	// New fields
	item.Key = c.Key
	item.ItemType, _ = strconv.Atoi(c.ItemType)
	// item.ItemType, err = strconv.Atoi(c.ItemType)
	// if err != nil {
	// 	return nil, fmt.Errorf("Error parsing Item Type: %v", err)
	// }
	item.Status, _ = strconv.Atoi(c.Status)
	// item.Status, err = strconv.Atoi(c.Status)
	// if err != nil {
	// 	return nil, fmt.Errorf("Error parsing Item Status: %v", err)
	// }
	item.Units = c.Units
	item.TemplateID, _ = strconv.Atoi(c.TemplateID)
	// item.TemplateID, err = strconv.Atoi(c.TemplateID)
	// if err != nil {
	// 	return nil, fmt.Errorf("Error parsing Template ID: %v", err)
	// }
	item.MasterItemID, _ = strconv.Atoi(c.MasterItemID)
	// item.MasterItemID, err = strconv.Atoi(c.MasterItemID)
	// if err != nil {
	// 	return nil, fmt.Errorf("Error parsing Master Item ID: %v", err)
	// }
	return item, err
}

// jItems is a slice of jItems structs.
type jItems []jItem

// Items returns a native Go slice of Items mapped from the given JSON ITEMS
// data.
func (c jItems) Items() ([]Item, error) {
	if c != nil {
		items := make([]Item, len(c))
		for i, jitem := range c {
			item, err := jitem.Item()
			if err != nil {
				return nil, fmt.Errorf("Error unmarshalling Item %d in JSON data: %v", i, err)
			}
			items[i] = *item
		}

		return items, nil
	}

	return nil, nil
}
