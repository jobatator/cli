package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"syscall"

	"github.com/TylerBrock/colorjson"
	prompt "github.com/c-bata/go-prompt"
	"github.com/jobatator/cli/pkg/connexion"
	"golang.org/x/crypto/ssh/terminal"
)

// LivePrefixState -
var LivePrefixState struct {
	LivePrefix string
	IsEnable   bool
}

var conn net.Conn
var reader *bufio.Reader
var suggests []prompt.Suggest
var enableRaw bool = true
var jsonFormatter *colorjson.Formatter

func read() string {
	out, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	return out
}

func handleOutput() {
	output := read()
	jsonStr := []byte("")
	if enableRaw {
		// pretty print json objects
		if output[0:2] == "{\"" && output[len(output)-2:len(output)-1] == "}" {
			// case of a json object
			var obj map[string]interface{}
			json.Unmarshal([]byte(output), &obj)
			jsonStr, _ = jsonFormatter.Marshal(obj)
		}
		if output[0:1] == "[" && output[len(output)-2:len(output)-1] == "]" {
			// case of a json array
			var obj []interface{}
			json.Unmarshal([]byte(output), &obj)
			jsonStr, _ = jsonFormatter.Marshal(obj)
		}
	}
	if len(jsonStr) == 0 {
		fmt.Print(output)
	} else {
		fmt.Println(string(jsonStr))
	}
}

func executor(in string) {
	if strings.ToUpper(in) == "QUIT" || strings.ToUpper(in) == "EXIT" {
		conn.Close()
		os.Exit(0)
	}
	fmt.Fprintf(conn, in+"\n")
	handleOutput()
}

func completer(in prompt.Document) []prompt.Suggest {
	if len(strings.Split(in.Text, " ")) == 1 && len(in.Text) != 0 {
		return prompt.FilterHasPrefix(suggests, in.GetWordBeforeCursor(), true)
	}
	return []prompt.Suggest{}
}

func changeLivePrefix() (string, bool) {
	return LivePrefixState.LivePrefix + "> ", true
}

func main() {
	var commandsToRun []string
	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	if !(fi.Mode()&os.ModeNamedPipe == 0) {
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		commandsToRun = strings.Split(string(input), ";")
	}
	var url string = ""
	for _, arg := range os.Args {
		if arg == "-r" || arg == "--raw" {
			enableRaw = false
		} else {
			url = arg
		}
	}

	if url == "" {
		url = "127.0.0.1"
	}
	options := connexion.ParseURL(url)

	if enableRaw {
		jsonFormatter = colorjson.NewFormatter()
		jsonFormatter.Indent = 4
	}

	// connect to socket
	tmpConn, err := net.Dial("tcp", options.Host+":"+options.Port)
	if err != nil {
		panic(err)
	}
	conn = tmpConn
	reader = bufio.NewReader(conn)

	if len(options.Username) > 0 && len(options.Password) == 0 {
		// ask for the password
		fmt.Print("Password for " + options.Username + "@" + options.Host + ":" + options.Port + " : ")
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			panic(err)
		}
		fmt.Println()
		options.Password = string(bytePassword)
	}

	isAuthenticated := false
	if len(options.Password) > 0 {
		conn.Write([]byte("AUTH " + options.Username + " " + options.Password + " "))
		out := read()
		out = out[:len(out)-1]
		if out != "Welcome!" {
			fmt.Println(out)
			conn.Close()
			os.Exit(0)
		} else {
			isAuthenticated = true
			LivePrefixState.LivePrefix = options.Username + "@"
		}
	}
	LivePrefixState.LivePrefix += options.Host + ":" + options.Port

	if len(options.Group) > 0 && isAuthenticated {
		// set group
		conn.Write([]byte("USE_GROUP " + options.Group + " "))
		out := read()
		out = out[:len(out)-1]
		if out != "OK" {
			fmt.Println(out)
			conn.Close()
			os.Exit(0)
		}
		LivePrefixState.LivePrefix += "/" + options.Group
	}

	// AUTOMATIC MODE
	if len(commandsToRun) > 0 {
		for key, cmd := range commandsToRun {
			if len(commandsToRun)-1 == key {
				cmd = cmd[0 : len(cmd)-1]
			}
			cmd = strings.Trim(cmd, " ") + " "
			conn.Write([]byte(cmd))
			handleOutput()
		}
		os.Exit(0)
	}

	// MANUAL MOOE
	// fetch all commands
	conn.Write([]byte("HELP "))
	out := read()
	out = out[:len(out)-1]
	var commands []string
	json.Unmarshal([]byte(out), &commands)
	for _, cmd := range commands {
		entry := prompt.Suggest{Text: cmd, Description: ""}
		suggests = append(suggests, entry)
	}

	fmt.Println("Use Ctrl+D to exit")
	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix(">>> "),
		prompt.OptionLivePrefix(changeLivePrefix),
		prompt.OptionTitle("jobatator "+conn.RemoteAddr().String()),
		prompt.OptionInputTextColor(prompt.Cyan),
		prompt.OptionPrefixTextColor(prompt.DarkGray),
	)
	p.Run()
}
