package client

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type ListGroupsResponse struct {
	groups []string
}

func ListGroups(c *kafka.Client) ListGroupsResponse {
	lgresp, err := c.ListGroups(context.Background(), &kafka.ListGroupsRequest{})
	if err != nil {
		panic(err)
	}

	ret := ListGroupsResponse{}

	ret.groups = make([]string, len(lgresp.Groups))

	for i, t := range lgresp.Groups {
		ret.groups[i] = t.GroupID
	}

	return ret
}

func (r *ListGroupsResponse) Print() {
	for _, t := range r.groups {
		fmt.Println(t)
	}
}
