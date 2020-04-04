package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/goccy/go-yaml"
	"github.com/mnowotnik/fzshell/internal/compl"
	"github.com/mnowotnik/fzshell/internal/output"
)

const LB_TOK = "{{"
const RB_TOK = "}}"
const REP_START = "[["
const REP_STOP = "]]"
const DefaultItemSeparator = " "

type Options struct {
	MatchMultipleSpaces bool `default:"True" yaml:"matchMultipleSpaces"`
}

type Config struct {
	Completions []compl.Completion `yaml:"completions"`
	Options     Options            `yaml:"options"`
}

type lineInfo struct {
	lBuffer string
	rBuffer string
}

func parseLine(lineBuffer string, chCursorPos int) (lineInfo, error) {

	if len(lineBuffer) == 0 {
		return lineInfo{"", ""}, nil
	}

	if chCursorPos >= len(lineBuffer) {
		return lineInfo{lineBuffer, ""}, nil
	}

	return lineInfo{lineBuffer[0:chCursorPos], lineBuffer[chCursorPos:]}, nil
}

func Run() int {

	var args struct {
		LineBuffer string `arg:"required,positional" help:"Current line buffer or '-' to read from stdin."`
		CursorPos  int    `arg:"--cursor" help:"Cursor position in the line buffer" default:"-1"`
		AllItems   bool   `arg:"--all" help:"Print all results instead of running finder"`
		ConfigPath string `arg:"-c,--config" help:"Path to a configuration file"`
	}
	arg.MustParse(&args)
	lineBuffer := args.LineBuffer
	chCursorPos := args.CursorPos
	if lineBuffer == "-" {
		bytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			output.Log().Println(err)
			return 1
		}
		lineBuffer = string(bytes)
	}
	if chCursorPos == -1 {
		chCursorPos = len(lineBuffer)
	}

	configPath := ""
	if args.ConfigPath != "" {
		configPath = args.ConfigPath
	} else if os.Getenv("FZSHELL_CONFIG") != "" {
		configPath = os.Getenv("FZSHELL_CONFIG")
	} else {
		configHome := os.Getenv("XDG_CONFIG_HOME")
		if configHome == "" {
			configHome = os.Getenv("HOME") + "/.config"
		}
		configPath = configHome + "/fzshell/fzshell.yaml"
	}
	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		output.Log().Printf("Could not read config file: %s\n", configPath)
	}
	config := Config{}
	if err := yaml.Unmarshal(yamlFile, &config); err != nil {
		output.Log().Printf("Could not unmarshal: %v", err)
		return 1
	}
	lineInfo, err := parseLine(lineBuffer, chCursorPos)
	if err != nil {
		output.Log().Println(err)
		return 1
	}
	for _, completion := range config.Completions {
		options := compl.CompletionOptions{
			MatchMultipleSpaces: config.Options.MatchMultipleSpaces,
			ReturnAll:           args.AllItems,
		}
		result, code := completion.MatchAndFind(lineInfo.lBuffer, options)
		switch code {
		case compl.Matched:
			if args.AllItems {
				for _, r := range result {
					fmt.Println(r)
				}
			} else {
				fmt.Print(result[0])
			}
			return 0
		case compl.Errored:
			return 1
		case compl.Aborted:
			return 0
		case compl.NotMatched:
			continue
		}
	}
	return 0
}
