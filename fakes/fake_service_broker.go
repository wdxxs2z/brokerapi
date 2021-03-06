package fakes

import "github.com/pivotal-cf/brokerapi"

type FakeServiceBroker struct {
	ProvisionDetails   brokerapi.ProvisionDetails
	UpdateDetails      brokerapi.UpdateDetails
	DeprovisionDetails brokerapi.DeprovisionDetails

	ProvisionedInstanceIDs   []string
	DeprovisionedInstanceIDs []string
	UpdatedInstanceIDs       []string

	BoundInstanceIDs    []string
	BoundBindingIDs     []string
	BoundBindingDetails brokerapi.BindDetails
	SyslogDrainURL      string
	RouteServiceURL     string
	VolumeMounts        []brokerapi.VolumeMount

	UnbindingDetails brokerapi.UnbindDetails

	InstanceLimit int

	ProvisionError     error
	BindError          error
	DeprovisionError   error
	LastOperationError error
	UpdateError        error

	BrokerCalled             bool
	LastOperationState       brokerapi.LastOperationState
	LastOperationDescription string

	AsyncAllowed bool

	ShouldReturnAsync     bool
	DashboardURL          string
	OperationDataToReturn string

	LastOperationInstanceID string
	LastOperationData       string
}

type FakeAsyncServiceBroker struct {
	FakeServiceBroker
	ShouldProvisionAsync bool
}

type FakeAsyncOnlyServiceBroker struct {
	FakeServiceBroker
}

func (fakeBroker *FakeServiceBroker) Services() []brokerapi.Service {
	fakeBroker.BrokerCalled = true

	return []brokerapi.Service{
		brokerapi.Service{
			ID:            "0A789746-596F-4CEA-BFAC-A0795DA056E3",
			Name:          "p-cassandra",
			Description:   "Cassandra service for application development and testing",
			Bindable:      true,
			PlanUpdatable: true,
			Plans: []brokerapi.ServicePlan{
				brokerapi.ServicePlan{
					ID:          "ABE176EE-F69F-4A96-80CE-142595CC24E3",
					Name:        "default",
					Description: "The default Cassandra plan",
					Metadata: &brokerapi.ServicePlanMetadata{
						Bullets:     []string{},
						DisplayName: "Cassandra",
					},
				},
			},
			Metadata: &brokerapi.ServiceMetadata{
				DisplayName:      "Cassandra",
				LongDescription:  "Long description",
				DocumentationUrl: "http://thedocs.com",
				SupportUrl:       "http://helpme.no",
			},
			Tags: []string{
				"pivotal",
				"cassandra",
			},
		},
	}
}

func (fakeBroker *FakeServiceBroker) Provision(instanceID string, details brokerapi.ProvisionDetails, asyncAllowed bool) (brokerapi.ProvisionedServiceSpec, error) {
	fakeBroker.BrokerCalled = true

	if fakeBroker.ProvisionError != nil {
		return brokerapi.ProvisionedServiceSpec{}, fakeBroker.ProvisionError
	}

	if len(fakeBroker.ProvisionedInstanceIDs) >= fakeBroker.InstanceLimit {
		return brokerapi.ProvisionedServiceSpec{}, brokerapi.ErrInstanceLimitMet
	}

	if sliceContains(instanceID, fakeBroker.ProvisionedInstanceIDs) {
		return brokerapi.ProvisionedServiceSpec{}, brokerapi.ErrInstanceAlreadyExists
	}

	fakeBroker.ProvisionDetails = details
	fakeBroker.ProvisionedInstanceIDs = append(fakeBroker.ProvisionedInstanceIDs, instanceID)
	return brokerapi.ProvisionedServiceSpec{DashboardURL: fakeBroker.DashboardURL}, nil
}

func (fakeBroker *FakeAsyncServiceBroker) Provision(instanceID string, details brokerapi.ProvisionDetails, asyncAllowed bool) (brokerapi.ProvisionedServiceSpec, error) {
	fakeBroker.BrokerCalled = true

	if fakeBroker.ProvisionError != nil {
		return brokerapi.ProvisionedServiceSpec{}, fakeBroker.ProvisionError
	}

	if len(fakeBroker.ProvisionedInstanceIDs) >= fakeBroker.InstanceLimit {
		return brokerapi.ProvisionedServiceSpec{}, brokerapi.ErrInstanceLimitMet
	}

	if sliceContains(instanceID, fakeBroker.ProvisionedInstanceIDs) {
		return brokerapi.ProvisionedServiceSpec{}, brokerapi.ErrInstanceAlreadyExists
	}

	fakeBroker.ProvisionDetails = details
	fakeBroker.ProvisionedInstanceIDs = append(fakeBroker.ProvisionedInstanceIDs, instanceID)
	return brokerapi.ProvisionedServiceSpec{IsAsync: fakeBroker.ShouldProvisionAsync, DashboardURL: fakeBroker.DashboardURL, OperationData: fakeBroker.OperationDataToReturn}, nil
}

func (fakeBroker *FakeAsyncOnlyServiceBroker) Provision(instanceID string, details brokerapi.ProvisionDetails, asyncAllowed bool) (brokerapi.ProvisionedServiceSpec, error) {
	fakeBroker.BrokerCalled = true

	if fakeBroker.ProvisionError != nil {
		return brokerapi.ProvisionedServiceSpec{}, fakeBroker.ProvisionError
	}

	if len(fakeBroker.ProvisionedInstanceIDs) >= fakeBroker.InstanceLimit {
		return brokerapi.ProvisionedServiceSpec{}, brokerapi.ErrInstanceLimitMet
	}

	if sliceContains(instanceID, fakeBroker.ProvisionedInstanceIDs) {
		return brokerapi.ProvisionedServiceSpec{}, brokerapi.ErrInstanceAlreadyExists
	}

	if !asyncAllowed {
		return brokerapi.ProvisionedServiceSpec{}, brokerapi.ErrAsyncRequired
	}

	fakeBroker.ProvisionDetails = details
	fakeBroker.ProvisionedInstanceIDs = append(fakeBroker.ProvisionedInstanceIDs, instanceID)
	return brokerapi.ProvisionedServiceSpec{IsAsync: true, DashboardURL: fakeBroker.DashboardURL}, nil
}

func (fakeBroker *FakeServiceBroker) Update(instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {
	fakeBroker.BrokerCalled = true

	if fakeBroker.UpdateError != nil {
		return brokerapi.UpdateServiceSpec{}, fakeBroker.UpdateError
	}

	fakeBroker.UpdateDetails = details
	fakeBroker.UpdatedInstanceIDs = append(fakeBroker.UpdatedInstanceIDs, instanceID)
	fakeBroker.AsyncAllowed = asyncAllowed
	return brokerapi.UpdateServiceSpec{IsAsync: fakeBroker.ShouldReturnAsync, OperationData: fakeBroker.OperationDataToReturn}, nil
}

func (fakeBroker *FakeServiceBroker) Deprovision(instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {
	fakeBroker.BrokerCalled = true

	if fakeBroker.DeprovisionError != nil {
		return brokerapi.DeprovisionServiceSpec{}, fakeBroker.DeprovisionError
	}

	fakeBroker.DeprovisionDetails = details
	fakeBroker.DeprovisionedInstanceIDs = append(fakeBroker.DeprovisionedInstanceIDs, instanceID)

	if sliceContains(instanceID, fakeBroker.ProvisionedInstanceIDs) {
		return brokerapi.DeprovisionServiceSpec{}, nil
	}
	return brokerapi.DeprovisionServiceSpec{IsAsync: false}, brokerapi.ErrInstanceDoesNotExist
}

func (fakeBroker *FakeAsyncOnlyServiceBroker) Deprovision(instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {
	fakeBroker.BrokerCalled = true

	if fakeBroker.DeprovisionError != nil {
		return brokerapi.DeprovisionServiceSpec{IsAsync: true}, fakeBroker.DeprovisionError
	}

	if !asyncAllowed {
		return brokerapi.DeprovisionServiceSpec{IsAsync: true}, brokerapi.ErrAsyncRequired
	}

	fakeBroker.DeprovisionedInstanceIDs = append(fakeBroker.DeprovisionedInstanceIDs, instanceID)
	fakeBroker.DeprovisionDetails = details

	if sliceContains(instanceID, fakeBroker.ProvisionedInstanceIDs) {
		return brokerapi.DeprovisionServiceSpec{IsAsync: true, OperationData: fakeBroker.OperationDataToReturn}, nil
	}

	return brokerapi.DeprovisionServiceSpec{IsAsync: true, OperationData: fakeBroker.OperationDataToReturn}, brokerapi.ErrInstanceDoesNotExist
}

func (fakeBroker *FakeAsyncServiceBroker) Deprovision(instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {
	fakeBroker.BrokerCalled = true

	if fakeBroker.DeprovisionError != nil {
		return brokerapi.DeprovisionServiceSpec{IsAsync: asyncAllowed}, fakeBroker.DeprovisionError
	}

	fakeBroker.DeprovisionedInstanceIDs = append(fakeBroker.DeprovisionedInstanceIDs, instanceID)
	fakeBroker.DeprovisionDetails = details

	if sliceContains(instanceID, fakeBroker.ProvisionedInstanceIDs) {
		return brokerapi.DeprovisionServiceSpec{IsAsync: asyncAllowed, OperationData: fakeBroker.OperationDataToReturn}, nil
	}

	return brokerapi.DeprovisionServiceSpec{OperationData: fakeBroker.OperationDataToReturn, IsAsync: asyncAllowed}, brokerapi.ErrInstanceDoesNotExist
}

func (fakeBroker *FakeServiceBroker) Bind(instanceID, bindingID string, details brokerapi.BindDetails) (brokerapi.Binding, error) {
	fakeBroker.BrokerCalled = true

	if fakeBroker.BindError != nil {
		return brokerapi.Binding{}, fakeBroker.BindError
	}

	fakeBroker.BoundBindingDetails = details

	fakeBroker.BoundInstanceIDs = append(fakeBroker.BoundInstanceIDs, instanceID)
	fakeBroker.BoundBindingIDs = append(fakeBroker.BoundBindingIDs, bindingID)

	return brokerapi.Binding{
		Credentials: FakeCredentials{
			Host:     "127.0.0.1",
			Port:     3000,
			Username: "batman",
			Password: "robin",
		},
		SyslogDrainURL:  fakeBroker.SyslogDrainURL,
		RouteServiceURL: fakeBroker.RouteServiceURL,
		VolumeMounts:    fakeBroker.VolumeMounts,
	}, nil
}

func (fakeBroker *FakeServiceBroker) Unbind(instanceID, bindingID string, details brokerapi.UnbindDetails) error {
	fakeBroker.BrokerCalled = true

	fakeBroker.UnbindingDetails = details

	if sliceContains(instanceID, fakeBroker.ProvisionedInstanceIDs) {
		if sliceContains(bindingID, fakeBroker.BoundBindingIDs) {
			return nil
		}
		return brokerapi.ErrBindingDoesNotExist
	}

	return brokerapi.ErrInstanceDoesNotExist
}

func (fakeBroker *FakeServiceBroker) LastOperation(instanceID, operationData string) (brokerapi.LastOperation, error) {
	fakeBroker.LastOperationInstanceID = instanceID
	fakeBroker.LastOperationData = operationData

	if fakeBroker.LastOperationError != nil {
		return brokerapi.LastOperation{}, fakeBroker.LastOperationError
	}

	return brokerapi.LastOperation{State: fakeBroker.LastOperationState, Description: fakeBroker.LastOperationDescription}, nil
}

type FakeCredentials struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func sliceContains(needle string, haystack []string) bool {
	for _, element := range haystack {
		if element == needle {
			return true
		}
	}
	return false
}
