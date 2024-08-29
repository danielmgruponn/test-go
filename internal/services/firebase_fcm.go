package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

type FCMClient struct {
	client *messaging.Client
}

type FirebaseCredentials struct {
	ProjectID string `json:"project_id"`
}

func NewFCMClient() (*FCMClient, error) {
	credentials := "files/config/credential.json"

	data, err := os.ReadFile(credentials)
	if err != nil {
		fmt.Printf("Error al leer el archivo de credenciales: %s\n", err)
	}

	var credential FirebaseCredentials

	err = json.Unmarshal(data, &credential)
	if err != nil {
		log.Fatalf("error unmarshalling credentials: %v\n", err)
	}

	projectID := credential.ProjectID

	config := &firebase.Config{
		ProjectID: projectID,
	}

	opt := option.WithCredentialsFile(credentials)
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		return nil, fmt.Errorf("error al inicializar la app de Firebase: %v", err)
	}

	client, err := app.Messaging(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error al inicializar el cliente de mensajería: %v", err)
	}

	return &FCMClient{client: client}, nil
}

func (s *FCMClient) SendMessage(token, title, body string, data map[string]string) error {
	mns := &messaging.Message{
		Token: token,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data: data,
	}

	response, err := s.client.Send(context.Background(), mns)
	if err != nil {
		return fmt.Errorf("error al enviar el mensaje: %v", err)
	}

	fmt.Printf("Mensaje enviado exitosamente: %v\n", response)
	return nil
}
