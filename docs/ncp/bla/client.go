package bla

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log/slog"
// 	"time"

// 	"github.com/nats-io/nats.go"
// )

// func SchemaPublish() error {
// 	data := PublishSchemaMsg{
// 		Name:    "test",
// 		Version: "1.0.0",
// 		Schema:  EventTxtar,
// 	}
// 	bData, err := json.Marshal(data)
// 	if err != nil {
// 		return fmt.Errorf("marshal: %w", err)
// 	}

// 	return clientSendMsg(bData, "actor.schema.publish")
// }

// func clientSendMsg(data []byte, subject string) error {
// 	nc, err := nats.Connect(nats.DefaultURL)
// 	if err != nil {
// 		return fmt.Errorf("connect: %w", err)
// 	}
// 	defer nc.Close()

// 	reply, err := nc.Request(subject, data, time.Minute)
// 	if err != nil {
// 		return fmt.Errorf("request: %w", err)
// 	}

// 	slog.Info("reply", "data", string(reply.Data))
// 	return nil
// }

// func clientSendMsg(data []byte, subject string) error {
// 	nc, err := nats.Connect(nats.DefaultURL)
// 	if err != nil {
// 		return fmt.Errorf("connect: %w", err)
// 	}
// 	defer nc.Close()
// 	js, err := jetstream.New(nc)
// 	if err != nil {
// 		return fmt.Errorf("jetstream: %w", err)
// 	}

// 	// Setup reply inbox and channel for receiving messages
// 	replyTo := nats.NewInbox()
// 	ch := make(chan *nats.Msg, 64)
// 	sub, err := nc.ChanSubscribe(replyTo, ch)
// 	if err != nil {
// 		return fmt.Errorf("subscribe: %w", err)
// 	}
// 	defer sub.Unsubscribe()

// 	// Create and publish message
// 	msg := nats.NewMsg(subject)
// 	msg.Data = data
// 	msg.Header.Add("Reply-To", replyTo)
// 	if _, err := js.PublishMsg(context.TODO(), msg); err != nil {
// 		return fmt.Errorf("publish: %w", err)
// 	}
// 	// fmt.Println("ack: ", ack)
// 	timec := time.After(time.Minute * 10)
// 	for {
// 		select {
// 		case <-timec:
// 			fmt.Println("CLIENT TIMEOUT")
// 			return nil
// 		case m := <-ch:
// 			fmt.Println("CLIENT MESSAGE: ", string(m.Data))
// 			var event ActionEvent
// 			if err := json.Unmarshal(m.Data, &event); err != nil {
// 				return fmt.Errorf("unmarshal: %w", err)
// 			}
// 			switch event.Type {
// 			case EventTypeInfo:
// 				slog.Info(event.Message)
// 			case EventTypeError:
// 				slog.Error(event.Message)
// 				return fmt.Errorf("error: %s", event.Message)
// 			case EventTypeCompleted:
// 				slog.Info("completed")
// 				return nil
// 			}
// 		}
// 	}
// }
