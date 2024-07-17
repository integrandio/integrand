package web

import (
	"encoding/json"
	"fmt"
	"integrand/services"
	"log/slog"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

var (
	topicsAllApi    = regexp.MustCompile(`^\/api/v1/topic[\/]*$`)
	topicSingleApi  = regexp.MustCompile(`^\/api/v1/topic\/(.*)$`)
	topicEventsApi  = regexp.MustCompile(`^\/api/v1/topic\/(.*)\/events$`)
	topicConsumeApi = regexp.MustCompile(`^\/api/v1/topic\/(.*)\/consume$`)
	//topicProduceApi = regexp.MustCompile(`^\/api/v1/topic\/(.*)\/produce$`)
)

type topicAPI struct {
	userID int
}

func (ea *topicAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userId, err := apiBrowserAPIAuthenticate(w, r)
	if err != nil {
		slog.Error(err.Error())
		notFoundApiError(w)
		return
	}
	ea.userID = userId
	switch {
	//Order matters....(Since this regex is a subset of the other)
	case r.Method == http.MethodGet && topicConsumeApi.MatchString(r.URL.Path):
		ea.streamEvents(w, r)
	default:
		ea.eventRestHandler(w, r)
	}
}

func (ea *topicAPI) eventRestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	switch {
	//Glue API Routes
	case r.Method == http.MethodGet && topicEventsApi.MatchString(r.URL.Path):
		ea.getEvents(w, r)
	case r.Method == http.MethodGet && topicsAllApi.MatchString(r.URL.Path):
		ea.getTopics(w, r)
		return
	case r.Method == http.MethodGet && topicSingleApi.MatchString(r.URL.Path):
		ea.getTopic(w, r)
		return
	case r.Method == http.MethodPost && topicsAllApi.MatchString(r.URL.Path):
		ea.createTopic(w, r)
	case r.Method == http.MethodDelete && topicSingleApi.MatchString(r.URL.Path):
		ea.deleteTopic(w, r)
		return
	default:
		notFoundApiError(w)
		return
	}
}

func (ta *topicAPI) getTopics(w http.ResponseWriter, _ *http.Request) {
	eventStreams, err := services.GetEventStreams(ta.userID)
	if err != nil {
		slog.Error(err.Error())
		internalServerError(w)
		return
	}
	resJsonBytes, err := generateSuccessMessage(eventStreams)
	if err != nil {
		internalServerError(w)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resJsonBytes)
}

func (ta *topicAPI) getTopic(w http.ResponseWriter, r *http.Request) {
	matches := topicSingleApi.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		notFoundApiError(w)
		return
	}
	eventStream, err := services.GetEventStream(matches[1])
	if err != nil {
		slog.Error(err.Error())
		internalServerError(w)
		return
	}
	resJsonBytes, err := generateSuccessMessage(eventStream)
	if err != nil {
		internalServerError(w)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resJsonBytes)
}

type CreateTopicBody struct {
	TopicName string `json:"topicName"`
}

func (ta *topicAPI) createTopic(w http.ResponseWriter, r *http.Request) {
	var createBody CreateTopicBody
	if err := json.NewDecoder(r.Body).Decode(&createBody); err != nil {
		slog.Error(err.Error())
		internalServerError(w)
		return
	}
	topic, err := services.CreateEventStream(createBody.TopicName)
	if err != nil {
		slog.Error(err.Error())
		internalServerError(w)
		return
	}
	w.WriteHeader(http.StatusOK)
	resJsonBytes, _ := generateSuccessMessage(topic)
	w.Write(resJsonBytes)
}

func (ta *topicAPI) deleteTopic(w http.ResponseWriter, r *http.Request) {
	matches := topicSingleApi.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		notFoundApiError(w)
		return
	}
	topicName := matches[1]
	err := services.DeleteEventStream(topicName, ta.userID)
	if err != nil {
		slog.Error(err.Error())
		internalServerError(w)
		return
	}
	w.WriteHeader(http.StatusOK)
	c := map[string]interface{}{"success": "successfully deleted topic"}
	resJsonBytes, _ := generateSuccessMessage(c)
	w.Write(resJsonBytes)
}

func (ta *topicAPI) getEvents(w http.ResponseWriter, r *http.Request) {
	offsetParam := r.URL.Query().Get("offset")
	limitParam := r.URL.Query().Get("limit")
	offset := 0
	if offsetParam != "" {
		slog.Info("offset provided....")
		var err error
		offset, err = strconv.Atoi(offsetParam)
		if err != nil {
			// TODO: Determine what to do if this is wrong...
			slog.Error(err.Error())
			offset = 0
		}
	}

	limit := 5
	if limitParam != "" {
		slog.Info("limit provided....")
		var err error
		limit, err = strconv.Atoi(limitParam)
		if err != nil {
			// TODO: Determine what to do if this is wrong...
			slog.Error(err.Error())
			limit = 5
		}
	}

	matches := topicEventsApi.FindStringSubmatch(r.URL.Path)
	eventStream, err := services.GetEventStream(matches[1])
	if err != nil {
		slog.Error(err.Error())
		internalServerError(w)
		return
	}

	var dataArray []interface{}
	// TODO: Error handle this loop. RN if anything fails, we will get success with null data
	for i := 0; i < limit; i++ {
		//Iterate and loop through this...
		jsonBytes, err := services.GetEvent(eventStream.TopicName, offset)
		if err != nil {
			//TODO: Fix this
			if err.Error() == "offset out of bounds" {
				break
			} else {
				slog.Error(err.Error())
				break
			}
		}
		var i interface{}
		err = json.Unmarshal(jsonBytes, &i)
		if err != nil {
			slog.Error(err.Error())
			break
		}
		//We need to convert the bytes to string
		dataArray = append(dataArray, i)
		offset++
	}

	w.WriteHeader(http.StatusOK)
	resJsonBytes, _ := generateSuccessMessage(dataArray)
	w.Write(resJsonBytes)
}

// TODO: I think something is messed up here. Observed a bug when I was testing and close the client.
func (ta *topicAPI) streamEvents(w http.ResponseWriter, r *http.Request) {
	acceptType := r.Header.Get("Accept")
	if acceptType != "text/event-stream" {
		slog.Error("invalid accept type")
		internalServerError(w)
		return
	}
	offset := 0
	offsetParam := r.URL.Query().Get("offset")
	if offsetParam != "" {
		slog.Info("offset provided....")
		var err error
		offset, err = strconv.Atoi(offsetParam)
		if err != nil {
			// TODO: Determine what to do if this is wrong...
			slog.Error(err.Error())
			offset = 0
		}
	}
	matches := topicConsumeApi.FindStringSubmatch(r.URL.Path)

	eventStream, err := services.GetEventStream(matches[1])
	if err != nil {
		slog.Error(err.Error())
		internalServerError(w)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		slog.Error("SSE is not avaliable....")
		internalServerError(w)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	for {
		//Iterate and loop through this...
		jsonBytes, err := services.GetEvent(eventStream.TopicName, offset)
		if err != nil {
			//TODO: Fix this
			if err.Error() == "offset out of bounds" {
				slog.Warn(err.Error())
				time.Sleep(2 * time.Second)
				continue
			} else {
				slog.Error(err.Error())
				break
			}
		}
		fmt.Fprintf(w, "data: %s\n\n", string(jsonBytes))
		flusher.Flush()
		offset++
	}
	fmt.Fprintf(w, "data: %s\n\n", "close")
	flusher.Flush()
}
