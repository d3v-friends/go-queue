package fnSqs

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"testing"
	"time"
)

func TestQueue(test *testing.T) {
	var conn, err = NewConn()
	if err != nil {
		test.Fatal(err)
	}

	var ls []string
	if ls, err = conn.List(context.Background()); err != nil {
		test.Fatal(err)
	}

	fmt.Printf("queue list: %s\n", ls)

	var queue *Queue
	if queue, err = conn.NewQueue(context.TODO(), "dev-sender"); err != nil {
		test.Fatal(err)
	}

	var chErr = make(chan error)
	defer close(chErr)

	go queue.Receiver(func(msg types.Message) (err error) {
		fmt.Printf("Received: %s\n", *msg.Body)
		return
	}, chErr)

	go func() {
		fmt.Printf("chErr: err=%s", <-chErr)
	}()

	test.Run("send message", func(t *testing.T) {
		var ctx = context.TODO()
		if err = queue.Send(ctx, "Hello world"); err != nil {
			t.Fatal(err)
		}

		time.Sleep(60 * time.Second)
	})
}
