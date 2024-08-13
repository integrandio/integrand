package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"integrand/persistence"
	"log"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"
)

const SLEEP_TIME int = 1
const MULTIPLYER int = 2
const MAX_BACKOFF int = 10

func init() {
	// Register all of our functions
	persistence.FUNC_MAP = map[string]interface{}{
		"ld_ld_sync":    ld_ld_sync,
		"calendly_sync": calendly_sync,
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
	currentOffset := workflow.Offset
	for {
		bytes, err := persistence.BROKER.ConsumeMessage(workflow.TopicName, currentOffset)
		if err != nil {
			if err.Error() == "offset out of bounds" {
				// This error is returned when we're given an offset thats ahead of the commitlog, so we can return for next cycle to begin
				slog.Debug(err.Error())
				break
			} else if err.Error() == "offset does not exist" {
				// This error is returned when we look for an offset and it does not exist becuase it can't be avaliable in the commitlog
				slog.Warn(err.Error())
				break // Exit the function, to be re-checked in the next cycle
			} else {
				slog.Error(err.Error())
				break // Something's wrong
			}
		}
		workflow.Call(bytes, workflow.SinkURL)
		currentOffset++
	}
	// We set offset here in case we create a new workflow with lots of messages in topic which would send redundant requests to update the offset
	if currentOffset != workflow.Offset {
		_, err := persistence.DATASTORE.SetOffsetOfWorkflow(workflow.Id, currentOffset)
		if err != nil {
			// This is a critical error. If we cannot set workflow's offset properly, our workflows will be out of sync forever
			slog.Error(err.Error())
		}
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

type CalendlyEventBody struct {
	Event   string  `json:"event"`
	Payload Payload `json:"payload"`
}

type Payload struct {
	FirstName      *string `json:"first_name"`
	LastName       *string `json:"last_name"`
	Name           string  `json:"name"`
	Email          string  `json:"email"`
	ScheduledEvent struct {
		Location struct {
			Location *string `json:"location"`
			Type     *string `json:"type"`
		} `json:"location"`
	} `json:"scheduled_event"`
}

type LeadDocketRequestBody struct {
	First   string `json:"First"`
	Last    string `json:"Last"`
	Phone   string `json:"Phone"`
	Email   string `json:"Email"`
	Summary string `json:"Summary"`
}

func calendly_sync(bytes []byte, sinkURL string) error {
	// Unmarshal the JSON byte array into the map
	var calendlyJson CalendlyEventBody
	err := json.Unmarshal(bytes, &calendlyJson)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	if calendlyJson.Event == "invitee.created" {
		err := sendCalendlyAppointment(calendlyJson, sinkURL)
		if err != nil {
			slog.Error("Error occurred while sending calendly appointment to CLF", "error", err)
			return err
		}
	}

	return nil
}

func sendCalendlyAppointment(calendlyJson CalendlyEventBody, sinkURL string) error {
	payload := calendlyJson.Payload
	var request LeadDocketRequestBody

	if payload.FirstName != nil {
		request.First = *payload.FirstName
		request.Last = *payload.LastName
	} else {
		nameParts := strings.Split(strings.TrimSpace(payload.Name), " ")
		request.First = nameParts[0]
		if len(nameParts) > 1 {
			request.Last = nameParts[1]
		} else {
			request.Last = ""
		}
	}

	if *payload.ScheduledEvent.Location.Type == "outbound_call" {
		request.Phone = *payload.ScheduledEvent.Location.Location
	}

	request.Email = payload.Email

	jsonBodyBytes, err := json.Marshal(request)
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
