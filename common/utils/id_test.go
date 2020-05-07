package utils

import "testing"

func BenchmarkWorker_GenerateId(b *testing.B) {
	worker, err := newWorker(2)
	if err != nil {
		b.Fatalf("NewWorker Err: %v", err)
	}
	for i := 0; i < b.N; i++ {
		id := worker.generateID()
		b.Log(id)
	}
}

func TestGetId(t *testing.T) {
	worker, err := newWorker(2)
	if err != nil {
		t.Fatalf("NewWorker Err: %v", err)
	}
	id := worker.generateID()
	t.Log(id)
	err = CheckID(id)
	if err != nil {
		t.Error(err)
	}

	t.Log(8191 | 10)
}

func TestW(t *testing.T) {
	a := 1       // 0001
	b := 2       // 0010
	c := 3       // 0011
	t.Log(a | b) // 0011
	t.Log(a | c) // 0011
	t.Log(b | c) // 0011
}
