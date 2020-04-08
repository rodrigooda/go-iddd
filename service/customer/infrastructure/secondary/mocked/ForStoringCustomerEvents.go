// Code generated by mockery v1.0.0. DO NOT EDIT.

// +build test

package mocked

import es "github.com/AntonStoeckl/go-iddd/service/lib/es"
import mock "github.com/stretchr/testify/mock"
import values "github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"

// ForStoringCustomerEvents is an autogenerated mock type for the ForStoringCustomerEvents type
type ForStoringCustomerEvents struct {
	mock.Mock
}

// Add provides a mock function with given fields: recordedEvents, id
func (_m *ForStoringCustomerEvents) Add(recordedEvents es.DomainEvents, id values.CustomerID) error {
	ret := _m.Called(recordedEvents, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(es.DomainEvents, values.CustomerID) error); ok {
		r0 = rf(recordedEvents, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateStreamFrom provides a mock function with given fields: recordedEvents, id
func (_m *ForStoringCustomerEvents) CreateStreamFrom(recordedEvents es.DomainEvents, id values.CustomerID) error {
	ret := _m.Called(recordedEvents, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(es.DomainEvents, values.CustomerID) error); ok {
		r0 = rf(recordedEvents, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// EventStreamFor provides a mock function with given fields: id
func (_m *ForStoringCustomerEvents) EventStreamFor(id values.CustomerID) (es.DomainEvents, error) {
	ret := _m.Called(id)

	var r0 es.DomainEvents
	if rf, ok := ret.Get(0).(func(values.CustomerID) es.DomainEvents); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(es.DomainEvents)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(values.CustomerID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Purge provides a mock function with given fields: id
func (_m *ForStoringCustomerEvents) Purge(id values.CustomerID) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(values.CustomerID) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
