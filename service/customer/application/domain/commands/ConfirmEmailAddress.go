package commands

import (
	"go-iddd/service/customer/application/domain/values"
	"go-iddd/service/lib"
	"reflect"
	"strings"

	"github.com/cockroachdb/errors"
)

type ConfirmEmailAddress struct {
	customerID       values.CustomerID
	emailAddress     values.EmailAddress
	confirmationHash values.ConfirmationHash
	isValid          bool
}

func NewConfirmEmailAddress(
	customerID string,
	emailAddress string,
	confirmationHash string,
) (ConfirmEmailAddress, error) {

	customerIDValue, err := values.BuildCustomerID(customerID)
	if err != nil {
		return ConfirmEmailAddress{}, err
	}

	emailAddressValue, err := values.BuildEmailAddress(emailAddress)
	if err != nil {
		return ConfirmEmailAddress{}, err
	}

	confirmationHashValue, err := values.BuildConfirmationHash(confirmationHash)
	if err != nil {
		return ConfirmEmailAddress{}, err
	}

	confirmEmailAddress := ConfirmEmailAddress{
		customerID:       customerIDValue,
		emailAddress:     emailAddressValue,
		confirmationHash: confirmationHashValue,
		isValid:          true,
	}

	return confirmEmailAddress, nil
}

func (confirmEmailAddress ConfirmEmailAddress) CustomerID() values.CustomerID {
	return confirmEmailAddress.customerID
}

func (confirmEmailAddress ConfirmEmailAddress) EmailAddress() values.EmailAddress {
	return confirmEmailAddress.emailAddress
}

func (confirmEmailAddress ConfirmEmailAddress) ConfirmationHash() values.ConfirmationHash {
	return confirmEmailAddress.confirmationHash
}

func (confirmEmailAddress ConfirmEmailAddress) ShouldBeValid() error {
	if !confirmEmailAddress.isValid {
		err := errors.Newf("%s: is not valid", confirmEmailAddress.commandName())

		return errors.Mark(err, lib.ErrCommandIsInvalid)
	}

	return nil
}

func (confirmEmailAddress ConfirmEmailAddress) commandName() string {
	commandType := reflect.TypeOf(confirmEmailAddress).String()
	commandTypeParts := strings.Split(commandType, ".")
	commandName := commandTypeParts[len(commandTypeParts)-1]

	return strings.Title(commandName)
}

func (confirmEmailAddress ConfirmEmailAddress) IsCommand() bool {
	return true
}
