package conf

const (
	PATH = "/usr/local/gocode/src/golang.bricks/"
)

func ServerFile() (s string) {
	return PATH + "data.txt"
}

func ClientFile() (s string) {
	return PATH + "bricks-client.txt"
}

func Empty() (e string) {
	return string([]byte{byte(0)})
}

func Finish() (e string) {
	return string([]byte{byte(1)})
}
