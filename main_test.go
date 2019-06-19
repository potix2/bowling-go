package main

import (
	"strconv"
	"strings"
	"testing"
)

type Bowling struct {
	bonusCounters []int
}

type Roll struct {
	frame int
	hit int
	isSpare bool
	isStrike bool
}

func newGame() *Bowling {
	return &Bowling{ []int{}}
}

func Parse(marks string) []Roll {
	frame := 0
	firstHitCurrentFrame := 0
	isFirst := true
	result := make([]Roll, 0)
	for _, mark := range strings.Split(marks, "") {
		if isFirst {
			frame = frame + 1
		}

		if mark == "/" {
			hit := 10 - firstHitCurrentFrame
			result = append(result, Roll{frame, hit, true, false})
			isFirst = true
		} else if mark == "X" {
			if frame <= 10 {
				result = append(result, Roll{frame, 10, false, true})
			} else {
				result = append(result, Roll{frame, 10, false, false})
			}
			isFirst = true
		} else {
			if i, err := strconv.ParseInt(mark, 10, 32); err == nil {
				result = append(result, Roll{frame, int(i), false, false})
				firstHitCurrentFrame = int(i)
			} else {
				result = append(result, Roll{frame, 0, false, false})
			}

			isFirst = !isFirst
		}
	}

	return result
}

func (b *Bowling) CalcScore(marks string) int {
	var score int
	for _, roll := range Parse(marks) {
		score = score + b.updateBonus(roll.hit)

		if roll.frame <= 10 {
			score = score + roll.hit

			if roll.isSpare {
				b.addBonusCounter(1)
			} else if roll.isStrike {
				b.addBonusCounter(2)
			}
		}
	}

	return score
}

func (b *Bowling) updateBonus(hit int) int {
	bonusTotal := 0
	newBonusCounters := make([]int, 0)
	for _, counter := range b.bonusCounters {
		bonusTotal = bonusTotal + hit
		if counter-1 > 0 {
			newBonusCounters = append(newBonusCounters, counter-1)
		}
	}
	b.bonusCounters = newBonusCounters
	return bonusTotal
}

func (b *Bowling) addBonusCounter(count int) {
	b.bonusCounters = append(b.bonusCounters, count)
}

func TestParse(t *testing.T) {
	testCases := []struct {
		in       string
		expected []Roll
	}{
		{"9-9-9-9-9-9-9-9-9-9-",
			[]Roll {
				{1, 9, false, false},
				{1, 0, false, false},
				{2, 9, false, false},
				{2, 0, false, false},
				{3, 9, false, false},
				{3, 0, false, false},
				{4, 9, false, false},
				{4, 0, false, false},
				{5, 9, false, false},
				{5, 0, false, false},
				{6, 9, false, false},
				{6, 0, false, false},
				{7, 9, false, false},
				{7, 0, false, false},
				{8, 9, false, false},
				{8, 0, false, false},
				{9, 9, false, false},
				{9, 0, false, false},
				{10, 9, false, false},
				{10, 0, false, false},
			},
		},
		{"5/5/5/5/5/5/5/5/5/5/5",
			[]Roll {
				{1, 5, false, false},
				{1, 5, true, false},
				{2, 5, false, false},
				{2, 5, true, false},
				{3, 5, false, false},
				{3, 5, true, false},
				{4, 5, false, false},
				{4, 5, true, false},
				{5, 5, false, false},
				{5, 5, true, false},
				{6, 5, false, false},
				{6, 5, true, false},
				{7, 5, false, false},
				{7, 5, true, false},
				{8, 5, false, false},
				{8, 5, true, false},
				{9, 5, false, false},
				{9, 5, true, false},
				{10, 5, false, false},
				{10, 5, true, false},
				{11, 5, false, false},
			},
		},
		{"XXXXXXXXXXXX",
			[]Roll {
				{1, 10, false, true},
				{2, 10, false, true},
				{3, 10, false, true},
				{4, 10, false, true},
				{5, 10, false, true},
				{6, 10, false, true},
				{7, 10, false, true},
				{8, 10, false, true},
				{9, 10, false, true},
				{10, 10, false, true},
				{11, 10, false, false},
				{12, 10, false, false},
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.in, func(t *testing.T) {
			actual := Parse(tt.in)
			for i, _ := range tt.expected {
				if actual[i] != tt.expected[i] {
					t.Errorf("got %v, expected %v", actual, tt.expected)
				}
			}
		})
	}
}

func TestCalcScore(t *testing.T) {
	testCases := []struct {
		in       string
		expected int
	}{
		{"9-9-9-9-9-9-9-9-9-9-", 90},
		{"8-8-8-8-8-8-8-8-8-8-", 80},
		{"--------------------", 0},
		{"5/5/5/5/5/5/5/5/5/5/5", 150},
		{"5/5-----------------", 20},
		{"5/52----------------", 22},
		{"5/-5----------------", 15},
		{"X------------------", 10},
		{"X1-----------------", 12},
		{"X11----------------", 14},
		{"X111---------------", 15},
		{"XXXXXXXXXXXX", 300},
	}

	for _, tt := range testCases {
		t.Run(tt.in, func(t *testing.T) {
			actual := newGame().CalcScore(tt.in)
			if actual != tt.expected {
				t.Errorf("got %v, expected %v", actual, tt.expected)
			}
		})
	}
}
