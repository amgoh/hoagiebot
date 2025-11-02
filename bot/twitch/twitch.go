package twitch

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// Twitch-related environment vars
var WebhookSecret string
var TwitchClientID string
var TwitchClientSecret string

var TwitchEventChan = make(chan TwitchEventPayload, 10)

// TODO: Create Types.go to hold all necessary types
type TwitchAppAccessToken struct {
	Token string `json:"access_token"`
	ExpiryTime int `json:"expires_in"`
}

type Condition struct {
	BroadcasterUserID string `json:"broadcaster_user_id"`
	UserID string `json:"user_id"`
}

type Transport struct {
	Method string `json:"method"`
	Callback string `json:"callback"`
	Secret string `json:"secret"`
}

type Subscription struct {
	ID string `json:"id"`
	Status string `json:"status"`
	Type string `json:"type"`
	Version string `json:"version"`
	Cost int `json:"cost"`
	Condition Condition `json:"condition"`
	Transport Transport	`json:"transport"`
	Timestamp string `json:"created_at"`
}

type Event struct {
	ID string `json:"id"`
	UserID string `json:"user_id"`
	UserLogin string `json:"user_login"`
	UserName string `json:"user_name"`
	BroadcasterUserID string `json:"broadcaster_user_id"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	BroadcasterUserName string `json:"broadcaster_user_name"`
	Type string `json:"type"`
	StartTime string `json:"started_at"`
}

// Twitch EventSub JSON request body
type TwitchEventPayload struct {
	Challenge string `json:"challenge"`
	Subscription Subscription `json:"subscription"`
 	Event Event `json:"event"`
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
	if err != nil {
		return
	}
	fmt.Println("request body read!")

	if !verifyMessage(req, body) {
		fmt.Println("Signatures don't match...")
		res.WriteHeader(403)
		return
	}
	fmt.Println("Signatures match...")

	reqBody := &TwitchEventPayload{}
	err = json.Unmarshal(body, reqBody)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("request body struct filled")

	msgType := req.Header.Get("Twitch-Eventsub-Message-Type")

	// Prepare HTTP response
	switch msgType {
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
	
			TwitchEventChan <- *reqBody
		case "revocation":
			fmt.Println("revocation")
			res.WriteHeader(http.StatusOK)
		default:
			fmt.Println("ok!")
			res.WriteHeader(http.StatusOK)
		}
}

/*
	Sends a POST request to the Twitch oauth2 token API to 
	retrieve app access token 
*/
func getAppAccessToken() (*TwitchAppAccessToken, error) {
	data := url.Values{}
	data.Set("client_id", TwitchClientID)
	data.Set("client_secret", TwitchClientSecret)
	data.Set("grant_type", "client_credentials")
	tokenReqBody := data.Encode()

	tokenRequest, err := http.NewRequest(http.MethodPost, "https://id.twitch.tv/oauth2/token", 
	strings.NewReader(tokenReqBody));
	if err != nil {
		return nil, err
	}	
	tokenRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(tokenRequest)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Got response with Status: %s\n", res.Status)

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	res.Body.Close()

	accessToken := &TwitchAppAccessToken{}
	err = json.Unmarshal(resBody, accessToken)
	if err != nil {
		return nil, err
	}

	return accessToken, nil
}

func subscribeToEventSub(accessToken *TwitchAppAccessToken) (error) {
	// TODO: GuildSettings
	broadcaster_user_id := getBroadcasterID("amoghiehoagie", accessToken)

	// TODO: Switch from ngrok to duckdns -> AWS domain
	subReqBody, err := json.Marshal(&Subscription{
		Type: "stream.online",
		Version: "1",
		Condition: Condition{
			BroadcasterUserID: broadcaster_user_id+"1133192559",
		},
		Transport: Transport{
			Method: "webhook",
			Callback: "https://hoagiebot.duckdns.org/eventsub",
			Secret: WebhookSecret,
		},
	})
	
	subRequest, err := http.NewRequest(http.MethodPost, "https://api.twitch.tv/helix/eventsub/subscriptions", bytes.NewReader(subReqBody));
	if err != nil {
		return err
	}
	subRequest.Header.Add("Content-Type", "application/json")
	subRequest.Header.Add("Authorization", "Bearer " + accessToken.Token)
	subRequest.Header.Add("Client-Id", TwitchClientID)

	res, err := http.DefaultClient.Do(subRequest)
	if err != nil {
		return err
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	res.Body.Close()

	fmt.Printf("Response body: %s\n", resBody)

	return nil
}

func getBroadcasterID(username string, accessToken *TwitchAppAccessToken) (string) {
	req, err := http.NewRequest(http.MethodGet, "https://api.twitch.tv/helix/users?login="+username, nil)
	if err != nil {
		fmt.Println("Broadcaster ID GET Request failed to create...")
		return ""
	}
	req.Header.Add("Authorization", "Bearer " + accessToken.Token)
	req.Header.Add("Client-Id", TwitchClientID)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Broadcaster ID GET Request failed to send...")
		return ""
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Broadcaster ID Response failed...")
		return ""
	}
	res.Body.Close()

	fmt.Printf("Response body: %s\n", resBody)
	
	// TODO: extract broadcaster ID from resBody, need helix api json struct

	return ""
}

func StartTwitchServer() {
	if (len(WebhookSecret) == 0 || len(TwitchClientID) == 0 || len(TwitchClientSecret) == 0) {
		fmt.Println("Twitch integration environment vars not set... Aborting")
		return
	}
	
	
	accessToken, err := getAppAccessToken()
	if err != nil {
		log.Fatal(err)
		return
	}

	err = subscribeToEventSub(accessToken)
	if err != nil {
		log.Fatal(err)
		return
	}

	http.HandleFunc("/eventsub", handler)
	fmt.Println("Listening on port 8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
		return 
	}
}
