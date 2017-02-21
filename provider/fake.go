package provider

type FakeProvider struct {
}

func (f FakeProvider) Scale(serviceId string, target int, direction bool) error {
	return nil
}
