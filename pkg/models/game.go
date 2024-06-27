package models

import (
	"github.com/google/uuid"
	"sync"
)

type Score struct {
	Fire  int `json:"fire"`
	Water int `json:"water"`
	Food  int `json:"food"`
}

type ClickMessage struct {
	UserId uuid.UUID `json:"userId"`
	Score
}

type PlayerScore struct {
	Score
	Id    uuid.UUID `json:"id"`
	Mutex sync.Mutex
}

type GlobalScore struct {
	Score
	Mutex sync.Mutex
}

type BroadcastScore struct {
	UserId uuid.UUID    `json:"userId"`
	Ps     *PlayerScore `json:"ps"`
	Gs     *GlobalScore `json:"gs"`
}

func (ps *PlayerScore) UpdateScore(fire, water, food int) {
	ps.Mutex.Lock()
	defer ps.Mutex.Unlock()
	ps.Fire += fire
	ps.Water += water
	ps.Food += food
}

func (gs *GlobalScore) UpdateScore(fire, water, food int) {
	gs.Mutex.Lock()
	defer gs.Mutex.Unlock()
	gs.Fire += fire
	gs.Water += water
	gs.Food += food
}
