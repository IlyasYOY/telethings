package thingsurl

type versionThingsURL struct {
	*thingsURL
}

func (tu *thingsURL) Version() *versionThingsURL {
	return &versionThingsURL{thingsURL: tu}
}

func (vtu *versionThingsURL) String() string {
	return "things:///version"
}
