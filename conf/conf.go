package conf

const (
	PATH        = "/usr/local/gocode/src/golang.bricks/"
	SERVER_FILE = PATH + "data.txt"
	CLIENT_FILE = PATH + "bricks-client.txt"
)

func Empty() (e string) {
	return string([]byte{byte(0)})
}

func Finish() (e string) {
	return string([]byte{byte(1)})
}
