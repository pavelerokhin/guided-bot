package main

import (
	"bufio"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"regexp"
	"strings"
)

type States struct {
	States []State `json:"states"`
}

func (s *States) GetState(id int64) *State {
	for _, state := range s.States {
		if state.ID == id {
			return &state
		}
	}

	return nil
}

type State struct {
	ID     int64  `yaml:"id"`
	Before string `yaml:"before"`
	Text   string `yaml:"text"`
	Input  string `yaml:"input"`
	After  string `yaml:"after"`
	Next   *Next  `yaml:"next"`
}

type Next struct {
	RightId int64  `yaml:"right"`
	RightIf string `yaml:"right-if"`
	LeftId  int64  `yaml:"left"`
}

func (t *Next) IsSimple() bool {
	return t.RightIf == "" && t.LeftId == 0
}

type memory map[string]string
type functions map[string]func(string) error
type filters map[string]func(string) bool

var (
	ff = functions{
		"print": func(s string) error {
			if s == "" {
				fmt.Println(s)
			}

			return nil
		},
	}
	fl = filters{
		"isEmpty": func(s string) bool {
			return s == ""
		},
	}
	m = make(memory)
)

func main() {
	// reade conversation.yml file and parse it to States struct
	var states States

	f, err := os.Open("./conversation.yml")
	yamlFile, err := io.ReadAll(f)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = yaml.Unmarshal(yamlFile, &states)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var i, j int64
	for {
		if i == 999 {
			break
		}

		state := states.GetState(i)
		if state == nil {
			fmt.Println("*** NO STATE ***", i)
			break
		}

		if state.Before != "" {
			// get variable name from the string state.Before, using pattern {(var_name)}
			key := extractVariableName(state.Before)
			functionName := extractFunctionName(state.Before)
			function := ff[functionName]
			if function != nil {
				err := function(m[key])
				if err != nil {
					fmt.Println("ERROR")
					fmt.Println(err)
					os.Exit(1)
				}
			}
		}
		key := extractVariableName(state.Text)
		if key != "" {
			state.Text = strings.ReplaceAll(state.Text, "{"+key+"}", m[key])
		}
		fmt.Println(state.Text)
		if state.Next.IsSimple() {
			j = state.Next.RightId
			fmt.Println("simple, next id", j)
		} else {
			key := state.Input
			value := readInputFromUser()
			fmt.Println("k ey", key, "value", value)
			m[key] = value

			fmt.Println("input", state.Input)
			fmt.Println("right id", state.Next.RightId)
			fmt.Println("left id", state.Next.LeftId)
			j = state.Next.RightId
		}
		if state.After != "" {
			fmt.Println(state.After)
			f := extractFunctionName(state.After)
			v := extractVariableName(state.After)
			fmt.Println("AFTER function ", f, "variable", v)
		}
		i = j
	}

	fmt.Println("end")
}

func extractVariableName(input string) string {
	// Define the regular expression pattern to match "{(var_name)}"
	re := regexp.MustCompile(`\{(.*?)\}`)

	// Find the first match in the input string
	match := re.FindStringSubmatch(input)

	// Extract and return the variable name from the match
	if len(match) > 1 {
		return match[1]
	}

	return "" // Return an empty string if no variable name found
}

func extractFunctionName(input string) string {
	// Define the regular expression pattern to match "(.*)\("
	re := regexp.MustCompile(`(.*)\(`)

	// Find the first match in the input string
	match := re.FindStringSubmatch(input)

	// Extract and return the variable name from the match
	if len(match) > 1 {
		return match[1]
	}

	return "" // Return an empty string if no variable name found
}

func readInputFromUser() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter a string and press Enter:")
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
