package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

//const defaultPort = "8080"

type jsBody struct {
	Name string `json:"key" validate:"required"`
	Age  int    `json:"age" validate:"required"`
}

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	if len(os.Args) == 2 {
		// helpCmd := flag.String("help", "", "provide http sub-commands")
		helpPtr := flag.Bool("help", false, "help")
		flag.Parse()

		if *helpPtr {
			fmt.Println("default command is http GET")
			fmt.Println("you can specify GET PUSH PUT DELETE for other purposes")
		} else {
			resp, e := http.Get("https://" + os.Args[1])
			if e != nil {
				log.Fatal(e)
			} else {
				fmt.Println(resp)
			}
		}
	} else {
		getCmd := flag.NewFlagSet("get", flag.ExitOnError)
		getHelp := getCmd.Bool("help", false, "help")
		postCmd := flag.NewFlagSet("post", flag.ExitOnError)
		postHelp := postCmd.Bool("help", false, "help")
		postJSON := postCmd.Bool("json", false, "json")
		deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
		deleteHelp := deleteCmd.Bool("help", false, "help")
		putCmd := flag.NewFlagSet("put", flag.ExitOnError)
		putHelp := putCmd.Bool("help", false, "help")

		switch os.Args[1] {
		case "get":
			getCmd.Parse(os.Args[2:])
			if *getHelp {
				fmt.Println("syntax: httpcli get <URL> [FLAGS...]")
			} else if len(os.Args) > 3 {
				var queryFlags arrayFlags
				var headerFlags arrayFlags
				getCmd.Var(&queryFlags, "query", "query")
				getCmd.Parse(os.Args[3:])
				url := "https://" + os.Args[2] + "?"
				for val := range queryFlags {
					url += queryFlags[val] + "&"
				}

				fmt.Println("Request Sent to ", url[:len(url)-1])
				client := http.Client{}
				req, _ := http.NewRequest("GET", url[:len(url)-1], nil)

				getCmd.Var(&headerFlags, "header", "header")
				getCmd.Parse(os.Args[3:])

				for val := range headerFlags {
					req.Header.Set(strings.Split(headerFlags[val], "=")[0], strings.Split(headerFlags[val], "=")[1])
				}

				resp, e := client.Do(req)
				if e != nil {
					log.Fatal(e)
				} else {
					fmt.Println(resp)
				}
			} else {
				resp, e := http.Get("https://" + os.Args[2])
				if e != nil {
					log.Fatal(e)
				} else {
					fmt.Println(resp)
				}
			}
		case "post":
			postCmd.Parse(os.Args[2:])

			if *postHelp {
				fmt.Println("syntax: httpcli post <URL> [FLAGS...]")
				return
			}

			postCmd.Parse(os.Args[3:])
			if *postJSON {
				s := os.Args[4]

				var js map[string]interface{}

				err := json.Unmarshal([]byte(s), &js)
				fmt.Println(s)
				fmt.Println([]byte(s))
				fmt.Println(js)
				if err != nil {
					fmt.Println("Wrong JSON format")
					return
				}

				var jsonBody = []byte(s)

				resp, e := http.NewRequest("POST", "https://"+os.Args[2], bytes.NewBuffer(jsonBody))
				if e != nil {
					log.Fatal(e)
				} else {
					fmt.Println("POST command with JSON format established")
					fmt.Println(resp)
				}
			} else {
				resp, e := http.NewRequest("POST", "https://"+os.Args[2], strings.NewReader(os.Args[3]))
				if e != nil {
					log.Fatal(e)
				} else {
					fmt.Println("POST command established")
					fmt.Println(resp)
				}
			}
		case "delete":
			deleteCmd.Parse(os.Args[2:])
			if *deleteHelp {
				fmt.Println("syntax: httpcli delete <URL> [FLAGS...]")
				return
			} else {
				resp, e := http.NewRequest("DELETE", "https://"+os.Args[2], strings.NewReader(os.Args[3]))
				if e != nil {
					log.Fatal(e)
				} else {
					fmt.Println("DELETE command established")
					fmt.Println(resp)
				}
			}
		case "put":
			putCmd.Parse(os.Args[2:])
			if *putHelp {
				fmt.Println("syntax: httpcli put <URL> [FLAGS...]")
				return
			}

			postCmd.Parse(os.Args[3:])

			if *postJSON {
				s := os.Args[4]

				var js jsBody

				err := json.Unmarshal([]byte(s), &js)
				if err != nil {
					fmt.Println("Wrong JSON format")
					return
				}

				var jsonBody = []byte(s)

				resp, e := http.NewRequest("PUT", "https://"+os.Args[2], bytes.NewBuffer(jsonBody))
				if e != nil {
					log.Fatal(e)
				} else {
					fmt.Println("PUT command with JSON established")
					fmt.Println(resp)
				}
			} else {
				resp, e := http.NewRequest("PUT", "https://"+os.Args[2], strings.NewReader(os.Args[3]))
				if e != nil {
					log.Fatal(e)
				} else {
					fmt.Println("PUT command established")
					fmt.Println(resp)
				}
			}
		default:
			fmt.Println("Command not found")
		}

	}
}
