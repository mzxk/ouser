package ouserClient

import (
	"bufio"
	"fmt"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

func ReadPwd(s string) (result string) {
	fmt.Println(s)
	pwd, err := terminal.ReadPassword(0)
	if err != nil {
		return err.Error()
	}
	return string(pwd)
}
func Readline(s string) string {
	fmt.Println(s)
	reader := bufio.NewReader(os.Stdin)
	data, _, _ := reader.ReadLine()
	return string(data)
}
