package eventstore

import (
	"database/sql"
	"go-iddd/service/customer/application/domain/values"
	"go-iddd/service/lib"
	"go-iddd/service/lib/es"
	"math"

	"github.com/cockroachdb/errors"
)

const streamPrefix = "customer"

type CustomerEventStore struct {
	eventStore es.EventStore
}

func NewCustomerEventStore(eventStore es.EventStore) *CustomerEventStore {
	return &CustomerEventStore{eventStore: eventStore}
}

func (customer *CustomerEventStore) EventStreamFor(id values.CustomerID) (es.DomainEvents, error) {
	eventStream, err := customer.eventStore.LoadEventStream(customer.streamID(id), 0, math.MaxUint32)
	if err != nil {
		return nil, errors.Wrap(err, "customerEventStore.EventStreamFor")
	}

	if len(eventStream) == 0 {
		err := errors.New("found empty event stream")
		return nil, lib.MarkAndWrapError(err, lib.ErrNotFound, "customerEventStore.EventStreamFor")
	}

	return eventStream, nil
}

func (customer *CustomerEventStore) CreateStreamFrom(recordedEvents es.DomainEvents, id values.CustomerID, tx *sql.Tx) error {
	if err := customer.eventStore.AppendEventsToStream(customer.streamID(id), recordedEvents, tx); err != nil {
		if errors.Is(err, lib.ErrConcurrencyConflict) {
			err = errors.New("found duplicate customer")
			return lib.MarkAndWrapError(err, lib.ErrDuplicate, "customerEventStore.CreateStreamFrom")
		}

		return errors.Wrap(err, "customerEventStore.CreateStreamFrom")
	}

	return nil
}

func (customer *CustomerEventStore) Add(recordedEvents es.DomainEvents, id values.CustomerID, tx *sql.Tx) error {
	if err := customer.eventStore.AppendEventsToStream(customer.streamID(id), recordedEvents, tx); err != nil {
		return errors.Wrap(err, "customerEventStore.Add")
	}

	return nil
}

func (customer *CustomerEventStore) Delete(id values.CustomerID) error {
	if err := customer.eventStore.PurgeEventStream(customer.streamID(id)); err != nil {
		return errors.Wrap(err, "customerEventStore.Delete")
	}

	return nil
}

func (customer *CustomerEventStore) streamID(id values.CustomerID) es.StreamID {
	return es.NewStreamID(streamPrefix + "-" + id.ID())
}
