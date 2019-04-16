package domain

import (
	"errors"
	"go-iddd/customer/domain/valueobjects"
	"go-iddd/shared"
)

/*** The Customer behaviour method to apply the ConfirmEmailAddress command ***/

func (customer *customer) confirmEmailAddress(confirmEmailAddress ConfirmEmailAddress) error {
	var err error

	if customer.emailAddress.IsConfirmed() {
		return nil
	}

	if !customer.emailAddress.Equals(confirmEmailAddress.EmailAddress()) {
		return errors.New("customer - emailAddress can not be confirmed because it has changed")
	}

	if customer.emailAddress, err = customer.emailAddress.Confirm(confirmEmailAddress.ConfirmationHash()); err != nil {
		return err
	}

	return nil
}

/*** The ConfirmEmailAddress command itself - struct, factory, own getters, shared.Command getters ***/

type ConfirmEmailAddress interface {
	ID() valueobjects.ID
	EmailAddress() valueobjects.EmailAddress
	ConfirmationHash() valueobjects.ConfirmationHash

	shared.Command
}

type confirmEmailAddress struct {
	id               valueobjects.ID
	emailAddress     valueobjects.EmailAddress
	confirmationHash valueobjects.ConfirmationHash
}

func NewConfirmEmailAddress(
	id valueobjects.ID,
	emailAddress valueobjects.EmailAddress,
	confirmationHash valueobjects.ConfirmationHash,
) (*confirmEmailAddress, error) {

	command := &confirmEmailAddress{
		id:               id,
		emailAddress:     emailAddress,
		confirmationHash: confirmationHash,
	}

	if err := shared.AssertAllPropertiesAreNotNil(command); err != nil {
		return nil, err
	}

	return command, nil
}

func (confirmEmailAddress *confirmEmailAddress) ID() valueobjects.ID {
	return confirmEmailAddress.id
}

func (confirmEmailAddress *confirmEmailAddress) EmailAddress() valueobjects.EmailAddress {
	return confirmEmailAddress.emailAddress
}

func (confirmEmailAddress *confirmEmailAddress) ConfirmationHash() valueobjects.ConfirmationHash {
	return confirmEmailAddress.confirmationHash
}

func (confirmEmailAddress *confirmEmailAddress) Identifier() string {
	return confirmEmailAddress.id.String()
}

func (confirmEmailAddress *confirmEmailAddress) CommandName() string {
	return shared.BuildNameFor(confirmEmailAddress)
}