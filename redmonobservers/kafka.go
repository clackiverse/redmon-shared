package redmonobservers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Shopify/sarama"
	uuid "github.com/satori/go.uuid"
	"github.com/turnage/graw/reddit"
)

func ToKafkaTopic(topic string) func(context.Context, interface{}) (interface{}, error) {
	config := sarama.NewConfig()

	config.ClientID = fmt.Sprintf("%s", uuid.NewV1())
	config.Version = sarama.V2_4_0_0
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(strings.Split("127.0.0.1:9092", ","), config)

	if err != nil {
		fmt.Println("cannot create producer", err)
	}

	return func(_ context.Context, o interface{}) (interface{}, error) {
		post, err := o.(reddit.Post)

		if err == false {
			fmt.Println("message can't be cast")
		}
		binary, _ := json.Marshal(post)

		producerMessage := sarama.ProducerMessage{
			Topic: "incoming_reddit",
			Key:   sarama.StringEncoder(post.ID),
			Value: sarama.StringEncoder(string(binary)),
		}

		producer.SendMessage(&producerMessage)

		return o, nil
	}

}
