package fnSqs

import (
	"context"
	"testing"
)

func TestQueue(test *testing.T) {
	var conn, err = NewConn()
	if err != nil {
		test.Fatal(err)
	}

	var queue *Queue
	if queue, err = conn.NewQueue(context.TODO(), "dev-sender"); err != nil {
		test.Fatal(err)
	}

}
