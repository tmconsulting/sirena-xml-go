package structs

import (
	"encoding/xml"
)

const (
	// Date and time formats
	Time        = "15:04"
	TimeSec     = "15:04:05"
	Date        = "02.01.2006"
	TimeDate    = "15:04 02.01.2006"
	DateTime    = "02.01.2006 15:04"
	TimeSecDate = "15:04:05 02.01.2006"
	DateTimeSec = "02.01.2006 15:04:05"
)

// BookingResponse is a Sirena response to <booking> request
type BookingResponse struct {
	Answer  BookingAnswer `xml:"answer"`
	XMLName xml.Name      `xml:"sirena" json:"-"`
}

// BookingAnswer is an <answer> section in Sirena booking response
type BookingAnswer struct {
	Pult     string               `xml:"pult,attr,omitempty"`
	MsgID    int                  `xml:"msgid,attr"`
	Time     string               `xml:"time,attr"`
	Instance string               `xml:"instance,attr"`
	Booking  BookingAnswerBooking `xml:"booking"`
}

// BookingAnswerBooking is a <booking> section in Sirena booking response
type BookingAnswerBooking struct {
	Regnum   string                `xml:"regnum,attr"`
	Agency   string                `xml:"agency,attr"`
	PNR      BookingAnswerPNR      `xml:"pnr"`
	Contacts BookingAnswerContacts `xml:"contacts"`
	Error    *Error                `xml:"error,omitempty"`
}

// BookingAnswerPNR is a <pnr> section in Sirena booking response
type BookingAnswerPNR struct {
	RegNum            string                      `xml:"regnum"`
	UTCTimeLimit      string                      `xml:"utc_timelimit"`
	TimeLimit         string                      `xml:"timelimit"`
	LatinRegistration bool                        `xml:"latin_registration"`
	Version           int                         `xml:"version"`
	Segments          []BookingAnswerPNRSegment   `xml:"segments>segment"`
	Passengers        []BookingAnswerPNRPassenger `xml:"passengers>passenger"`
	Prices            BookingAnswerPNRPrices      `xml:"prices"`
}

// BookingAnswerPNRSegment is a <segment> section in <booking> answer
type BookingAnswerPNRSegment struct {
	ID           int                        `xml:"id,attr,omitempty"`
	Company      string                     `xml:"company"`
	Flight       string                     `xml:"flight"`
	SubClass     string                     `xml:"subclass"`
	Class        string                     `xml:"class"`
	BaseClass    string                     `xml:"baseclass"`
	SeatCount    int                        `xml:"seatcount"`
	Airplane     string                     `xml:"airplane"`
	Legs         []PNRSegmentLeg            `xml:"legs>leg"`
	Departure    PNRSegmentDepartureArrival `xml:"departure"`
	Arrival      PNRSegmentDepartureArrival `xml:"arrival"`
	Status       PNRSegmentStatus           `xml:"status"`
	FlightTime   string                     `xml:"flightTime"`
	RemoteRecloc string                     `xml:"remote_recloc"`
	Cabin        string                     `xml:"cabin"`
}

// PNRSegmentLeg is a <leg> entry in <segment> section
type PNRSegmentLeg struct {
	Airplane string              `xml:"airplane,attr"`
	Dep      PNRSegmentLegDepArr `xml:"dep"`
	Arr      PNRSegmentLegDepArr `xml:"arr"`
}

// PNRSegmentLegDepArr is <sep> and <arr> entries in <leg> section
type PNRSegmentLegDepArr struct {
	TimeLocal string `xml:"time_local,attr"`
	TimeUTC   string `xml:"time_utc,attr"`
	Term      string `xml:"term,attr"`
	Value     string `xml:",chardata"`
}

// PNRSegmentDepartureArrival is <departure> and <arrival> entries in <segment> section
type PNRSegmentDepartureArrival struct {
	City     string `xml:"city"`
	Airport  string `xml:"airport"`
	Date     string `xml:"date"`
	Time     string `xml:"time"`
	Terminal string `xml:"terminal"`
}

// PNRSegmentStatus is a <status> entry in a <segment> section
type PNRSegmentStatus struct {
	Text   string `xml:"text,attr"`
	Status string `xml:",chardata"`
}

// BookingAnswerPNRPassenger is a <passenger> section in Sirena booking response
type BookingAnswerPNRPassenger struct {
	ID          string               `xml:"id,attr,omitempty"`
	LeadPass    bool                 `xml:"lead_pass,attr"`
	Name        string               `xml:"name"`
	Surname     string               `xml:"surname"`
	Sex         string               `xml:"sex"`
	Birthdate   string               `xml:"birthdate"`
	Age         int                  `xml:"age"`
	DocCode     string               `xml:"doccode"`
	Doc         string               `xml:"doc"`
	PspExpire   string               `xml:"pspexpire"`
	Category    PNRPassengerCategory `xml:"category"`
	DocCountry  string               `xml:"doc_country"`
	Nationality string               `xml:"nationality"`
	Residence   string               `xml:"residence"`
	Contacts    []Contact            `xml:"contacts>contact"`
}

// PNRPassengerCategory is a <category> entry in <passenger> section
type PNRPassengerCategory struct {
	RBM   int    `xml:"rbm,attr"`
	Value string `xml:",chardata"`
}

// BookingAnswerPNRPrices is a <prices> section in <booking> answer
type BookingAnswerPNRPrices struct {
	TickSer      string                  `xml:"tick_ser,attr"`
	FOP          string                  `xml:"fop,attr"`
	Prices       []BookingAnswerPNRPrice `xml:"price"`
	VariantTotal PNRVariantTotal         `xml:"variant_total"`
}

func (b *BookingAnswerPNRPrices) GetTotalPaxCost(paxType string) *float64 {
	// Passenger total cost will be added
	var paxTotalCost = new(float64)

	// Variant object contains objects of all passengers.
	// Since it is necessary to add the price for only one passenger of the transferred type, therefore passengerID is determined and the price is added only for a passenger of that passengerID.
	var passengerID int

	for _, price := range b.Prices {
		if price.Code == paxType {

			// Attach ID of first passenger of appropriate type to passengerID
			passengerID = price.PassengerID
			break
		}
	}

	for _, price := range b.Prices {
		// Check if it is needed passenger type
		if price.Code == paxType && price.PassengerID == passengerID {
			//allocate a new zero-valued paxTotalCost

			*paxTotalCost += price.Total
		}
	}

	return paxTotalCost
}

func (b *BookingAnswerPNRPrices) GetTaxesVariantCost() *float64 {
	// Passenger total cost will be added
	var variantTaxesCost = new(float64)

	for _, price := range b.Prices {

		for _, tax := range price.Taxes {

			*variantTaxesCost += tax.Value.Value
		}
	}

	return variantTaxesCost
}

func (b *BookingAnswerPNRPrices) GetFareVariantCost() *float64 {
	// Passenger total cost will be added
	var variantFareCost = new(float64)

	for _, price := range b.Prices {

		*variantFareCost += price.Fare.Value.Value
	}

	return variantFareCost
}

func (b *BookingAnswerPNRPrices) GetRawTaxPax(paxType string) []PNRPriceTax {
	// Passenger total cost will be added
	var paxRawTaxes []PNRPriceTax

	// Variant object contains objects of all passengers.
	// Since it is necessary to add the price for only one passenger of the transferred type, therefore passengerID is determined and the price is added only for a passenger of that passengerID.
	var passengerID int

	for _, price := range b.Prices {
		if price.Code == paxType {

			// Attach ID of first passenger of appropriate type to passengerID
			passengerID = price.PassengerID
			break
		}
	}

	for _, price := range b.Prices {
		// Check if it is needed passenger type
		if price.Code == paxType && price.PassengerID == passengerID {
			//allocate a new zero-valued paxTotalCost

		TAXES_LOOP:
			for _, tax := range price.Taxes {

				for _, containsPax := range paxRawTaxes {

					if tax.Value == containsPax.Value {
						continue TAXES_LOOP
					}
				}

				paxRawTaxes = append(paxRawTaxes, tax)
			}
		}
	}

	return paxRawTaxes
}

func (b *BookingAnswerPNRPrices) GetRawVatPax(paxType string) []*Vat {
	// Passenger total cost will be added
	var paxRawVats []*Vat

	// Variant object contains objects of all passengers.
	// Since it is necessary to add the price for only one passenger of the transferred type, therefore passengerID is determined and the price is added only for a passenger of that passengerID.
	var passengerID int

	for _, price := range b.Prices {
		if price.Code == paxType {

			// Attach ID of first passenger of appropriate type to passengerID
			passengerID = price.PassengerID
			break
		}
	}

	for _, price := range b.Prices {
		// Check if it is needed passenger type
		if price.Code == paxType && price.PassengerID == passengerID {
			//allocate a new zero-valued paxTotalCost

			paxRawVats = append(paxRawVats, price.Vat)

		}
	}

	return paxRawVats
}

func (b *BookingAnswerPNRPrices) GetFarePaxCost(paxType string) *float64 {
	// Passenger total cost will be added
	var paxFareCost = new(float64)

	// Variant object contains objects of all passengers.
	// Since it is necessary to add the price for only one passenger of the transferred type, therefore passengerID is determined and the price is added only for a passenger of that passengerID.
	var passengerID int

	for _, price := range b.Prices {
		if price.Code == paxType {

			// Attach ID of first passenger of appropriate type to passengerID
			passengerID = price.PassengerID
			break
		}
	}

	for _, price := range b.Prices {
		// Check if it is needed passenger type
		if price.Code == paxType && price.PassengerID == passengerID {
			//allocate a new zero-valued paxFareCost

			*paxFareCost += price.Fare.Value.Value
		}
	}

	return paxFareCost
}

// GetTaxesPaxCost func return tax amount for passenger of given type
func (b *BookingAnswerPNRPrices) GetTaxesPaxCost(paxType string) *float64 {
	// Passenger total cost will be added
	var paxTaxesCost = new(float64)

	// Variant object contains objects of all passengers.
	// Since it is necessary to add the price for only one passenger of the transferred type, therefore passengerID is determined and the price is added only for a passenger of that passengerID.
	var passengerID int

	for _, price := range b.Prices {
		if price.Code == paxType {

			// Attach ID of first passenger of appropriate type to passengerID
			passengerID = price.PassengerID
			break
		}
	}

	for _, price := range b.Prices {
		// Check if it is needed passenger type
		if price.Code == paxType && price.PassengerID == passengerID {
			//allocate a new zero-valued paxTaxesCost
			for _, tax := range price.Taxes {

				*paxTaxesCost += tax.Value.Value
			}
		}
	}

	return paxTaxesCost
}

func (b *BookingAnswerPNRPrices) GetTotalVariantCost() PNRVariantTotal {
	return b.VariantTotal
}

// BookingAnswerPNRPrice is a <price> entry in Sirena booking response
type BookingAnswerPNRPrice struct {
	SegmentID         int            `xml:"segment-id,attr"`
	PassengerID       int            `xml:"passenger-id,attr"`
	Code              string         `xml:"code,attr"`
	OrigCode          string         `xml:"orig_code,attr"`
	Count             int            `xml:"count,attr"`
	Currency          string         `xml:"currency,attr"`
	TourCode          string         `xml:"tour_code,attr"`
	FC                string         `xml:"fc,attr"`
	Baggage           string         `xml:"baggage,attr"`
	Ticket            string         `xml:"ticket,attr"`
	ValidatingCompany string         `xml:"validating_company,attr"`
	ACCode            string         `xml:"accode,attr"`
	DocType           string         `xml:"doc_type,attr"`
	DocID             string         `xml:"doc_id,attr"`
	Brand             string         `xml:"brand,attr"`
	Fare              PNRPriceFare   `xml:"fare"`
	Taxes             []PNRPriceTax  `xml:"taxes>tax"`
	PaymentInfo       PNRPaymentInfo `xml:"payment_info>payment"`
	Total             float64        `xml:"total"`
	Vat               *Vat           `xml:"vat"`
}

// PNRPriceFare is a <fare> entry in a <price> section
type PNRPriceFare struct {
	Remark      string           `xml:"remark,attr"`
	FareExpDate string           `xml:"fare_expdate,attr"`
	Value       PNRPriceValue    `xml:"value"`
	Code        PNRPriceFareCode `xml:"code"`
}

// PNRPriceValue is a <value> entry in a <fare> section
type PNRPriceValue struct {
	Value    float64 `xml:",chardata"`
	Currency string  `xml:"currency,attr"`
}

// PNRPriceFareCode is a <code> entry in a <fare> section
type PNRPriceFareCode struct {
	Code     string `xml:",chardata"`
	BaseCode string `xml:"base_code,attr"`
}

// PNRPriceTax is a <tax> entry in a <price> section
type PNRPriceTax struct {
	Owner string        `xml:"owner,attr"`
	Code  string        `xml:"code"`
	Value PNRPriceValue `xml:"value"`
}

// PNRPaymentInfo is a <payment_info> entry in a <price> section
type PNRPaymentInfo struct {
	FOP     string  `xml:"fop,attr"`
	Curr    string  `xml:"curr,attr"`
	Payment float64 `xml:",chardata"`
}

// PNRVariantTotal is a <variant_total> entry in Sirena booking response
type PNRVariantTotal struct {
	Currency string  `xml:"currency,attr"`
	Value    float64 `xml:",chardata"`
}

// BookingAnswerContacts is a <contacts> entry in <booking> section
type BookingAnswerContacts struct {
	Contacts []Contact        `xml:"contact"`
	Customer ContactsCustomer `xml:"customer"`
}

// ContactsCustomer is a <customer> entry in <contacts> section
type ContactsCustomer struct {
	FirstName string `xml:"firstname"`
	LastName  string `xml:"lastname"`
}
