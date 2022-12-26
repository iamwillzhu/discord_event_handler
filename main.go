package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bwmarrin/discordgo"
)

func Handler(request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	log.Printf("request: %v", request)

	discordPublicKey := os.Getenv("DISCORD_PUBLIC_KEY")

	if discordPublicKey == "" {
		log.Println("[Handler] discord public key is not found")
		return &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal Error: something went wrong",
		}, nil
	}

	key, err := hex.DecodeString(discordPublicKey)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Internal Error: %s", err.Error()),
		}, nil
	}

	if ok := VerifyInteraction(request, key); !ok {
		return &events.APIGatewayProxyResponse{
			StatusCode: 401,
		}, nil
	}

	interaction, err := getInteraction(request)

	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf("Internal Error: %s", err.Error()),
		}, nil
	}

	interactionResponse, err := handleInteraction(interaction)

	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf("Internal Error: %s", err.Error()),
		}, nil
	}

	interactionResponseStr, _ := json.Marshal(interactionResponse)

	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(interactionResponseStr),
	}, nil
}

func getInteraction(request *events.APIGatewayProxyRequest) (*discordgo.Interaction, error) {
	interaction := &discordgo.Interaction{}

	if err := interaction.UnmarshalJSON([]byte(request.Body)); err != nil {
		log.Printf("[getIteraction] UnmarshalJSON error: %v", err)
		return nil, err
	}

	return interaction, nil
}

func handleInteraction(interaction *discordgo.Interaction) (*discordgo.InteractionResponse, error) {
	if interaction.Type == discordgo.InteractionPing {
		return &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponsePong,
		}, nil
	}

	if interaction.Type == discordgo.InteractionApplicationCommand {
		applicationCommandData := interaction.ApplicationCommandData()

		if applicationCommandData.Name == "foo" {
			return &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "bar",
				},
			}, nil
		}

		return &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "invalid slash command",
			},
		}, nil

	}

	return nil, fmt.Errorf("interaction type not supported")
}

func main() {
	lambda.Start(Handler)
}
