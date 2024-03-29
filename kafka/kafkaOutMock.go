package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

const (
	brokerAddress = "localhost:9094"
	topic         = "outbound"
)

func main() {
	fmt.Print("Kafka mock starting")
	logger := log.Default()
	wr := &kafka.Writer{
		Addr:                   kafka.TCP(brokerAddress),
		Topic:                  topic,
		Balancer:               nil,
		WriteTimeout:           10 * time.Second,
		Async:                  false,
		Logger:                 logger,
		ErrorLogger:            logger,
		AllowAutoTopicCreation: false,
	}
	defer wr.Close()

	driverId := "e4142f80-2d8c-4864-9b45-8f9eaf60dc1f"
	tripId := "550e8400-e29b-41d4-a716-446655440000"
	driverId = "2b9afec4-ee95-4849-a16e-dfd7d3144050"
	tripId = "071aacff-dd3f-40fd-b135-486d198a3f64"

	var (
		lat float32 = 1
		lng float32 = 1
	)
	createMessage := MakeCreateCommand(TripCommandCreate, tripId, lat, lng)
	produceCreateMessage(wr, createMessage)

	//var message TripCommand
	//message = MakeCommand(TripCommandAccept, driverId, tripId)
	//produceMessage(wr, message)
	//
	//message = MakeCommand(TripCommandStart, driverId, tripId)
	//produceMessage(wr, message)
	//
	//message = MakeCommand(TripCommandEnd, driverId, tripId)
	//produceMessage(wr, message)
	log.Println("Exited")
	log.Println(driverId)
}

func generateRequestID() string {
	id := uuid.New()
	return id.String()
}

func MakeCreateCommand(commandType CommandType, tripId string, lat float32, lng float32) CreatedTripCommand {
	reqId := generateRequestID()
	command := CreatedTripCommand{
		// real.
		//ID:              "h4110iam-v4ry-t1r4-d4nd-w4ntt0s1eep",
		RequestId:       reqId,
		Source:          Source,
		Type:            string(commandType),
		DataContentType: "application/json",
		Time:            time.Now(),
		Data: CreatedTripData{
			TripID:  tripId,
			OfferID: "h4110iam-v4ry-t1r4-d4nd-w4ntt51eep1",
			Price: TripPrice{
				Amount:   987,
				Currency: "RUB",
			},
			From: LatLngLiteral{
				Lat: lat,
				Lng: lng,
			},
			To: LatLngLiteral{
				Lat: 1.2,
				Lng: 1.2},
			Status: "DRIVER_SEARCH",
		},
	}
	return command
}

func MakeCommand(commandType CommandType, driverId string, tripId string) TripCommand {
	reqId := generateRequestID()
	command := TripCommand{
		RequestId:       reqId,
		Source:          Source,
		Type:            commandType,
		DataContentType: "application/json",
		Time:            time.Now(),
		Data: Data{
			TripID:   tripId,
			DriverId: driverId,
		},
	}
	return command
}

func produceCreateMessage(wr *kafka.Writer, message CreatedTripCommand) {

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %s", err)
		return
	}
	kafkaMessage := kafka.Message{
		Value: jsonMessage,
	}
	err = wr.WriteMessages(context.Background(), kafkaMessage)
	if err != nil {
		log.Printf("Error writing message: %s", err)
	}
}

func produceMessage(wr *kafka.Writer, message TripCommand) {

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %s", err)
		return
	}
	kafkaMessage := kafka.Message{
		Value: jsonMessage,
	}
	err = wr.WriteMessages(context.Background(), kafkaMessage)
	if err != nil {
		log.Printf("Error writing message: %s", err)
	}
}
