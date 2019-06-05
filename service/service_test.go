package service

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tmconsulting/sirenaxml-golang-sdk/configuration"
	"github.com/tmconsulting/sirenaxml-golang-sdk/logs"
	"github.com/tmconsulting/sirenaxml-golang-sdk/sdk"
	"github.com/tmconsulting/sirenaxml-golang-sdk/structs"
)

// AvailabilityXML is a test availability XML
func TestService(t *testing.T) {
	t.Run("test no zipped request", func(t *testing.T) {
		customSirenConfig := sc
		testRequest(t, customSirenConfig)
	})
	t.Run("test zipped request/response", func(t *testing.T) {
		customSirenConfig := sc
		customSirenConfig.ZippedMessaging = true
		testRequest(t, customSirenConfig)
	})
	t.Run("test error params", func(t *testing.T) {
		logger := logs.NewNullLog()
		customSirenConfig := sc
		customSirenConfig.ClientID = 1111
		_, err := sdk.NewClient(&customSirenConfig, logger)
		if !assert.Error(t, err) {
			t.FailNow()
		}
	})
}

func TestService_RawRequest(t *testing.T) {
	logger := logs.NewNullLog()
	sdkClient, err := sdk.NewClient(&sc, logger)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	service := NewSKD(sdkClient)
	checkKeyData(t, sdkClient)
	t.Run("success", func(t *testing.T) {
		xmlReq := []byte(`<?xml version="1.0" encoding="UTF-8"?><sirena><query><key_info/></query></sirena>`)
		response, err := service.RawRequest(xmlReq)
		if !assert.NoError(t, err) {
			t.FailNow()
		}
		// Decode XML and make sure Sirena Public Key is returned
		var keyInfoResponse structs.KeyInfoResponse
		err = xml.Unmarshal(response, &keyInfoResponse)
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		if !assert.NotEmpty(t, keyInfoResponse.Answer.KeyInfo.KeyManager.ServerPubliKey) {
			t.FailNow()
		}
	})
}

func checkKeyData(t *testing.T, c SirenaSDK) {
	if !assert.NotEmpty(t, c.GetKeyData().Key) {
		t.FailNow()
	}
	if !assert.NotZero(t, c.GetKeyData().ID) {
		t.FailNow()
	}
}

func TestService_Avalability(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		logger := logs.NewNullLog()
		sdkClient, err := sdk.NewClient(&sc, logger)
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		service := NewSKD(sdkClient)
		checkKeyData(t, sdkClient)
		availabiliteReq := &structs.AvailabilityRequest{
			Query: structs.AvailabilityRequestQuery{
				Availability: structs.Availability{
					Departure: "MOW",
					Arrival:   "LED",
					AnswerParams: structs.AvailabilityAnswerParams{
						ShowFlighttime: true,
					},
				},
			},
		}

		_, err = service.Avalability(availabiliteReq)
		if !assert.NoError(t, err) {
			t.FailNow()
		}
	})
}

func TestService_KeyInfo(t *testing.T) {
	logger := logs.NewNullLog()
	sdkClient, err := sdk.NewClient(&sc, logger)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	service := NewSKD(sdkClient)
	checkKeyData(t, sdkClient)
	t.Run("success", func(t *testing.T) {
		_, err = service.KeyInfo()
		if !assert.NoError(t, err) {
			t.FailNow()
		}
	})
}

func testRequest(t *testing.T, sc sirenaXML.Config) {
	logger := logs.NewNullLog()
	sdkClient, err := sdk.NewClient(&sc, logger)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	service := NewSKD(sdkClient)
	checkKeyData(t, sdkClient)

	var (
		respChan = make(chan *structs.KeyInfoResponse)
		errChan  = make(chan error)
	)
	for i := 0; i < int(sc.MaxConnections); i++ {
		go func() {
			response, err := service.KeyInfo()
			if err != nil {
				errChan <- err
				return
			}
			respChan <- response
		}()
	}

	select {
	case response := <-respChan:
		if !assert.NotEmpty(t, response.Answer.KeyInfo.KeyManager.ServerPubliKey) {
			t.FailNow()
		}
	case err := <-errChan:
		if !assert.NoError(t, err) {
			t.FailNow()
		}
	}

}
