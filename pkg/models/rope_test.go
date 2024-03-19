package models

import (
	"testing"
)

const (
	START_ROPE_STR_SIZE int = 20000
	ROPE_TEST_STR           = "Hello_World_I_am_a_Rope!"
	ROPE_TEST_STR_LEN       = len(ROPE_TEST_STR)
)

func Test_RopeCollect(t *testing.T) {

	rope := CreateRandomRope(10, START_ROPE_STR_SIZE)

	_, ropeStr := rope.Collect()

	if len(ropeStr) != START_ROPE_STR_SIZE {
		t.Error()
	}
}

func Test_RopeIndex(t *testing.T) {

	rope := CreateRope(4, ROPE_TEST_STR, nil)
	rope, _, err := rope.Index(8)

	if err != nil || rope.str != "Wor" {
		t.Error()
	}

}

func Test_RopeConcat(t *testing.T) {

	rope := CreateRope(4, ROPE_TEST_STR, nil)
	rope2 := CreateRope(4, ROPE_TEST_STR, nil)

	out := rope.Concat(rope2)

	_, ropeStr := out.Collect()

	if out.weight != ROPE_TEST_STR_LEN {
		t.Error("invalid weight")
	}

	if ropeStr != ROPE_TEST_STR+ROPE_TEST_STR {
		t.Error("unexpected rope string")
	}
}

func Test_RopeSplit(t *testing.T) {

	rope := CreateRope(4, ROPE_TEST_STR, nil)

	// "Hello_World_I_am_a_Rope!"
	_, newStr, _, err := rope.Split(10)

	if err != nil || newStr != "Hello_World" {
		t.Error()
	}
}

func Test_RopeRebalance(t *testing.T) {

	rope1 := CreateRope(4, ROPE_TEST_STR, nil)
	rope2 := CreateRope(7, ROPE_TEST_STR, nil)
	rope3 := CreateRope(3, ROPE_TEST_STR, nil)

	ropeList := []*Rope{rope1, rope2, rope3}

	rebalancedRope, err := RebalanceRope(ropeList, 0, len(ropeList)-1)

	if err != nil {
		panic(err)
	}

	_, ropeStr := rebalancedRope.Collect()

	if ropeStr != "Hello_World_I_am_a_Rope!Hello_World_I_am_a_Rope!Hello_World_I_am_a_Rope!" {
		t.Error("rope not balanced correctly")
	}
}

func Test_RopeInsert(t *testing.T) {

	rope1 := CreateRope(4, ROPE_TEST_STR, nil)

	newRope := rope1.Insert(2, "INSERT")

	_, ropeStr := newRope.Collect()

	if ropeStr != "HelINSERTlo_World_I_am_a_Rope!" {
		t.Error("rope not inserted correctly")
	}

	rope1 = CreateRope(4, ROPE_TEST_STR, nil)

	newRope = rope1.Insert(10, "INSERT")

	_, ropeStr = newRope.Collect()

	if ropeStr != "Hello_WorldINSERT_I_am_a_Rope!" {
		t.Error("rope not inserted correctly")
	}

	rope1 = CreateRope(4, ROPE_TEST_STR, nil)

	newRope = rope1.Insert(20, "INSERT")

	_, ropeStr = newRope.Collect()

	if ropeStr != "Hello_World_I_am_a_RoINSERTpe!" {
		t.Error("rope not inserted correctly")
	}
}
