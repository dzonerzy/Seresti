package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

type CGIOutput struct {
	was_error   bool
	stderr      string
	stdout      string
	return_code int
}

func RunCGI(path string, params map[string]string) (*CGIOutput, error) {
	out := new(CGIOutput)
	cmd := exec.Command("sh", path)
	cmd.Env = os.Environ()
	for k, v := range params {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
	cmd.Env = append(cmd.Env, "SERESTISEP=\t\t\t")
	cmd.Env = append(cmd.Env, fmt.Sprintf("SERESTIVER=%s", VERSION))
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				if status.ExitStatus() != 0 {
					out.was_error = true
					out.return_code = status.ExitStatus()
				}
			}
		} else {
			return nil, err
		}
	}
	out.stderr = stderr.String()
	out.stdout = stdout.String()
	return out, nil
}

func ParseCGIOutput(service SerestiService, w http.ResponseWriter, config *SerestiConfig, output *CGIOutput) {
	result := make(map[string]interface{})
	if output.was_error {
		result["code"] = output.return_code
		if config.global.verbose_error {
			result["error_message"] = output.stderr
		} else {
			result["error_message"] = ""
		}
		result["result"] = nil
		response, err := json.MarshalIndent(result, "", "\t")
		if err != nil {
			HTTP_500(w, err)
		} else {
			fmt.Fprintf(w, string(response))
		}
	} else {
		result["code"] = 0
		result["result"] = make(map[string]interface{})
		output_lines := strings.Split(output.stdout, "\n")
		if output_lines[1] == "" {
			results := make(map[string]string)
			lines_val := strings.Split(output_lines[0], "\t\t\t")
			if len(lines_val) != len(service.output_parameters) {
				HTTP_500(w, fmt.Errorf("Mismatched number of output arguments!"))
				return
			}
			for i, output_variable := range service.output_parameters {
				results[output_variable] = lines_val[i]
			}
			result["result"] = results
		} else {
			result["code"] = 0
			result["result"] = make(map[string]interface{})
			results := make([]map[string]string, 0)
			for _, line := range output_lines {
				single_result := make(map[string]string)
				if line != "" {
					lines_val := strings.Split(line, "\t\t\t")
					if len(lines_val) != len(service.output_parameters) {
						HTTP_500(w, fmt.Errorf("Mismatched number of output arguments!"))
						return
					}
					for i, output_variable := range service.output_parameters {
						single_result[output_variable] = lines_val[i]
					}
					results = append(results, single_result)
				}
			}
			result["result"] = results[0 : len(output_lines)-1]
		}
		response, err := json.MarshalIndent(result, "", "\t")
		if err != nil {
			HTTP_500(w, err)
		} else {
			fmt.Fprintf(w, string(response))
		}
	}
}
