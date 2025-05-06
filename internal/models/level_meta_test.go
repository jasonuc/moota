package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLevelMeta(t *testing.T) {
	t.Run("add xp but not enough to increase level", func(t *testing.T) {
		lm := &LevelMeta{Level: 1, XP: 0}
		expectedLevelMeta := NewLeveLMeta(1, 25)

		lm.addXp(25)

		assert.Equal(t, lm.XP, expectedLevelMeta.XP)
		assert.Equal(t, lm.Level, expectedLevelMeta.Level)
	})

	t.Run("add enough xp for one level up", func(t *testing.T) {
		lm := &LevelMeta{Level: 1, XP: 0}
		expectedLevelMeta := NewLeveLMeta(2, 25)

		xpToAdd := xpRequiredForLevel(expectedLevelMeta.Level) + 25
		lm.addXp(xpToAdd)

		assert.Equal(t, lm.XP, expectedLevelMeta.XP)
		assert.Equal(t, lm.Level, expectedLevelMeta.Level)
	})

	t.Run("add enough xp for to level up to occur multiple times", func(t *testing.T) {
		lm := &LevelMeta{Level: 1, XP: 0}
		expectedLevelMeta := NewLeveLMeta(5, 10)

		xpToAdd := xpRequiredForLevel(2) + xpRequiredForLevel(3) + xpRequiredForLevel(4) + xpRequiredForLevel(5) + expectedLevelMeta.XP
		lm.addXp(xpToAdd)

		assert.Equal(t, lm.XP, expectedLevelMeta.XP)
		assert.Equal(t, lm.Level, expectedLevelMeta.Level)
	})
}
