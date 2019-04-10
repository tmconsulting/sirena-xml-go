package sirena

import "encoding/xml"

// SetAgentDiscountRequest is a set_agent_discount request
type SetAgentDiscountRequest struct {
	Query   SetAgentDiscountQuery `xml:"query"`
	XMLName xml.Name              `xml:"sirena"`
}

// SetAgentDiscountQuery is a <query> section in set_agent_discount request
type SetAgentDiscountQuery struct {
	SetAgentDiscount *SetAgentDiscount `xml:"set_agent_discount"`
}

// SetAgentDiscount is a body of set_agent_discount request
type SetAgentDiscount struct {
	Regnum        *SetAgentDiscountRegnum `xml:"regnum"`
	Unit          []*SetAgentDiscountUnit `xml:"unit"`
	RequestParams OrderRequestParams      `xml:"request_params"`
}

// SetAgentDiscountRegnum is a Regnum (PNR number and version) in set_agent_discount request
type SetAgentDiscountRegnum struct {
	Version int    `xml:"version,attr"`
	Value   string `xml:",chardata"`
}

// SetAgentDiscountUnit is a <unit> element of the set_agent_discount request
type SetAgentDiscountUnit struct {
	SetAgentDiscountFare *SetAgentDiscountFare `xml:"fare"`
}

// SetAgentDiscountFare is a fare element in set_agent_discount request
type SetAgentDiscountFare struct {
	Discount int    `xml:"discount,attr"`
	Brand    string `xml:"brand,attr"`
	Value    string `xml:",chardata"`
}

// SetAgentDiscountRequestParams is a <request_params> section in set_agent_discount request
type SetAgentDiscountRequestParams struct {
	TickSer string `xml:"tick_ser"`
}
