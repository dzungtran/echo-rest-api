package usecases

import "go.uber.org/dig"

// Inject api handlers
func Inject(container *dig.Container) error {
	_ = container.Provide(NewUserUsecase)
	_ = container.Provide(NewOrgUsecase)
	// Auto generate
	// DO NOT DELETE THIS LINE ABOVE
	return nil
}
