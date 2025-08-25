package twitch

import (
	"crypto"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

var secret string = "this is the secret for the webhooks"


var MESSAGE_TYPE string = "Twitch-Eventsub-Message-Type"
var MESSAGE_TYPE_VERIFICATION string = "webhook_callback_verification"
var MESSAGE_TYPE_NOTIFICATION string = "notification"
var MESSAGE_TYPE_REVOCATION string = "revocation"


//app.post('/eventsub', (req, res) => {
//    let secret = getSecret();
//    let message = getHmacMessage(req);
//    let hmac = HMAC_PREFIX + getHmac(secret, message);  // Signature to compare
//
//    if (true === verifyMessage(hmac, req.headers[TWITCH_MESSAGE_SIGNATURE])) {
//        console.log("signatures match");
//
//        // Get JSON object from body, so you can process the message.
//        let notification = JSON.parse(req.body);
//        
//        if (MESSAGE_TYPE_NOTIFICATION === req.headers[MESSAGE_TYPE]) {
//            // TODO: Do something with the event's data.
//
//            console.log(`Event type: ${notification.subscription.type}`);
//            console.log(JSON.stringify(notification.event, null, 4));
//            
//            res.sendStatus(204);
//        }
//        else if (MESSAGE_TYPE_VERIFICATION === req.headers[MESSAGE_TYPE]) {
//            res.set('Content-Type', 'text/plain').status(200).send(notification.challenge);
//        }
//        else if (MESSAGE_TYPE_REVOCATION === req.headers[MESSAGE_TYPE]) {
//            res.sendStatus(204);
//
//            console.log(`${notification.subscription.type} notifications revoked!`);
//            console.log(`reason: ${notification.subscription.status}`);
//            console.log(`condition: ${JSON.stringify(notification.subscription.condition, null, 4)}`);
//        }
//        else {
//            res.sendStatus(204);
//            console.log(`Unknown message type: ${req.headers[MESSAGE_TYPE]}`);
//        }
//    }
//    else {
//        console.log('403');    // Signatures didn't match.
//        res.sendStatus(403);
//    }
//})
//  
//app.listen(port, () => {
//  console.log(`Example app listening at http://localhost:${port}`);
//})
//

func verifyMessage(req *http.Request, body []byte) bool {
	messageId := req.Header.Get("Twitch-Eventsub-Message-Id")
	timestamp := req.Header.Get("Twitch-Eventsub-Message-Timestamp")
	signature := req.Header.Get("Twitch-Eventsub-Message-Signature")

	
	message := messageId + timestamp + string(body)
	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write([]byte(message))

	givenSignature := "sha256=" + hex.EncodeToString(hash.Sum(nil))

	return hmac.Equal([]byte(givenSignature), []byte(signature))
}

func callback(req *http.Request, res *http.Response) {
	body, err := io.ReadAll(req.Body) 
	if err != nil {
		return
	}


	if verifyMessage(req, body) {
		fmt.Println("Signatures match...")

		notification, err := io.ReadAll(req.Body)
		if err != nil {
			return
		}

		fmt.Println(notification)

		messageType := req.Header.Get("Twitch-Eventsub-Message-Type")

		switch messageType {
		case "webhook_callback_verification":
			fmt.Println("callback")
		case "notification":
			fmt.Println("notification")
		case "revocation":
			fmt.Println("revocation")
		default:
			fmt.Println("ok!")
		}
	
	} else {
		fmt.Println("not ok!")
	}
}
