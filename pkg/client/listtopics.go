package client

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type ListTopicsResponse struct {
	topics []string
}

func ListTopics(c *kafka.Client) ListTopicsResponse {
	mresp, err := c.Metadata(context.Background(), &kafka.MetadataRequest{})
	if err != nil {
		panic(err)
	}

	ret := ListTopicsResponse{}

	ret.topics = make([]string, len(mresp.Topics))

	for i, t := range mresp.Topics {
		ret.topics[i] = t.Name
	}

	return ret
}

func (r *ListTopicsResponse) Print() {
	for _, t := range r.topics {
		fmt.Println(t)
	}
}
