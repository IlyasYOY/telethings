package thingsurl

type thingsURL struct {
	authToken string
}

func New(authToken string) *thingsURL {
	return &thingsURL{
		authToken: authToken,
	}
}
