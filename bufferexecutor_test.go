package bufferexecutor

import (
	"reflect"
	"runtime"
	"testing"
	//"time"
)

func Test_Put1(t *testing.T) {

	var value []interface{}

	for i := 0; i < 99; i++ {
		value = append(value, i)
	}

	var expect []interface{}
	p := New(3, 1000000, func(data []interface{}) {
		expect = append(expect, data...)
	})

	for _, v := range value {
		p.Put(v)
	}


	runtime.Gosched()

	if reflect.DeepEqual(value, expect) {

	} else {
		t.Errorf("expect %v got %v", value, expect)
	}
}

func Test_Put(t *testing.T) {
	var value = []interface{}{4, 5, 6}
	var expect []interface{}
	p := New(3, 1000000, func(data []interface{}) {
		expect = append(expect, data...)
	})

	for _, v := range value {
		p.Put(v)

	}
	runtime.Gosched()

	if reflect.DeepEqual(value, expect) {

	} else {
		t.Errorf("expect %v got %v", value, expect)
	}
}

func Test_Close(t *testing.T) {

	var value = []interface{}{4, 5, 6}
	var expect []interface{}
	p := New(1000, 1000000, func(data []interface{}) {
		expect = append(expect, data...)
	})

	for _, v := range value {
		p.Put(v)
	}
	p.Close()
	runtime.Gosched()

	if reflect.DeepEqual(value, expect) {

	} else {
		t.Errorf("expect %v got %v", value, expect)
	}
}

func Test_Signal(t *testing.T) {
	var value = []interface{}{4}
	var expect []interface{}
	p := New(1000, 1000000, func(data []interface{}) {
		expect = append(expect, data...)
	})

	for _, v := range value {
		p.Put(v)

	}

	runtime.Gosched()
	p.Signal()

	runtime.Gosched()

	if reflect.DeepEqual(value, expect) {

	} else {
		t.Errorf("expect %v got %v", value, expect)
	}

}

