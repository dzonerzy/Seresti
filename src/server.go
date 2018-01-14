package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func HTTP_500(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(fmt.Sprintf("500 - Something bad happened!\nError: %s", err)))
}

func HTTP_404(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 page not found"))
}

func GenericHandler(service SerestiService, config *SerestiConfig) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", fmt.Sprintf("%s %s %s", config.global.server_name, SERVER_POSTFIX, VERSION))
		if r.Method == service.method {
			var params = make(map[string]string)
			if r.Method == "GET" {
				for _, param := range service.input_parameters {
					fullparam := VAR_PREFIX + strings.ToUpper(param)
					params[fullparam] = EscapeShell(r.URL.Query().Get(param))
				}
			} else {
				r.ParseForm()
				for _, param := range service.input_parameters {
					fullparam := VAR_PREFIX + strings.ToUpper(param)
					params[fullparam] = EscapeShell(r.FormValue(param))
				}
			}
			LOG(config, "Serving script", service.sh)
			output, err := RunCGI(service.sh, params)
			if err != nil {
				LOG(config, "Error while executing", service.sh)
				HTTP_500(w, err)
			} else {
				w.Header().Set("Content-Type", "application/json")
				ParseCGIOutput(service, w, config, output)
			}
		} else {
			HTTP_404(w)
		}
	}
}

func StartServer(log *log.Logger, config *SerestiConfig) {
	LOGGER = log
	LOG(config, "Initializing Server on", config.global.listen)
	r := mux.NewRouter()
	fs := fmt.Sprintf("/%s", config.global.serve_files_uri)
	if config.global.index_path != "" {
		LOG(config, "Registering handler for index", config.global.index_path)
		r.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(config.global.index_path))))
	}
	r.HandleFunc("/api/v1/version",
		func(w http.ResponseWriter, r *http.Request) {
			result := make(map[string]interface{})
			result["code"] = 0
			result["result"] = make(map[string]interface{})
			results := make(map[string]string)
			results["app"] = APP_NAME
			results["version"] = VERSION
			result["result"] = results
			response, err := json.MarshalIndent(result, "", "\t")
			if err != nil {
				HTTP_500(w, err)
			} else {
				fmt.Fprintf(w, string(response))
			}
		})
	if config.global.serve_files {
		LOG(config, "Registering handler for filesystem", fs)
		r.PathPrefix(fs).Handler(http.StripPrefix(fs, http.FileServer(http.Dir(config.global.serve_files_path))))
	}
	for _, service := range config.services {
		if service.enabled && config.global.serve_api_version == service.version {
			LOG(config, fmt.Sprintf("Registering handler for /api/v%d/%s/%s", service.version, service.group, service.name))
			r.HandleFunc(fmt.Sprintf("/api/v%d/%s/%s", service.version, service.group, service.name),
				GenericHandler(service, config))
		}
	}
	LOG(config, "Starting server")
	http.ListenAndServe(config.global.listen, r)
}
