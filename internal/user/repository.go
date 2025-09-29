package user

type UserRepository interface {
	Get(uint64)
}
