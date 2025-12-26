package ports

type AuthUseCase interface {
	LoginFromPassword(username, password string) (string, error)
}
