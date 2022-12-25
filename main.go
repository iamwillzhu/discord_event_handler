package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/iamwillzhu/discord_event_handler/interaction"
)

func Handler(request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	log.Printf("request: %v", request)

	discordPublicKey := os.Getenv("DISCORD_PUBLIC_KEY")

	if discordPublicKey == "" {
		return &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal Error: discord public key is not found",
		}, nil
	}

	key, err := hex.DecodeString(discordPublicKey)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Internal Error: %s", err.Error()),
		}, nil
	}

	if ok := interaction.VerifyInteraction(request, key); !ok {
		return &events.APIGatewayProxyResponse{
			StatusCode: 401,
		}, nil
	}

	interactionResponse := interaction.InteractionResponse{
		Type: interaction.InteractionResponsePong,
	}

	interactionResponseStr, _ := json.Marshal(interactionResponse)

	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(interactionResponseStr),
	}, nil
}

func main() {
	lambda.Start(Handler)
}
