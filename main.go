package main

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// InteractionResponseType is type of interaction response.
type InteractionResponseType uint8

// Interaction response types.
const (
	// InteractionResponsePong is for ACK ping event.
	InteractionResponsePong InteractionResponseType = 1
	// InteractionResponseChannelMessageWithSource is for responding with a message, showing the user's input.
	InteractionResponseChannelMessageWithSource InteractionResponseType = 4
	// InteractionResponseDeferredChannelMessageWithSource acknowledges that the event was received, and that a follow-up will come later.
	InteractionResponseDeferredChannelMessageWithSource InteractionResponseType = 5
	// InteractionResponseDeferredMessageUpdate acknowledges that the message component interaction event was received, and message will be updated later.
	InteractionResponseDeferredMessageUpdate InteractionResponseType = 6
	// InteractionResponseUpdateMessage is for updating the message to which message component was attached.
	InteractionResponseUpdateMessage InteractionResponseType = 7
	// InteractionApplicationCommandAutocompleteResult shows autocompletion results. Autocomplete interaction only.
	InteractionApplicationCommandAutocompleteResult InteractionResponseType = 8
	// InteractionResponseModal is for responding to an interaction with a modal window.
	InteractionResponseModal InteractionResponseType = 9
)

// InteractionResponse represents a response for an interaction event.
type InteractionResponse struct {
	Type InteractionResponseType `json:"type,omitempty"`
}

// VerifyInteraction implements message verification of the discord interactions api
// signing algorithm, as documented here:
// https://discord.com/developers/docs/interactions/receiving-and-responding#security-and-authorization
func VerifyInteraction(r *events.APIGatewayProxyRequest, key ed25519.PublicKey) bool {
	var msg bytes.Buffer

	signature, ok := r.Headers["X-Signature-Ed25519"]
	if !ok || signature == "" {
		return false
	}

	sig, err := hex.DecodeString(signature)
	if err != nil {
		return false
	}

	if len(sig) != ed25519.SignatureSize {
		return false
	}

	timestamp, ok := r.Headers["X-Signature-Timestamp"]
	if !ok || timestamp == "" {
		return false
	}

	msg.WriteString(timestamp)

	var body bytes.Buffer
	bodyReader := strings.NewReader(r.Body)

	// copy body into buffers
	_, err = io.Copy(&msg, io.TeeReader(bodyReader, &body))
	if err != nil {
		return false
	}

	return ed25519.Verify(key, msg.Bytes(), sig)
}

func Handler(request *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {

	log.Printf("request: %v", request)

	discordPublicKey := os.Getenv("DISCORD_PUBLIC_KEY")

	if discordPublicKey == "" {
		return &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal Error: discord public key is not found",
		}
	}

	key, err := hex.DecodeString(discordPublicKey)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Internal Error: %s", err.Error()),
		}
	}

	if ok := VerifyInteraction(request, key); !ok {
		return &events.APIGatewayProxyResponse{
			StatusCode: 401,
		}
	}

	interactionResponse := InteractionResponse{
		Type: InteractionResponsePong,
	}

	interactionResponseStr, _ := json.Marshal(interactionResponse)

	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(interactionResponseStr),
	}
}

func main() {
	lambda.Start(Handler)
}
