package twitch

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Twitch-related environment vars
var WebhookSecret string
var TwitchClientID string
var TwitchClientSecret string

// Twitch EventSub JSON request body
type RequestBody struct {
	Challenge string `json:"challenge"`
	Subscription struct {
		ID string `json:"id"`
		Status string `json:"status"`
		Type string `json:"type"`
		Version string `json:"version"`
		Cost int `json:"cost"`
		Condition struct {
			BroadcasterID string `json:"broadcaster_user_id"`
		} `json:"condition"`
		Transport struct {
			Method string `json:"webhook"`
			Callback string `json:"callback"`
		}	`json:"transport"`
		Timestamp string `json:"created_at"`
	} `json:"subscription"`
	Event struct {
		ID string `json:"id"`
		UserID string `json:"user_id"`
		UserLogin string `json:"user_login"`
		UserName string `json:"user_name"`
		BroadcasterUserID string `json:"broadcaster_user_id"`
		BroadcasterUserLogin string `json:"broadcaster_user_login"`
		BroadcasterUserName string `json:"broadcaster_user_name"`
		Type string `json:"type"`
		StartTime string `json:"started_at"`
	} `json:"event"`
}

// Verify the request using Twitch's given signature before responding
func verifyMessage(req *http.Request, body []byte) bool {
	messageId := req.Header.Get("Twitch-Eventsub-Message-Id")
	timestamp := req.Header.Get("Twitch-Eventsub-Message-Timestamp")
	signature := req.Header.Get("Twitch-Eventsub-Message-Signature")

	message := messageId + timestamp + string(body)
	hash := hmac.New(sha256.New, []byte(WebhookSecret))
	hash.Write([]byte(message))

	givenSignature := "sha256=" + hex.EncodeToString(hash.Sum(nil))

	return hmac.Equal([]byte(givenSignature), []byte(signature))
}

// Handle request from Twitch
func handler(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Handle request")
	
	body, err := io.ReadAll(req.Body) 
	fmt.Println("request body read!")
	
	if err != nil {
		return
	}

	if verifyMessage(req, body) {
		fmt.Println("Signatures match...")

		reqBody := &RequestBody{}
		err = json.Unmarshal(body, reqBody)
		if err != nil {
			log.Fatal(err)
			return
		}

		fmt.Println("request body struct filled")

		messageType := req.Header.Get("Twitch-Eventsub-Message-Type")

		// Prepare HTTP response
		switch messageType {
		case "webhook_callback_verification":
			fmt.Println("verification")
			res.Header().Set("Content-Type", "text/plain")
			n, err := res.Write([]byte(reqBody.Challenge))
			
			if err != nil {
				fmt.Printf("only %d bytes written... failed to respond\n", n)
				return
			}

		case "notification":
			fmt.Println("notification")
			res.WriteHeader(http.StatusOK)
		case "revocation":
			fmt.Println("revocation")
			res.WriteHeader(http.StatusOK)
		default:
			fmt.Println("ok!")
			res.WriteHeader(http.StatusOK)
		}
	
	} else {
		fmt.Println("not ok!")
		res.WriteHeader(403)
	}
}

// Subscribe to Twitch EventSub through Webhook
func SubscribeAndListen() {
	if (len(WebhookSecret) == 0 || len(TwitchClientID) == 0 || len(TwitchClientSecret) == 0) {
		fmt.Println("Twitch integration environment vars not set... Aborting")
		return
	}

	fmt.Println("Listening...")
	http.HandleFunc("/eventsub", handler)
	fmt.Println("Set handler function...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err);
		return 
	}
}
