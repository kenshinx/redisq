package redisq

import "testing"

var (
	msgs = []string{"a", "b", "c", "d"}
)

func setup() {
	rq.Clear()
	rq.Put(msgs)
}

func TestPut(t *testing.T) {
	setup()
	err := rq.Put(msgs)
	if err != nil {
		t.Fatal(err)
	}

}

func TestGet(t *testing.T) {
	setup()
	msg, err := rq.Get(true, 1)
	if err != nil {
		t.Fatal(err)
	}
	if string(msg) != msgs[0] {
		t.Fatalf("Set %s,but get %s", msgs[0], string(msg))
	}
}

func TestConsume(t *testing.T) {
	setup()

}

func TestLength(t *testing.T) {
	setup()
	if rq.Length() != 4 {
		t.Fatalf("Expect redis queue length is %d, but get %d", len(msgs), rq.Length())
	}
}

func TestEmpty(t *testing.T) {
	setup()
	rq.Clear()
	if rq.Empty() != true {
		t.Fatal("Exception,Redis queue is empty now")
	}
	rq.Put(msgs)
	if rq.Empty() {
		t.Fatal("Exception,Redis queue is not empty now")
	}
}
