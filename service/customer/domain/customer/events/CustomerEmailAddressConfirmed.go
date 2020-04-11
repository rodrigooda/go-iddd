package events

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
)

type CustomerEmailAddressConfirmed struct {
	customerID   values.CustomerID
	emailAddress values.EmailAddress
	meta         es.EventMeta
}

func BuildCustomerEmailAddressConfirmed(
	customerID values.CustomerID,
	emailAddress values.EmailAddress,
	streamVersion uint,
) CustomerEmailAddressConfirmed {

	event := CustomerEmailAddressConfirmed{
		customerID:   customerID,
		emailAddress: emailAddress,
	}

	event.meta = es.BuildEventMeta(event, streamVersion)

	return event
}

func RebuildCustomerEmailAddressConfirmed(
	customerID values.CustomerID,
	emailAddress values.EmailAddress,
	meta es.EventMeta,
) CustomerEmailAddressConfirmed {

	event := CustomerEmailAddressConfirmed{
		customerID:   customerID,
		emailAddress: emailAddress,
		meta:         meta,
	}

	return event
}

func (event CustomerEmailAddressConfirmed) CustomerID() values.CustomerID {
	return event.customerID
}

func (event CustomerEmailAddressConfirmed) EmailAddress() values.EmailAddress {
	return event.emailAddress
}

func (event CustomerEmailAddressConfirmed) Meta() es.EventMeta {
	return event.meta
}

func (event CustomerEmailAddressConfirmed) IndicatesAnError() (bool, string) {
	return false, ""
}
