package storage

type ServiceProvider struct {
}

func (database *ServiceProvider) Register() {
	defaultStorage = NewStorage("")
}

func (database *ServiceProvider) Boot() {

}
