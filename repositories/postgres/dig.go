package postgres

import "go.uber.org/dig"

// Inject api repositories
func Inject(container *dig.Container) error {
	_ = container.Provide(NewSqlxTransaction)
	_ = container.Provide(NewPgsqlUserRepository)
	_ = container.Provide(NewPgsqlOrgRepository)
	_ = container.Provide(NewPgsqlUserOrgRepository)
	// Auto generate
	// DO NOT DELETE THIS LINE ABOVE
	return nil
}
