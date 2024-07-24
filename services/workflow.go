package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"reflect"
)

var caseTypeMapping = map[string]int{
	"Motor Vehicle Accident (MVA)": 4,
	"Premises Liability":           15,
	"Dog Bite":                     14,
	"Other":                        2,
}

var Workflows []Workflow

type Workflow struct {
	TopicName    string
	Offset       int
	FunctionName string
	Enabled      bool
}

type funcMap map[string]interface{}

var FUNC_MAP = funcMap{}

func init() {
	// Register all of our functions
	FUNC_MAP = map[string]interface{}{
		"ld_ld_sync": ld_ld_sync,
	}
}

func (workflow Workflow) Call(params ...interface{}) (result interface{}, err error) {
	f := reflect.ValueOf(FUNC_MAP[workflow.FunctionName])
	if len(params) != f.Type().NumIn() {
		err = errors.New("the number of params is out of index")
		return
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	res := f.Call(in)
	result = res[0].Interface()
	return
}

func ld_ld_sync(bytes []byte) error {
	log.Println("Executing")

	// Unmarshal the JSON byte array into the map
	var jsonBody map[string]interface{}
	err := json.Unmarshal(bytes, &jsonBody)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	if jsonBody["LeadStatus"] == "Referred Out" &&
		jsonBody["LeadSubstatus"] == "Pending Review" &&
		jsonBody["LeadReferredTo"] == "The Capital Law Firm" {
		err := sendLeadToClf(jsonBody)
		if err != nil {
			slog.Error("Error occurred while sending lead to CLF", "error", err)
			return err
		}
	}
	return nil
}

// Should move to utils later

func GetOrDefaultString(m map[string]interface{}, key string, defaultStr string) string {
	if value, ok := m[key]; ok {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return defaultStr
}

func GetOrDefaultInt(m map[string]int, key string, defaultInt int) int {
	if num, ok := m[key]; ok {
		return num
	}
	return defaultInt
}

func sendLeadToClf(jsonBody map[string]interface{}) error {
	defaultStr := ""
	sinkUrl := os.Getenv("SINK_URL")

	leadCaseTypeStr := GetOrDefaultString(jsonBody, "LeadCaseType", "")

	requestBody := map[string]interface{}{
		"First":         GetOrDefaultString(jsonBody, "ContactFirstName", defaultStr),
		"Last":          GetOrDefaultString(jsonBody, "ContactLastName", defaultStr),
		"Phone":         GetOrDefaultString(jsonBody, "ContactMobilePhone", defaultStr),
		"Email":         GetOrDefaultString(jsonBody, "ContactEmail", defaultStr),
		"Summary":       GetOrDefaultString(jsonBody, "LeadSummary", defaultStr),
		"Case_Type":     GetOrDefaultInt(caseTypeMapping, leadCaseTypeStr, 2),
		"Incident_Date": GetOrDefaultString(jsonBody, "LeadIncidentDate", defaultStr),
	}

	jsonBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	resp, err := http.Post(sinkUrl, "application/json", bytes.NewBuffer(jsonBodyBytes))
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("HTTP request failed",
			"status_code", resp.StatusCode,
			"status_text", http.StatusText(resp.StatusCode),
		)
		return errors.New("HTTP request failed with status code: " + http.StatusText(resp.StatusCode))
	}

	var responseBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		slog.Error(err.Error())
		return err
	}

	log.Printf("Status Code: %d", resp.StatusCode)
	log.Printf("Response Body: %v", responseBody)
	return nil
}
