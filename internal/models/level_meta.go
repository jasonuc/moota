package models

import "math"

type LevelMeta struct {
	Level int64 `json:"level"`
	XP    int64 `json:"xp"`
}

func NewLeveLMeta(initalLevel, initialXp int64) LevelMeta {
	return LevelMeta{
		Level: initalLevel,
		XP:    initialXp,
	}
}

func (l *LevelMeta) levelUp() {
	l.Level++
}

func (l *LevelMeta) addXp(xp int64) {
	l.XP += xp

	// Handles scenario where plant needs to level up multiple times in one go
	for {
		nextLevelXpReq := xpRequiredForLevel(l.Level + 1)
		if l.XP >= nextLevelXpReq {
			l.XP -= nextLevelXpReq
			l.levelUp()
		} else {
			break
		}
	}
}

func xpRequiredForLevel(level int64) int64 {
	return int64(math.Round(75 * (math.Pow(float64(level), 2.0) - float64(level))))
}
