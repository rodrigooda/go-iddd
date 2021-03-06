package cmd

import (
	"database/sql"

	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/hexagon/application"
	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/hexagon/application/domain/customer"
	customergrpc "github.com/AntonStoeckl/go-iddd/service/customeraccounts/infrastructure/adapter/grpc"
	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/infrastructure/adapter/postgres"
	"github.com/AntonStoeckl/go-iddd/service/shared"
	"github.com/AntonStoeckl/go-iddd/service/shared/es"
	"github.com/cockroachdb/errors"
)

const (
	eventStoreTableName           = "eventstore"
	uniqueEmailAddressesTableName = "unique_email_addresses"
)

type DIContainer struct {
	postgresDBConn                    *sql.DB
	customerEventStore                *postgres.CustomerEventStore
	marshalCustomerEvent              es.MarshalDomainEvent
	unmarshalCustomerEvent            es.UnmarshalDomainEvent
	buildUniqueEmailAddressAssertions customer.ForBuildingUniqueEmailAddressAssertions
	customerCommandHandler            *application.CustomerCommandHandler
	customerQueryHandler              *application.CustomerQueryHandler
	customerGRPCServer                customergrpc.CustomerServer
}

func NewDIContainer(
	postgresDBConn *sql.DB,
	marshalCustomerEvent es.MarshalDomainEvent,
	unmarshalCustomerEvent es.UnmarshalDomainEvent,
	buildUniqueEmailAddressAssertions customer.ForBuildingUniqueEmailAddressAssertions,
) (*DIContainer, error) {

	if postgresDBConn == nil {
		return nil, errors.Mark(errors.New("newDIContainer: postgres DB connection must not be nil"), shared.ErrTechnical)
	}

	container := &DIContainer{
		postgresDBConn:                    postgresDBConn,
		marshalCustomerEvent:              marshalCustomerEvent,
		unmarshalCustomerEvent:            unmarshalCustomerEvent,
		buildUniqueEmailAddressAssertions: buildUniqueEmailAddressAssertions,
	}

	container.init()

	return container, nil
}

func (container DIContainer) init() {
	container.GetCustomerEventStore()
	container.GetCustomerCommandHandler()
	container.GetCustomerQueryHandler()
	container.GetCustomerGRPCServer()
}

func (container DIContainer) GetPostgresDBConn() *sql.DB {
	return container.postgresDBConn
}

func (container DIContainer) GetCustomerEventStore() *postgres.CustomerEventStore {
	if container.customerEventStore == nil {
		container.customerEventStore = postgres.NewCustomerEventStore(
			container.postgresDBConn,
			eventStoreTableName,
			container.marshalCustomerEvent,
			container.unmarshalCustomerEvent,
			uniqueEmailAddressesTableName,
			container.buildUniqueEmailAddressAssertions,
		)
	}

	return container.customerEventStore
}

func (container DIContainer) GetCustomerCommandHandler() *application.CustomerCommandHandler {
	if container.customerCommandHandler == nil {
		container.customerCommandHandler = application.NewCustomerCommandHandler(
			container.GetCustomerEventStore().RetrieveEventStream,
			container.GetCustomerEventStore().StartEventStream,
			container.GetCustomerEventStore().AppendToEventStream,
		)
	}

	return container.customerCommandHandler
}

func (container DIContainer) GetCustomerQueryHandler() *application.CustomerQueryHandler {
	if container.customerQueryHandler == nil {
		container.customerQueryHandler = application.NewCustomerQueryHandler(
			container.GetCustomerEventStore().RetrieveEventStream,
		)
	}

	return container.customerQueryHandler
}

func (container DIContainer) GetCustomerGRPCServer() customergrpc.CustomerServer {
	if container.customerGRPCServer == nil {
		container.customerGRPCServer = customergrpc.NewCustomerServer(
			container.GetCustomerCommandHandler().RegisterCustomer,
			container.GetCustomerCommandHandler().ConfirmCustomerEmailAddress,
			container.GetCustomerCommandHandler().ChangeCustomerEmailAddress,
			container.GetCustomerCommandHandler().ChangeCustomerName,
			container.GetCustomerCommandHandler().DeleteCustomer,
			container.GetCustomerQueryHandler().CustomerViewByID,
		)
	}

	return container.customerGRPCServer
}
