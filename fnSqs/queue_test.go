package fnSqs

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestQueue(test *testing.T) {
	var conn *Conn
	var err error

	if conn, err = NewConnWithAuth(&INewConn{
		Access: "",
		Secret: "",
		Region: "ap-northeast-2",
	}); err != nil {
		test.Fatal(err)
	}

	var newQueueNm = func() string {
		return fmt.Sprintf("queue-%s", uuid.NewString()[8:])
	}

	test.Run("create, delete queue", func(t *testing.T) {
		var ctx = context.TODO()
		var queueNm = newQueueNm()
		if err = conn.Create(ctx, queueNm); err != nil {
			t.Fatal(err)
		}

		if err = conn.Delete(ctx, queueNm); err != nil {
			t.Fatal(err)
		}
	})

	test.Run("send, receive queue", func(t *testing.T) {
		var ctx = context.TODO()
		var queueNm = newQueueNm()
		if err = conn.Create(ctx, queueNm); err != nil {
			t.Fatal(err)
		}

		defer func() {
			if err = conn.Delete(ctx, queueNm); err != nil {
				t.Fatal(err)
			}
		}()

		type Body struct {
			Name string `json:"name"`
		}

		var body = &Body{
			Name: fmt.Sprintf("hello-%s", uuid.NewString()[8:]),
		}

		var queue *Queue[Body]
		if queue, err = NewQueue[Body](ctx, conn, queueNm); err != nil {
			t.Fatal(err)
		}

		if err = queue.Send(ctx, body); err != nil {
			t.Fatal(err)
		}

		time.Sleep(3 * time.Second)
		queue.Receiver(func(msg *Body) (err error) {
			assert.Equal(t, body.Name, msg.Name)

			fmt.Printf("test done!\n")
			t.Skip()
			return
		})

	})

}
