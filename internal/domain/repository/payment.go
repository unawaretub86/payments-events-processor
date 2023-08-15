package repository

func (repositoryPayment repositoryPayment) CreatePayment(orderId, requestId string) (*string, error) {
	return repositoryPayment.database.CreatePayment(orderId, requestId)
}

func (repositoryPayment repositoryPayment) UpdatePayment(orderId, requestId string) (*string, *string, error) {
	return repositoryPayment.database.UpdatePayment(orderId, requestId)
}
