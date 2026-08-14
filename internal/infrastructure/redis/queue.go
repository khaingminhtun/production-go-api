package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

const (
	EmailStream        = "email:queue"
	EmailConsumerGroup = "email-workers"
)

type EmailJob struct {
	ID        string         `json:"id"`
	To        string         `json:"to"`
	Subject   string         `json:"subject"`
	Template  string         `json:"template"`
	Data      map[string]any `json:"data"`
	CreatedAt time.Time      `json:"created_at"`
}

type EmailJobMessage struct {
	MessageID string
	Job       EmailJob
}

type EmailQueue interface {
	Publish(
		ctx context.Context,
		job EmailJob,
	) error

	Consume(
		ctx context.Context,
		consumerName string,
		count int,
		block time.Duration,
	) ([]EmailJobMessage, error)

	Ack(
		ctx context.Context,
		messageID string,
	) error

	EnsureConsumerGroup(
		ctx context.Context,
	) error
}

type emailQueue struct {
	client *goredis.Client
}

func NewEmailQueue(client *goredis.Client) EmailQueue {
	return &emailQueue{
		client: client,
	}
}

func (q *emailQueue) Publish(
	ctx context.Context,
	job EmailJob,
) error {

	data, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("marshal email job: %w", err)
	}

	_, err = q.client.XAdd(
		ctx,
		&goredis.XAddArgs{
			Stream: EmailStream,
			ID:     "*",
			Values: map[string]any{
				"job": string(data),
			},
		},
	).Result()

	if err != nil {
		return fmt.Errorf("publish email job: %w", err)
	}

	return nil
}

func (q *emailQueue) EnsureConsumerGroup(
	ctx context.Context,
) error {

	err := q.client.XGroupCreateMkStream(
		ctx,
		EmailStream,
		EmailConsumerGroup,
		"0",
	).Err()

	if err == nil {
		return nil
	}

	// Consumer group already exists.
	// This is expected when the application restarts.
	if strings.Contains(err.Error(), "BUSYGROUP") {
		return nil
	}

	return fmt.Errorf(
		"create email consumer group: %w",
		err,
	)
}

func (q *emailQueue) Consume(
	ctx context.Context,
	consumerName string,
	count int,
	block time.Duration,
) ([]EmailJobMessage, error) {

	result, err := q.client.XReadGroup(
		ctx,
		&goredis.XReadGroupArgs{
			Group:    EmailConsumerGroup,
			Consumer: consumerName,
			Streams: []string{
				EmailStream,
				">",
			},
			Count: int64(count),
			Block: block,
			NoAck: false,
		},
	).Result()

	if err != nil {

		if err == goredis.Nil {
			return nil, nil
		}

		return nil, fmt.Errorf(
			"consume email jobs: %w",
			err,
		)
	}

	var messages []EmailJobMessage

	for _, stream := range result {

		for _, message := range stream.Messages {

			rawJob, ok := message.Values["job"].(string)
			if !ok {
				return nil, fmt.Errorf(
					"invalid email job format",
				)
			}

			var job EmailJob

			if err := json.Unmarshal(
				[]byte(rawJob),
				&job,
			); err != nil {
				return nil, fmt.Errorf(
					"unmarshal email job: %w",
					err,
				)
			}

			messages = append(
				messages,
				EmailJobMessage{
					MessageID: message.ID,
					Job:       job,
				},
			)
		}
	}

	return messages, nil
}

func (q *emailQueue) Ack(
	ctx context.Context,
	messageID string,
) error {

	if err := q.client.XAck(
		ctx,
		EmailStream,
		EmailConsumerGroup,
		messageID,
	).Err(); err != nil {

		return fmt.Errorf(
			"ack email job %s: %w",
			messageID,
			err,
		)
	}

	return nil
}
