package compl

import (
	"bytes"
	"os/exec"
	"strconv"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/pkg/errors"
)

func cmd(name string, args ...string) (string, error) {
	out, err := exec.Command(name, args...).Output()
	if err != nil {
		errStr := err.(*exec.ExitError).Stderr
		return "", errors.Wrap(err, "error running command: "+
			strings.Join(append([]string{name}, args[:]...), " ")+"Error: "+string(errStr))
	}
	return string(out), nil
}

func cmdPipe(name string, args ...string) (string, error) {
	cmd := exec.Command(name, (args[:len(args)-1])...)
	buf := bytes.Buffer{}
	buf.Write([]byte(args[len(args)-1]))
	cmd.Stdin = &buf
	out, err := cmd.Output()
	if err != nil {
		errStr := err.(*exec.ExitError).Stderr
		return "", errors.Wrap(err, "error running command: "+
			strings.Join(append([]string{name}, args[:]...), " ")+"Error: "+string(errStr))
	}
	return string(out), nil
}

func shell(cmd ...string) (string, error) {
	cArgs := []string{"-c", strings.Join(cmd, "")}
	out, err := exec.Command("bash", cArgs...).Output()
	if err != nil {
		errStr := err.(*exec.ExitError).Stderr
		return "", errors.Wrap(err, "error running command: "+strings.Join(cmd, "")+"Error: "+string(errStr))
	}
	return string(out), nil
}

func getFuncMap() template.FuncMap {
	funcMap := sprig.TxtFuncMap()
	funcMap["shell"] = shell
	funcMap["cmd"] = cmd
	funcMap["cmdPipe"] = cmdPipe
	return funcMap
}

func getTemplate() *template.Template {
	return template.New("name").Funcs(getFuncMap())
}

func parseTemplate(templateStr string) (*template.Template, error) {
	tmpl, err := getTemplate().Parse(templateStr)
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

func render(tmpl string, args []string, kwargs map[string]string) (bytes.Buffer, error) {
	t, err := parseTemplate(tmpl)
	if err != nil {
		return bytes.Buffer{}, err
	}
	return renderFromTemplate(t, args, kwargs)
}

func renderFromTemplate(t *template.Template, args []string, kwargs map[string]string) (bytes.Buffer, error) {
	data := make(map[string]string)
	for k, v := range kwargs {
		data[k] = v
	}
	for i, v := range args {
		data["_"+strconv.Itoa(i+1)] = v
	}
	var buf bytes.Buffer
	err := t.Execute(&buf, data)
	if err != nil {
		return bytes.Buffer{}, err
	}
	return buf, nil
}
