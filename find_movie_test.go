package rotten_tomatoes

import (
	"testing"
)

func findMovie(title string, year int, c *Config, t *testing.T) *RottenTomatoMovie {
	match, err := c.FindMovie(title, year)
	if err != nil {
		t.Error(err)
	} else if len(match) == 0 {
		t.Error("No match found", match)
		return nil
	} else if len(match) > 1 {
		if year != 0 {
			t.Error("Ambiguous results dispite year", match)
		}
	}

	return match[0]
}

func testInit(t *testing.T) *Config {
	c, err := InitSetup()
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestFindMovie(t *testing.T) {
	c := testInit(t)

	findMovie("up", 2009, c, t)
	findMovie("toy story", 0, c, t)
	findMovie("toy story", 1995, c, t)
	findMovie("toy story 3", 0, c, t)
	findMovie("A Good Year", 2006, c, t)
	findMovie("Waltz with Bashir", 0, c, t)
	findMovie("Anna Karenina", 2012, c, t)
	findMovie("The Informers", 2008, c, t)
	findMovie("The Informers", 2009, c, t)
	findMovie("The Sting", 1973, c, t)
	findMovie("Gainsbourg", 2010, c, t)

	// Ambiguous
	findMovie("Anna Karenina", 0, c, t)
	findMovie("The Thomas Crown Affair", 0, c, t)
	findMovie("Dredd", 0, c, t)
	findMovie("Eulogy", 0, c, t)
	findMovie("Gainsbourg", 0, c, t)
}

func TestActors(t *testing.T) {
	c := testInit(t)
	m := findMovie("Let me in", 2010, c, t)
	if len(m.AbridgedCast) != 5 {
		t.Error("Cast not right", len(m.AbridgedCast), m.AbridgedCast)
	}
	err := m.GetFullCast(c)
	if err != nil || len(m.FullCast) != 30 {
		t.Error("Full cast not right", len(m.FullCast), m.FullCast, err)
	}
}

func TestScoredMovies(t *testing.T) {
	m := &RottenTomatoMovie{Year: 1980}
	s := []scoredMovie{
		{4, 0.0, m},
		{4, 0.0, m},
		{2, 0.1, m},
		{2, 0.7, m},
		{1, 0.9, m},
		{1, 1.3, m},
		{0, 2.0, m},
	}

	best1 := selectBest(s)
	if len(best1) != 2 {
		t.Error("Wrong best1", best1)
	}
	best2 := selectBest(s[2:])
	if len(best2) != 1 {
		t.Error("Wrong best2", best2)
	}
	best3 := selectBest(s[3:])
	if len(best3) != 4 {
		t.Error("Wrong best3", best3)
	}
	best4 := selectBest(s[4:])
	if len(best4) != 3 {
		t.Error("Wrong best4", best4)
	}
}
