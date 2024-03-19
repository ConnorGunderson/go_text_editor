package models

import (
	"fmt"
	"math/big"
	"math/rand"
	"strings"
)

type Rope struct {
	weight int
	str    string
	right  *Rope
	left   *Rope
	parent *Rope
}

// method to find the leaf of the given index
func (r *Rope) Index(i int) (*Rope, int, error) {

	// if there is no left or right leaf then then the recursive indexing is at the desired leaf
	if r.right == nil && r.left == nil {
		return r, i, nil
	}

	// if the index is greater than the current leaf weight, subtract the weight and traverse right
	if i > r.weight {
		r, rIndex, err := r.right.Index(i - r.weight)

		if err != nil {
			panic(err)
		}

		return r, rIndex, nil
	}

	// if the index is not greater than the right leaf's weight, use the left leaf
	r, lIndex, err := r.left.Index(i)

	if err != nil {
		panic(err)
	}

	return r, lIndex, nil
}

// Split the string starting at the index designated
// re-traverse the tree and subtract the weight and set right-side nodes to nil
func (r *Rope) Split(index int) (leftRope *Rope, ropeStr string, rightRope *Rope, err error) {

	_, ropeStr = r.Collect()

	if index > len(ropeStr) {
		return r, "", nil, fmt.Errorf("Index too high")
	}

	rightRopeList := []*Rope{}
	weightToSub := 0

	// closure to adjust parents weight and remove right node
	AdjustParentLeaf := func(leaf *Rope) *Rope {
		if leaf.parent != nil {

			if leaf.parent.right != nil && leaf.parent.right != leaf {
				rightRopeList = append(rightRopeList, leaf.parent.right)
				weightToSub += leaf.parent.right.weight

				if leaf.parent.parent != nil && leaf.parent.parent.left == leaf.parent {
					leaf.parent.parent.weight -= weightToSub
				}

				leaf.parent.right = nil
			}

			return leaf.parent
		}
		return leaf
	}

	// get desired leaf that includes index
	r, leftoverIndex, err := r.Index(index)

	if leftoverIndex > 0 {

		ropeFromLeftovers := CreateRope(leftoverIndex, r.str, r.parent)

		if r == r.parent.left {
			r.parent.left = ropeFromLeftovers
		} else {
			r.parent.right = ropeFromLeftovers
		}

		r = ropeFromLeftovers

		_, asdf := r.Collect()

		fmt.Println(asdf, r.left, leftoverIndex)

		r, _, err = r.Index(index)

		if err != nil {
			panic(err)
		}
	}

	if err != nil {
		panic(err)
	}

	for {
		r = AdjustParentLeaf(r)

		if r.parent == nil {
			// we reconstruct the right rope from the rightRopeList populated by the AdjustParentLeaf closure
			rightRope, err = RebalanceRope(rightRopeList, 0, len(rightRopeList))

			if err != nil {
				panic(err)
			}

			_, ropeStr := r.Collect()

			return r, ropeStr, rightRope, nil
		}
	}

}

func RebalanceRope(rl []*Rope, start int, end int) (*Rope, error) {

	if end > len(rl) {
		return nil, fmt.Errorf("end is greater than the length of the rope list")
	}

	rnge := end - start

	if rnge == 1 {
		return rl[0], nil
	}

	var outRope *Rope

	for _, rope := range rl {
		if outRope == nil {
			outRope = rope
		} else {
			outRope = outRope.Concat(rope)
		}

	}

	return outRope, nil
}

// function to collect all leafs in a rope and return the full string
func (r *Rope) Collect() ([]*Rope, string) {
	rCopy := *r

	ropeStack := []Rope{}
	ropeList := []*Rope{}

	getLeftMostNodeVal := func(rl *Rope) {
		for {
			ropeStack = append(ropeStack, *rl)

			if rl.left == nil {
				break
			}

			rl = rl.left
		}

		ropeList = append(ropeList, rl)
	}

	getLeftMostNodeVal(&rCopy)

	for len(ropeStack) != 0 {
		parent := ropeStack[len(ropeStack)-1]
		ropeStack = ropeStack[:len(ropeStack)-1]

		if parent.right != nil {
			getLeftMostNodeVal(parent.right)
		}
	}

	var stringList []string

	for i := 0; i < len(ropeList); i++ {
		stringList = append(stringList, ropeList[i].str)
	}

	return ropeList, strings.Join(stringList, "")
}

func (r *Rope) Concat(newR *Rope) *Rope {
	rLeafList, _ := r.Collect()

	var newWeight int

	for _, leaf := range rLeafList {
		newWeight += leaf.weight
	}

	newRope := &Rope{
		weight: newWeight,
		right:  newR,
		left:   r,
	}

	newRope.right.parent = newRope
	newRope.left.parent = newRope

	return newRope
}

func (r *Rope) Insert(idx int, str string) *Rope {
	left, _, right, err := r.Split(idx)

	if err != nil {
		panic(err)
	}

	left = left.Concat(CreateRope(10, str, nil))

	left = left.Concat(right)

	return left
}

func CreateRope(maxStrLen int, str string, parent *Rope) *Rope {

	if len(str) <= maxStrLen {
		fmt.Println(str)
		return &Rope{
			str:    str,
			weight: len(str),
			parent: parent,
		}
	}

	splitNum := len(str) / 2
	fmt.Println(splitNum)

	str1 := str[:splitNum]
	str2 := str[splitNum:]

	newRopeLeaf := Rope{
		parent: parent,
	}

	if parent == nil {
		newRopeLeaf.weight = len(str)
		newRopeLeaf.left = CreateRope(maxStrLen, str, &newRopeLeaf)
		newRopeLeaf.right = nil

	} else {
		newRopeLeaf.weight = len(str1)
		newRopeLeaf.left = CreateRope(maxStrLen, str1, &newRopeLeaf)
		newRopeLeaf.right = CreateRope(maxStrLen, str2, &newRopeLeaf)
	}

	return &newRopeLeaf
}

func CreateRandomRope(maxStrLen int, strLen int) *Rope {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, strLen)

	for i := range bytes {
		bi := big.NewInt(int64(len(charset))).Int64()

		num := rand.Intn(int(bi))

		bytes[i] = charset[num]
	}

	return CreateRope(maxStrLen, string(bytes), nil)
}
