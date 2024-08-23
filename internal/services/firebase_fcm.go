package services

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

type FCMClient struct {
	client *messaging.Client
}


func NewFCMClient() (*FCMClient, error) {
	// credentials := filepath.Join("../../files/config/credentilals.json")

	data, err := os.Open("../../files/config/credentilals.json")
	if err != nil {
		fmt.Println("Error al leer el archivo de credenciales: %v", err)
	}

	fmt.Printf(data)

	fmt.Println(data)
	opt := option.WithCredentialsFile(credentials)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error al inicializar la app de Firebase: %v", err)
	}

	client, err := app.Messaging(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error al inicializar el cliente de mensajería: %v", err)
	}

	return &FCMClient{client: client}, nil
}

func (s *FCMClient) SendMessage(token, title, body string) error {
	mns := &messaging.Message{
		Token: token,
		Notification: &messaging.Notification{
			Title: title,
			Body: body,
		},
	}

	response, err := s.client.Send(context.Background(), mns)
	if err != nil {
		return fmt.Errorf("error al enviar el mensaje: %v", err)
	}

	fmt.Printf("Mensaje enviado exitosamente: %v\n", response)
	return nil
}
