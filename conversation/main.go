package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"OpenAI-api/api"
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
		"imageKitchen": func(s string) error {
			s = fmt.Sprintf("make a photorealistic peacture of kitchen, well detailed, sotisfying the following description: {%s}", s)
			resp, err := api.HandleImageCreateL(s)
			if err != nil {
				return err
			}

			// print in blue color
			fmt.Printf("\x1b[94m%s\033[0m\n", "Here is the kitchen's sketch: kitchen.png")
			m["imageKitchen"] = "kitchen.png"

			return SaveBase64Image(resp[0], "kitchen.png")
		},
		"editKitchen": func(s string) error {
			s = fmt.Sprintf("change the  kitchen image following the suggestions: {%s}", s)
			resp, err := api.HandleImageEditL(m["imageKitchen"], s)
			if err != nil {
				return err
			}

			// print in blue color
			fmt.Printf("\x1b[94m%s\033[0m\n", "Here is the kitchen's editing: kitchen.png")
			m["imageKitchen"] = "kitchen.png"

			return SaveBase64Image(resp[0], "kitchen.png")
		},
		"imageLivingRoom": func(s string) error {
			s = fmt.Sprintf("make a photorealistic peacture of the living room, well detailed, sotisfying the following description: {%s}", s)
			resp, err := api.HandleImageCreateL(s)
			if err != nil {
				return err
			}

			// print in blue color
			fmt.Printf("\x1b[94m%s\033[0m\n", "Here is the living room's sketch: living_room.png")
			m["imageLivingRoom"] = "living_room.png"

			return SaveBase64Image(resp[0], "living_room.png")
		},
		"editLivingRoom": func(s string) error {
			s = fmt.Sprintf("change the  kitchen image following the suggestions: {%s}", s)
			resp, err := api.HandleImageEditL(m["imageLivingRoom"], s)
			if err != nil {
				return err
			}

			// print in blue color
			fmt.Printf("\x1b[94m%s\033[0m\n", "Here is the living room's sketch: living_room.png")
			m["imageLivingRoom"] = "living_room.png"

			return SaveBase64Image(resp[0], "living_room.png")
		},
		"imageTerrace": func(s string) error {
			s = fmt.Sprintf("make a photorealistic peacture of the terrace, well detailed, sotisfying the following description: {%s}", s)
			resp, err := api.HandleImageCreateL(s)
			if err != nil {
				return err
			}

			// print in blue color
			fmt.Printf("\x1b[94m%s\033[0m\n", "Here is the terrace's sketch: terrace.png")
			m["imageTerrace"] = "terrace.png"

			return SaveBase64Image(resp[0], "terrace.png")
		},
		"editTerrace": func(s string) error {
			s = fmt.Sprintf("change the  kitchen image following the suggestions: {%s}", s)
			resp, err := api.HandleImageEditL(m["imageTerrace"], s)
			if err != nil {
				return err
			}

			// print in blue color
			fmt.Printf("\x1b[94m%s\033[0m\n", "Here is the terrace's sketch: terrace.png")
			m["imageTerrace"] = "terrace.png"

			return SaveBase64Image(resp[0], "terrace.png")
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
	viper.SetConfigFile("config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("failed to read config file: %s", err))
	}

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
			break
		}

		// before hook
		if state.Before != "" {
			// get variable name from the string state.Before, using pattern {(var_name)}
			key := extractVariableName(state.Before)
			functionName := extractFunctionName(state.Before)
			function := ff[functionName]
			if function != nil {
				err := function(m[key])
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}
		}

		// message
		key := extractVariableName(state.Text)
		if key != "" {
			state.Text = strings.ReplaceAll(state.Text, "{"+key+"}", m[key])
		}
		// print message in green color
		fmt.Printf("\033[32m%s\033[0m\n", state.Text)

		if state.Input != "" {
			key := state.Input
			value := readInputFromUser()
			m[key] = value

		}
		// filter
		if state.Next.RightIf != "" {
			filterKey := extractVariableName(state.Next.RightIf)
			filterName := extractFunctionName(state.Next.RightIf)
			filter := fl[filterName]
			if filter != nil {
				if filter(m[filterKey]) {
					j = state.Next.RightId
				} else {
					j = state.Next.LeftId
				}
			} else {
				j = state.Next.RightId
			}
		} else {
			j = state.Next.RightId
		}

		if state.After != "" {
			key := extractVariableName(state.After)
			functionName := extractFunctionName(state.After)
			function := ff[functionName]
			if function != nil {
				err := function(m[key])
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}
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
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
