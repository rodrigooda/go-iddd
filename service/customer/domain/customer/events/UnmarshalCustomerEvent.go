package events

import (
	"github.com/AntonStoeckl/go-iddd/service/lib"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
	"github.com/cockroachdb/errors"
)

func UnmarshalCustomerEvent(name string, payload []byte, streamVersion uint) (es.DomainEvent, error) {
	switch name {
	case "CustomerRegistered":
		return UnmarshalCustomerRegisteredFromJSON(payload, streamVersion), nil
	case "CustomerEmailAddressConfirmed":
		return UnmarshalCustomerEmailAddressConfirmedFromJSON(payload, streamVersion), nil
	case "CustomerEmailAddressConfirmationFailed":
		return UnmarshalCustomerEmailAddressConfirmationFailedFromJSON(payload, streamVersion), nil
	case "CustomerEmailAddressChanged":
		return UnmarshalCustomerEmailAddressChangedFromJSON(payload, streamVersion), nil
	case "CustomerNameChanged":
		return UnmarshalCustomerNameChangedFromJSON(payload, streamVersion), nil
	case "CustomerDeleted":
		return UnmarshalCustomerDeletedFromJSON(payload, streamVersion), nil
	}

	err := errors.Mark(
		errors.Wrapf(errors.New("event is unknown"), "unmarshalDomainEvent [%s] failed", name),
		lib.ErrUnmarshalingFailed,
	)

	return nil, err
}