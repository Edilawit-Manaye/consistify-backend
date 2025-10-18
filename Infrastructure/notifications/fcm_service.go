package notifications

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go"          
	"firebase.google.com/go/messaging" 
)
type FCMService interface {
	SendNotification(ctx context.Context, token string, title, body string, data map[string]string) error
}

type fcmService struct {
	messagingClient *messaging.Client 
}
func NewFCMService(app *firebase.App) FCMService {
	client, err := app.Messaging(context.Background())
	if err != nil {
		log.Fatalf("Error getting Firebase Messaging client: %v", err)
	}
	log.Println("Firebase Messaging client initialized successfully.")
	return &fcmService{
		messagingClient: client,
	}
}


func (s *fcmService) SendNotification(ctx context.Context, token string, title, body string, data map[string]string) error {
	if token == "" {
		return fmt.Errorf("FCM token cannot be empty, skipping notification")
	}


	message := &messaging.Message{
		Token: token, 
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data: data, 

	
		Android: &messaging.AndroidConfig{
			Priority: "high",
		},
		
		APNS: &messaging.APNSConfig{
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					
					ContentAvailable: true,
				},
			},
		},
	
	}

	
	response, err := s.messagingClient.Send(ctx, message)
	if err != nil {
		
		log.Printf("Failed to send FCM message to token %s: %v", token, err)
		return fmt.Errorf("FCM send failed: %w", err)
	}


	log.Printf("Successfully sent FCM message to token %s. Message ID: %s", token, response)
	return nil
}



