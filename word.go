package main

type Word struct {
	Word         string `json:"word"`
	Score        int64  `json:"score"`
	NumSyllables int    `json:"numSyllables"`
}
