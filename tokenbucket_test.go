package TokenBucket

import (
	"testing"
	"time"
	"fmt"
)

func TestTokenBucket(t *testing.T) {
	var rate uint32 = 50
	token := NewTokenBucket(rate)
	if token == nil {
		t.Error("Failed to create token bucket")
	}
		
	routine := func(name string) {
		for {
			e := token.GetToken()
			if e == nil {
				fmt.Printf("%s got token\n", name)
			} else {
				fmt.Printf("%s must wait\n", name)
				time.Sleep(time.Millisecond * 100)
			}
		}
	}
	
	go routine("Thread 1")
	go routine("Thread 2")
	time.Sleep(2*time.Second)
	newToken := NewTokenBucket(rate*2)
	oldToken := token
	token = newToken
	fmt.Printf("Changing token, double the rate...")
	oldToken.Stop()
	time.Sleep(2*time.Second)
	token.Stop()
	fmt.Printf("Stopping...")
	time.Sleep(1*time.Second)
}
