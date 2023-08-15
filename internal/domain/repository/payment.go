package repository

func (repositoryPayment repositoryPayment) CreatePayment(orderId string, requestId string) (*string, error) {
	return repositoryPayment.database.CreatePayment(orderId, requestId)
}
