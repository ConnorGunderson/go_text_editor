package models

import (
	"testing"
)

const (
	START_ROPE_STR_SIZE int = 20000
	ROPE_TEST_STR           = "Hello_World_I_am_a_Rope!"
)

func Test_RopeCollect(t *testing.T) {

	rope := CreateRandomRope(10, START_ROPE_STR_SIZE)

	if len(rope.Collect()) != START_ROPE_STR_SIZE {
		t.Error()
	}
}

func Test_RopeIndex(t *testing.T) {

	rope := CreateRope(4, ROPE_TEST_STR, nil)
	rope, err := rope.Index(8)

	if err != nil || rope.str != "Wor" {
		t.Error()
	}
}

func Test_RopeConcat(t *testing.T) {

	rope := CreateRope(4, ROPE_TEST_STR, nil)
	rope2 := CreateRope(4, ROPE_TEST_STR, nil)

	out := rope.Concat(rope2)

	if out.Collect() != ROPE_TEST_STR+ROPE_TEST_STR {
		t.Error()
	}
}

func Test_RopeSplit(t *testing.T) {

	rope := CreateRope(4, ROPE_TEST_STR, nil)

	// "Hello_World_I_am_a_Rope!"
	_, newStr, err := rope.Split(20)

	if err != nil || newStr != "Hello_World_I_am_a_R" {
		t.Error()
	}
}
