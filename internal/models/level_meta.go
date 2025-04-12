package models

import "math"

type LevelMeta struct {
	Level int64
	Xp    int64
}

func NewLeveLMeta(initalLevel, initialXp int64) LevelMeta {
	return LevelMeta{
		Level: initalLevel,
		Xp:    initialXp,
	}
}

func (l *LevelMeta) levelUp() {
	l.Level++
}

func (l *LevelMeta) addXp(xp int64) {
	l.Xp += xp

	// Handles scenario where plant needs to level up multiple times in one go
	for {
		nextLevelXpReq := xpRequiredForLevel(l.Level + 1)
		if l.Xp >= nextLevelXpReq {
			l.Xp -= nextLevelXpReq
			l.levelUp()
		} else {
			break
		}
	}
}

func xpRequiredForLevel(level int64) int64 {
	return int64(math.Round(75 * (math.Pow(float64(level), 2.0) - float64(level))))
}
