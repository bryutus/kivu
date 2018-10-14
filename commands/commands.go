package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"

	"github.com/manifoldco/promptui"
	"github.com/urfave/cli"
)

type property struct {
	Scripts map[string]string `json:"scripts"`
}

type lazyloadProperty struct {
	Scripts json.RawMessage `json:"scripts"`
}

type script struct {
	Alias   string
	Command string
}

var Action = func(c *cli.Context) error {
	in, err := ioutil.ReadFile("package.json")
	if err != nil {
		log.Fatalln("Failed to read package.json:", err)
		os.Exit(1)
	}

	var property property
	if err := json.Unmarshal(in, &property); err != nil {
		log.Fatalln("Failed to parse package.json:", err)
		os.Exit(1)
	}

	var lazyload lazyloadProperty
	if err := json.Unmarshal(in, &lazyload); err != nil {
		log.Fatalln("Failed to parse package.json:", err)
		os.Exit(1)
	}

	var buf bytes.Buffer
	if err := json.Indent(&buf, lazyload.Scripts, "", "  "); err != nil {
		log.Fatalln("Failed to parse package.json:", err)
		os.Exit(1)
	}
	scriptsJSON := buf.String()

	keys := extractKeys(scriptsJSON)

	selects, err := listSelects(property, keys)
	if err != nil {
		log.Fatalln("Failed to make selects:", err)
		os.Exit(1)
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "▸ {{ .Alias | yellow }} [{{ .Command | cyan }}]",
		Inactive: "  {{ .Alias }} [{{ .Command | faint }}]",
		Selected: "✔ {{ .Alias | yellow }} [{{ .Command | cyan }}]",
	}

	prompt := promptui.Select{
		Label:     "Select npm run script",
		Items:     selects,
		Templates: templates,
		Size:      15,
	}

	i, _, err := prompt.Run()

	if err != nil {
		log.Fatalln("Failed to:", err)
		os.Exit(1)
	}

	if isKeyword(selects[i].Alias) {
		out, _ := exec.Command("npm", selects[i].Alias).CombinedOutput()
		fmt.Println(string(out))
		return nil
	}

	out, _ := exec.Command("npm", "run-script", selects[i].Alias).CombinedOutput()
	fmt.Println(string(out))

	return nil
}

func extractKeys(s string) (keys []string) {
	r := regexp.MustCompile(`(?m)"(.+)":`)
	matched := r.FindAllStringSubmatch(s, -1)
	for _, v := range matched {
		keys = append(keys, v[1])
	}
	return
}

func listSelects(p property, orderKeys []string) (selects []script, err error) {
	selects = []script{}
	for _, key := range orderKeys {
		var script script
		script.Alias = key
		script.Command = p.Scripts[key]
		selects = append(selects, script)
	}
	return
}

func isKeyword(alias string) bool {
	keywords := []string{
		"start",
		"restart",
		"stop",
		"test",
		"publish",
		"install",
		"uninstall",
		"update",
	}

	for _, keyword := range keywords {
		if alias == keyword {
			return true
		}
	}
	return false
}
