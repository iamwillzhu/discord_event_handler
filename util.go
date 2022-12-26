package main

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"io"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

// VerifyInteraction implements message verification of the discord interactions api
// signing algorithm, as documented here:
// https://discord.com/developers/docs/interactions/receiving-and-responding#security-and-authorization
func VerifyInteraction(r *events.APIGatewayProxyRequest, key ed25519.PublicKey) bool {
	var msg bytes.Buffer

	signature, ok := r.Headers["x-signature-ed25519"]
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

	timestamp, ok := r.Headers["x-signature-timestamp"]
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
