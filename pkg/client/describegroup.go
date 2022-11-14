package client

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type DescribeGroupResponse struct {
	topics []DescribedTopic
}

type DescribedTopic struct {
	parts []DescribedPartitions
}

type DescribedPartitions struct {
	ID            int
	CurrentOffset int
	LogEndOffset  int
}

func DescribeGroup(c *kafka.Client, group string) DescribeGroupResponse {
	resp, err := c.DescribeGroups(context.Background(), &kafka.DescribeGroupsRequest{
		GroupIDs: []string{group},
	})
	if err != nil {
		panic(err)
	}

	qOffsets := make(chan map[int]int64)
	qLogEndOffsets := make(chan map[int]int64)

	fmt.Printf("GROUP TOPIC PARTITION CURRENT-OFFSET LOG-END-OFFSET LAG CONSUMER-ID HOST CLIENT-ID\n")
	for _, g := range resp.Groups {
		for _, m := range g.Members {
			for _, t := range m.MemberAssignments.Topics {
				go Offsets(c, t, g.GroupID, qOffsets)
				go LogEndOffsets(c, t.Topic, t.Partitions, qLogEndOffsets)

				offsets := <-qOffsets
				logEndOffsets := <-qLogEndOffsets
				for p, commitedOffset := range offsets {
					fmt.Printf("%s %s %d %d %d %d %s %s %s\n", g.GroupID, t.Topic, p, commitedOffset, logEndOffsets[p], logEndOffsets[p]-commitedOffset, m.MemberID, m.ClientHost, m.ClientID)
				}
			}
		}
	}

	return DescribeGroupResponse{}
}

func (r *DescribeGroupResponse) Print() {
	fmt.Println("w00t")
}

func Offsets(c *kafka.Client, t kafka.GroupMemberTopic, group string, q chan map[int]int64) {
	ofresp, err := c.OffsetFetch(context.Background(), &kafka.OffsetFetchRequest{
		GroupID: group,
		Topics:  map[string][]int{t.Topic: t.Partitions},
	})
	if err != nil {
		panic(err)
	}

	ret := make(map[int]int64)
	for _, p := range ofresp.Topics[t.Topic] {
		ret[p.Partition] = p.CommittedOffset
	}

	q <- ret
}

func LogEndOffsets(c *kafka.Client, topic string, parts []int, q chan map[int]int64) {
	t := make(map[string][]kafka.OffsetRequest)

	partitions := make([]kafka.OffsetRequest, len(parts))
	for i, p := range parts {
		partitions[i].Partition = p
		partitions[i].Timestamp = -1
	}

	t[topic] = partitions

	loresp, err := c.ListOffsets(context.Background(), &kafka.ListOffsetsRequest{
		Topics: t,
	})
	if err != nil {
		panic(err)
	}

	offsets := loresp.Topics[topic]

	ret := make(map[int]int64)
	for _, x := range offsets {
		ret[x.Partition] = x.LastOffset
	}

	q <- ret
}
