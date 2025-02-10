package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// WebSocket Upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins (for testing purposes only!)
	},
}

// WebSocket Message Structure
type WebSocketMessage struct {
	Message string `json:"message"`
}

// Load Type Structure
type LoadType struct {
	Type string `json:"type"`
}

// WebSocket Handler
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	log.Println("WebSocket connection established")

	for {
		// Read message from the client
		_, p, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Unexpected WebSocket close error: %v", err)
			}
			break
		}

		// Parse the incoming message to determine the load type
		var loadType LoadType
		err = json.Unmarshal(p, &loadType)
		if err != nil {
			log.Printf("Error unmarshaling JSON: %v", err)
			errorMsg := WebSocketMessage{Message: "Invalid JSON format"}
			errorBytes, _ := json.Marshal(errorMsg)
			conn.WriteMessage(websocket.TextMessage, errorBytes)
			continue
		}

		// Handle the load type
		switch loadType.Type {
		case "simple":
			log.Println("Simulating simple load for WebSocket")
			time.Sleep(1 * time.Second) // Simulate processing time
			responseMessage := WebSocketMessage{Message: "Simple load test completed"}
			responseBytes, _ := json.Marshal(responseMessage)
			if err := conn.WriteMessage(websocket.TextMessage, responseBytes); err != nil {
				log.Printf("Error writing message: %v", err)
				return
			}

		case "medium":
			log.Println("Simulating medium load for WebSocket")
			for i := 0; i < 10; i++ {
				responseMessage := WebSocketMessage{
					Message: fmt.Sprintf("Medium load message %d processed", i+1),
				}
				responseBytes, _ := json.Marshal(responseMessage)
				if err := conn.WriteMessage(websocket.TextMessage, responseBytes); err != nil {
					log.Printf("Error writing message: %v", err)
					return
				}
				time.Sleep(1 * time.Second) // Simulate processing time for each message
			}
			// Send completion message
			completionMsg := WebSocketMessage{Message: "Medium load test completed"}
			completionBytes, _ := json.Marshal(completionMsg)
			if err := conn.WriteMessage(websocket.TextMessage, completionBytes); err != nil {
				log.Printf("Error writing completion message: %v", err)
				return
			}

		default:
			log.Println("Invalid load type for WebSocket")
			errorMsg := WebSocketMessage{Message: "Invalid load type"}
			errorBytes, _ := json.Marshal(errorMsg)
			if err := conn.WriteMessage(websocket.TextMessage, errorBytes); err != nil {
				log.Printf("Error writing error message: %v", err)
				return
			}
		}
	}
}

// Complex HTTP Handler (simulates resource utilization)
func handleComplexHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var msg WebSocketMessage
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Simulate resource utilization
	startTime := time.Now()
	log.Printf("Simulating resource utilization for message: %s\n", msg.Message)

	// Simulate CPU-intensive work
	result := 0
	for i := 0; i < 1000000; i++ {
		result += i * i
	}

	// Simulate a delay to mimic database interaction
	time.Sleep(2 * time.Second)

	duration := time.Since(startTime)

	response := map[string]interface{}{
		"message":         fmt.Sprintf("Message processed: %s", msg.Message),
		"computation":     result,
		"processing_time": duration.String(),
		"timestamp":       time.Now().Format(time.RFC3339),
		"status":          "success",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Load Test Handler for HTTP
func handleLoadTestHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var loadType LoadType
	if err := json.NewDecoder(r.Body).Decode(&loadType); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	startTime := time.Now()

	switch loadType.Type {
	case "simple":
		log.Println("Simulating simple load for HTTP")
		time.Sleep(1 * time.Second)
		response := map[string]interface{}{
			"message":         "Simple load test completed",
			"processing_time": time.Since(startTime).String(),
			"timestamp":       time.Now().Format(time.RFC3339),
			"status":          "success",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

	case "medium":
		log.Println("Simulating medium load for HTTP")
		messages := make([]string, 0)
		for i := 0; i < 10; i++ {
			message := fmt.Sprintf("Medium load message %d processed", i+1)
			messages = append(messages, message)
			time.Sleep(1 * time.Second)
		}
		response := map[string]interface{}{
			"message":         "Medium load test completed",
			"messages":        messages,
			"processing_time": time.Since(startTime).String(),
			"timestamp":       time.Now().Format(time.RFC3339),
			"status":          "success",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

	default:
		http.Error(w, "Invalid load type", http.StatusBadRequest)
	}
}

func main() {
	router := mux.NewRouter()

	// WebSocket Route
	router.HandleFunc("/ws", handleWebSocket)

	// HTTP Routes
	router.HandleFunc("/complex", handleComplexHTTP).Methods("POST")
	router.HandleFunc("/load-test-http", handleLoadTestHTTP).Methods("POST")

	// Start the server
	fmt.Println("Server listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
