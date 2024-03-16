package models

import (
	"fmt"
	"math/big"
	"math/rand"
	"strings"
	"time"
)

type Rope struct {
	weight int
	str    string
	right  *Rope
	left   *Rope
	parent *Rope
}

// method to find the leaf of the given index
func (r *Rope) Index(i int) (*Rope, error) {

	// if there is no left or right leaf then then the recursive indexing is at the desired leaf
	if r.right == nil && r.left == nil {
		return r, nil
	}

	// if the index is greater than the current leaf weight, subtract the weight and traverse right
	if i > r.weight {
		r, err := r.right.Index(i - r.weight)

		if err != nil {
			panic(err)
		}

		return r, nil
	}

	// if the index is not greater than the right leaf's weight, use the left leaf
	r, err := r.left.Index(i)

	if err != nil {
		panic(err)
	}

	return r, nil
}

// Split the string starting at the index designated
// re-traverse the tree and subtract the weight and set right-side nodes to nil
func (r *Rope) Split(index int) (*Rope, string, error) {
	time.Sleep(500 * time.Millisecond)
	if index > len(r.Collect()) {
		return r, "", fmt.Errorf("Index too high")
	}

	// closure to adjust parents weight and remove right node
	AdjustParentLeaf := func(leaf *Rope) *Rope {

		if leaf.parent != nil {
			newParentWeight := leaf.parent.weight
			if leaf.parent.right != nil && leaf.parent.right != leaf && leaf.left != nil {
				newParentWeight -= leaf.parent.right.weight
				leaf.parent.right = nil
			}

			leaf.parent.weight = newParentWeight
			return leaf.parent
		}
		return leaf
	}

	// get desired leaf that includes index
	r, err := r.Index(index)

	if err != nil {
		panic(err)
	}

	for {
		r = AdjustParentLeaf(r)

		if r.parent == nil {
			return r, r.Collect()[:index], nil
		}
	}

}

// function to collect all leafs in a rope and return the full string
func (r *Rope) Collect() string {
	rCopy := *r

	ropeStack := []Rope{}
	ropeList := []Rope{}

	getLeftMostNodeVal := func(rl Rope) {
		for {
			ropeStack = append(ropeStack, rl)

			if rl.left == nil {
				break
			}

			rl = *rl.left
		}
		ropeList = append(ropeList, rl)
	}

	getLeftMostNodeVal(rCopy)

	for len(ropeStack) != 0 {
		parent := ropeStack[len(ropeStack)-1]
		ropeStack = ropeStack[:len(ropeStack)-1]

		if parent.right != nil {
			getLeftMostNodeVal(*parent.right)
		}
	}

	var stringList []string

	for i := 0; i < len(ropeList); i++ {
		stringList = append(stringList, ropeList[i].str)
	}

	return strings.Join(stringList, "")
}

func (r *Rope) Concat(newR *Rope) *Rope {

	// shallow copy the rope so we don't have to traverse back up when concat-ing a new rope
	rCopy := *r
	var newWeight int

	// the new weight of the root node is based on the sum of the child nodes of the left rope
	// in the instance of this method, we are taking the original rope that is calling the method
	for {
		if rCopy.left == nil {
			if rCopy.parent.left != nil && rCopy.parent.right != nil {
				leftW := rCopy.parent.left.weight
				newWeight = leftW
				break
			}
		}
		rCopy = *rCopy.left
	}

	return &Rope{
		weight: newWeight,
		right:  newR,
		left:   r,
	}

}

func CreateRope(maxStrLen int, str string, parent *Rope) *Rope {
	if len(str) <= maxStrLen {
		return &Rope{
			str:    str,
			weight: len(str),
			parent: parent,
		}
	}

	splitNum := len(str) / 2

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
