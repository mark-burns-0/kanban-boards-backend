package auth

type UserRepository interface {
	Get(uint64)
}
