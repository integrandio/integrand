package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"integrand/persistence"
	"log"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

const SLEEP_TIME int = 1
const MULTIPLYER int = 2
const MAX_BACKOFF int = 10

func init() {
	// Register all of our functions
	persistence.FUNC_MAP = map[string]interface{}{
		"ld_ld_sync": ld_ld_sync,
	}
}

func Workflower() error {
	log.Println("Workflower started")
	for {
		time.Sleep(100 * time.Millisecond)
		currentWorkflows, _ := GetEnabledWorkflows()

		var wg sync.WaitGroup
		for _, workflow := range currentWorkflows {
			wg.Add(1)
			go processWorkflow(&wg, workflow)
		}
		wg.Wait()
	}
}

func processWorkflow(wg *sync.WaitGroup, workflow persistence.Workflow) {
	defer wg.Done()
	sleep_time := SLEEP_TIME
	for {
		bytes, err := persistence.BROKER.ConsumeMessage(workflow.TopicName, workflow.Offset)
		if err != nil {
			if err.Error() == "offset out of bounds" {
				// This error is returned when we're given an offset thats ahead of the commitlog
				slog.Debug(err.Error())
				time.Sleep(time.Duration(sleep_time) * time.Second)
				continue
			} else if err.Error() == "offset does not exist" {
				// This error is returned when we look for an offset and it does not exist becuase it can't be avaliable in the commitlog
				slog.Warn(err.Error())
				time.Sleep(time.Duration(sleep_time) * time.Second)
				return // Exit the function, to be re-checked in the next cycle
			} else {
				slog.Error(err.Error())
				return // Something's wrong
			}
		}
		workflow.Call(bytes, workflow.SinkURL)
		workflow.Offset++
		sleep_time = SLEEP_TIME
	}
}

var caseTypeMapping = map[string]int{
	"Motor Vehicle Accident (MVA)": 4,
	"Premises Liability":           15,
	"Dog Bite":                     14,
	"Other":                        2,
}

func ld_ld_sync(bytes []byte, sinkURL string) error {
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
		err := sendLeadToClf(jsonBody, sinkURL)
		if err != nil {
			slog.Error("Error occurred while sending lead to CLF", "error", err)
			return err
		}
	}
	return nil
}

func sendLeadToClf(jsonBody map[string]interface{}, sinkURL string) error {
	defaultStr := ""

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

	resp, err := http.Post(sinkURL, "application/json", bytes.NewBuffer(jsonBodyBytes))
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
