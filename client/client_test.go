package client

import (
	"fmt"
	"testing"
	"time"
)

func TestClientWrite(t *testing.T) {
	raftNodes := []string{"http://127.0.0.1:7001/newTx", "http://127.0.0.1:8001/newTx", "http://127.0.0.1:9001/newTx"}

	client := NewClient("client1", raftNodes)

	for i := 0; i < 10; i++ {
		// time.Sleep(1 * time.Second)
		time.Sleep(5 * time.Nanosecond)
		res := client.CommonWrite([]string{fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i)})

		fmt.Printf("本轮(%d)结果 : %s\n", i, res)
	}
}

func TestClientRead(t *testing.T) {
	raftNodes := []string{"http://127.0.0.1:7001/newTx", "http://127.0.0.1:8001/newTx", "http://127.0.0.1:9001/newTx"}

	client := NewClient("client1", raftNodes)

	// res := client.CommonWrite([]string{"test_key", "test_val"})
	// fmt.Printf("写入结果 : %s\n", res)

	for i := 0; i < 10; i++ {
		// time.Sleep(1 * time.Second)
		time.Sleep(5 * time.Nanosecond)
		res := client.CommonRead([]string{"111"})

		fmt.Printf("本轮(%d)结果 : %s\n", i, res)
	}
}
