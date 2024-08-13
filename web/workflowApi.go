package web

import (
	"encoding/json"
	"integrand/services"
	"log/slog"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
)

var (
	workflowFunctionApi = regexp.MustCompile(`^\/api/v1/workflow/functions[\/]*$`)
	workflowsAllApi     = regexp.MustCompile(`^\/api/v1/workflow[\/]*$`)
	workflowSingleApi   = regexp.MustCompile(`^\/api/v1/workflow\/(.*)$`)
)

type workflowAPI struct {
}

func (wf *workflowAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := apiBrowserAPIAuthenticate(w, r)
	if err != nil {
		slog.Error(err.Error())
		apiMessageResponse(w, http.StatusUnauthorized, "Authentication needed")
		return
	}
	switch {
	case r.Method == http.MethodGet && workflowFunctionApi.MatchString(r.URL.Path):
		wf.getWorkflowFunctions(w, r)
	case r.Method == http.MethodGet && workflowsAllApi.MatchString(r.URL.Path):
		wf.getWorkflows(w, r)
		return
	case r.Method == http.MethodGet && workflowSingleApi.MatchString(r.URL.Path):
		wf.getWorkflow(w, r)
		return
	case r.Method == http.MethodPut && workflowSingleApi.MatchString(r.URL.Path):
		wf.updateWorkflow(w, r)
		return
	case r.Method == http.MethodDelete && workflowSingleApi.MatchString(r.URL.Path):
		wf.deleteWorkflow(w, r)
		return
	case r.Method == http.MethodPost && workflowsAllApi.MatchString(r.URL.Path):
		wf.createWorkflow(w, r)
		return
	default:
		apiMessageResponse(w, http.StatusNotFound, "not found")
		return
	}
}

func (wf *workflowAPI) getWorkflows(w http.ResponseWriter, _ *http.Request) {
	workflows, err := services.GetWorkflows()
	if err != nil {
		slog.Error(err.Error())
		apiMessageResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}
	resJsonBytes, err := generateSuccessMessage(workflows)
	if err != nil {
		apiMessageResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resJsonBytes)
}

func (wf *workflowAPI) getWorkflow(w http.ResponseWriter, r *http.Request) {
	matches := workflowSingleApi.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		apiMessageResponse(w, http.StatusBadRequest, "incorrect request sent")
		return
	}
	id, err := strconv.Atoi(matches[1])
	if err != nil {
		apiMessageResponse(w, http.StatusBadRequest, "incorrect request sent")
		return
	}
	workflow, err := services.GetWorkflow(id)
	if err != nil {
		slog.Error(err.Error())
		apiMessageResponse(w, http.StatusBadRequest, "incorrect request sent")
		return
	}
	resJsonBytes, err := generateSuccessMessage(workflow)
	if err != nil {
		apiMessageResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resJsonBytes)
}

type CreateWorkflowBody struct {
	TopicName    string `json:"topicName"`
	FunctionName string `json:"functionName"`
	SinkURL      string `json:"sinkURL"`
}

func (wf *workflowAPI) createWorkflow(w http.ResponseWriter, r *http.Request) {
	var createBody CreateWorkflowBody
	if err := json.NewDecoder(r.Body).Decode(&createBody); err != nil {
		slog.Error(err.Error())
		apiMessageResponse(w, http.StatusBadRequest, "invalid body sent")
		return
	}
	_, err := url.ParseRequestURI(createBody.SinkURL)
	if err != nil {
		slog.Error("invalid url")
		apiMessageResponse(w, http.StatusBadRequest, "invalid sink url sent")
		return
	}
	workflow, err := services.CreateWorkflow(createBody.TopicName, createBody.FunctionName, createBody.SinkURL)
	if err != nil {
		slog.Error(err.Error())
		apiMessageResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}
	resJsonBytes, err := generateSuccessMessage(workflow)
	if err != nil {
		apiMessageResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resJsonBytes)
}

func (wf *workflowAPI) updateWorkflow(w http.ResponseWriter, r *http.Request) {
	matches := workflowSingleApi.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		apiMessageResponse(w, http.StatusBadRequest, "incorrect request sent")
		return
	}
	id, err := strconv.Atoi(matches[1])
	if err != nil {
		apiMessageResponse(w, http.StatusBadRequest, "incorrect request sent")
		return
	}
	workflow, err := services.UpdateWorkflow(id)
	if err != nil {
		slog.Error(err.Error())
		apiMessageResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}
	resJsonBytes, err := generateSuccessMessage(workflow)
	if err != nil {
		apiMessageResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resJsonBytes)
}

func (wf *workflowAPI) deleteWorkflow(w http.ResponseWriter, r *http.Request) {
	matches := workflowSingleApi.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		apiMessageResponse(w, http.StatusBadRequest, "incorrect request sent")
		return
	}
	id, err := strconv.Atoi(matches[1])
	if err != nil {
		apiMessageResponse(w, http.StatusBadRequest, "incorrect request sent")
		return
	}
	_, err = services.DeleteWorkflow(id)
	if err != nil {
		slog.Error(err.Error())
		apiMessageResponse(w, http.StatusBadRequest, "internal server error")
		return
	}
	c := map[string]interface{}{"success": "successfully deleted workflow"}
	resJsonBytes, err := generateSuccessMessage(c)
	if err != nil {
		apiMessageResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resJsonBytes)
}

func (wf *workflowAPI) getWorkflowFunctions(w http.ResponseWriter, _ *http.Request) {
	workflowFunctions := services.GetAvaliableWorkflowFunctions()
	resJsonBytes, err := generateSuccessMessage(workflowFunctions)
	if err != nil {
		apiMessageResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resJsonBytes)
}
