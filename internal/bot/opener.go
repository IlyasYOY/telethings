package bot

//go:generate go run github.com/gojuno/minimock/v3/cmd/minimock@v3.4.7 -i opener -o opener_mock_test.go -g

type opener interface {
	Open(url string) error
}
