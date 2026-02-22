package bot

type opener interface {
	Open(url string) error
}
