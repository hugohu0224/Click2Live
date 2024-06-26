package models

import "sync"

type Score struct {
	Fire  int `json:"fire"`
	Water int `json:"water"`
	Food  int `json:"food"`
}

type ClickMessage struct {
	Score
}

type PlayerScore struct {
	Score
	Mutex sync.Mutex
}

type GlobalScore struct {
	Score
	Mutex sync.Mutex
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
