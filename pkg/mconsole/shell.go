package mconsole

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func RunShell(sockpath string) error {
	conn, err := openConn(sockpath)
	if err != nil {
		return fmt.Errorf("Couldn't connect to socket %s: %w", sockpath, err)
	}
	hostname, err := SendCommand(Proc("sys/kernel/hostname"), *conn)
	if err != nil {
		return fmt.Errorf("Failed to run command on socket %s: %w",
			sockpath, err)
	}
	_, err = SendCommand(Proc("sys/kernel/hostname"), *conn)
	if err != nil {
		fmt.Println("damn cant do multiple commands")
		return err
	}
	prompt := fmt.Sprintf("\n[%s]# ", hostname)
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(prompt)
		cmd, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		cmd = strings.Trim(cmd, "\n\r\x00")
		switch cmd {
		case "quit":
			return nil
		case "mconsole-version":
			fmt.Printf("uml_mconsole client version %d\n", MCONSOLE_VERSION)
		default:
			output, err := SendCommand(cmd, *conn)
			if err != nil {
				return err
			}
			fmt.Print(output)
		}
	}
}
