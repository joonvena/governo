package server

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"governo/internal/auditlog"
)

type PubSubPushMessage struct {
	Message struct {
		Data string `json:"data"`
	} `json:"message"`
}

type Server struct {
	dispatcher *auditlog.Dispatcher
}

func New(dispatcher *auditlog.Dispatcher) *Server {
	return &Server{dispatcher: dispatcher}
}

func (s *Server) pubsubPushHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("failed to close request body: %v", err)
		}
	}()
	body, _ := io.ReadAll(r.Body)

	var msg PubSubPushMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	data, err := base64.StdEncoding.DecodeString(msg.Message.Data)
	if err != nil {
		http.Error(w, "decode error", http.StatusBadRequest)
		return
	}

	log.Println(string(data))

	if err := s.dispatcher.Dispatch(r.Context(), data); err != nil {
		fmt.Println("dispatch error:", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (s *Server) Run() {
	http.HandleFunc("/pubsub/push", s.pubsubPushHandler)
	http.HandleFunc("/healthz", s.healthzHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Listening on :%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}
