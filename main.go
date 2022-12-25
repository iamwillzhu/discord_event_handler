package main

import (
	"encoding/json"

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

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	interactionResponse := InteractionResponse{
		Type: InteractionResponsePong,
	}

	interactionResponseStr, err := json.Marshal(interactionResponse)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 401,
		}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(interactionResponseStr),
	}, nil
}

func main() {
	lambda.Start(Handler)
}
