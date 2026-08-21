package main

import (
	"math/rand"
	"strconv"
	"time"
)

type some struct {
	ID     string
	Name   string
	Images []string

	Status uint8

	someChild someChild

	user user
}

type user struct {
	ID        string
	FirstName string
	LastName  string
	Email     string
	Childs    []int64

	someChild someChild
}

type someChild struct {
	ID     string
	Status uint64

	Boxes [][]string

	Place Coordinates
}

type Coordinates struct {
	X, Y int64
}

func randString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func NewSomeSlice(countItems int) []*some {
	// // Чтобы структура была реально большой — генерим большие срезы строк
	// images := make([]string, 1000)
	// for i := range images {
	// 	images[i] = randString(256) // ~256 байт на строку
	// }

	// boxes := make([][]string, 50)
	// for i := range boxes {
	// 	row := make([]string, 100)
	// 	for j := range row {
	// 		row[j] = randString(128) // ~128 байт
	// 	}
	// 	boxes[i] = row
	// }

	// childs := make([]int64, 5000)
	// for i := range childs {
	// 	childs[i] = rand.Int63()
	// }

	res := make([]*some, countItems)

	for i := range res {
		res[i] = &some{
			ID:     strconv.FormatInt(time.Now().UnixNano(), 36),
			Name:   randString(6),
			Images: nil,
			Status: uint8(rand.Intn(255)),
			someChild: someChild{
				ID:     strconv.FormatInt(time.Now().UnixNano(), 36),
				Status: rand.Uint64(),
				Boxes:  nil,
				Place: Coordinates{
					X: rand.Int63(),
					Y: rand.Int63(),
				},
			},
			user: user{
				ID:        strconv.FormatInt(time.Now().UnixNano(), 36),
				FirstName: randString(6),
				LastName:  randString(6),
				Email:     randString(12),
				Childs:    nil,
				someChild: someChild{
					ID:     strconv.FormatInt(time.Now().UnixNano(), 36),
					Status: rand.Uint64(),
					Boxes:  nil,
					Place: Coordinates{
						X: rand.Int63(),
						Y: rand.Int63(),
					},
				},
			},
		}
	}

	return res
}

func NewSome() *some {
	// Чтобы структура была реально большой — генерим большие срезы строк
	images := make([]string, 1000)
	for i := range images {
		images[i] = randString(256) // ~256 байт на строку
	}

	boxes := make([][]string, 50)
	for i := range boxes {
		row := make([]string, 100)
		for j := range row {
			row[j] = randString(128) // ~128 байт
		}
		boxes[i] = row
	}

	childs := make([]int64, 5000)
	for i := range childs {
		childs[i] = rand.Int63()
	}

	return &some{
		ID:     strconv.FormatInt(time.Now().UnixNano(), 36),
		Name:   randString(64),
		Images: images,
		Status: uint8(rand.Intn(255)),
		someChild: someChild{
			ID:     strconv.FormatInt(time.Now().UnixNano(), 36),
			Status: rand.Uint64(),
			Boxes:  boxes,
			Place: Coordinates{
				X: rand.Int63(),
				Y: rand.Int63(),
			},
		},
		user: user{
			ID:        strconv.FormatInt(time.Now().UnixNano(), 36),
			FirstName: randString(32),
			LastName:  randString(32),
			Email:     randString(64),
			Childs:    childs,
			someChild: someChild{
				ID:     strconv.FormatInt(time.Now().UnixNano(), 36),
				Status: rand.Uint64(),
				Boxes:  boxes,
				Place: Coordinates{
					X: rand.Int63(),
					Y: rand.Int63(),
				},
			},
		},
	}
}

func Equal[T any](a, b []*T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
