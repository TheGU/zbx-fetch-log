package zabbix

type HostInterfaceGetResponse []struct {
	Interfaceid string `json:"interfaceid"`
	Hostid      string `json:"hostid"`
	Main        string `json:"main"`
	Type        string `json:"type"`
	Useip       string `json:"useip"`
	IP          string `json:"ip"`
	DNS         string `json:"dns"`
	Port        string `json:"port"`
	Bulk        string `json:"bulk"`
}

type HostInterfaceGetParams struct {
	Hostids       []string               `json:"hostids"`
	Interfaceids  []string               `json:"interfaceids"`
	Itemids       []string               `json:"itemids"`
	Triggerids    []string               `json:"triggerids"`
	SelectItems   SelectQuery            `json:"selectItems,omitempty"`
	SelectHosts   SelectQuery            `json:"selectHosts,omitempty"`
	LimitSelects  int64                  `json:"limitSelects,omitempty"`
	SortFiels     []string               `json:"sortfield,omitempty"`
	CountOutput   bool                   `json:"countOutput,omitempty"`
	Editable      bool                   `json:"editable,omitempty"`
	ExcludeSearch bool                   `json:"excludeSearch,omitempty"`
	Filter        map[string]interface{} `json:"filter,omitempty"`
	Limit         int64                  `json:"limit,omitempty"`
	Nodeids       []string               `json:"nodeids,omitempty"`
	Output        SelectQuery            `json:"output,omitempty"`
	PreserveKeys  bool                   `json:"preserveKeys,omitempty"`
}

func (c *Session) HostInterfaceGet(params HostInterfaceGetParams) (HostInterfaceGetResponse, error) {
	var respData HostInterfaceGetResponse
	err := c.Get("hostinterface.get", params, &respData)
	return respData, err
}
