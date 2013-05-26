package rotten_tomatoes

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"utils/intmath"
	"utils/json"
	"utils/strtools"
)

type RottenTomatoActor struct {
	Name       string
	Id         string
	Characters []string
}

type RottenTomatoMovie struct {
	Id      string
	Title   string
	Year    int
	Runtime int
	Ratings struct {
		CriticsScore  int `json:"critics_score"`
		AudienceScore int `json:"audience_score"`
	}
	Posters struct {
		Thumbnail string
		Profile   string
		Detailed  string
		Original  string
	}
	AbridgedCast []RottenTomatoActor `json:"abridged_cast"`
	FullCast     []RottenTomatoActor `json:"cast"`
	AlternateIds struct {
		Imdb string
	} `json:"alternate_ids"`

	Links struct {
		Self      string
		Alternate string
		Cast      string
		Clips     string
		Reviews   string
		Similar   string
	}
}

func (r *RottenTomatoMovie) GetFullCast(c *Config) error {
	if r.Links.Cast == "" {
		return nil
	}

	url, err := c.AddKey(r.Links.Cast)
	if err != nil {
		return err
	}
	log.Println(url)
	json.FromUrl(url, r)
	return nil
}

type RottenTomatoMovies struct {
	Total int

	Movies []*RottenTomatoMovie

	Links struct {
		Self string
		Next string
		Prev string
	}

	Error string
}

type movie_match struct {
	title string
	year  int
}

func (m *movie_match) scoreMatch(movie *RottenTomatoMovie) (int, float32) {
	titledist := strtools.LevenshteinDistanceCustom(m.title, movie.Title,
		1.0, 1.0, 0.3, strtools.EqualRuneFold)

	diffyear := intmath.Abs(m.year - movie.Year)
	shortlength := intmath.Min(len(m.title), len(movie.Title))
	tolerance := float32(shortlength) * 0.1
	//log.Println(titledist, diffyear, shortlength, tolerance)

	if titledist == 0.0 && (diffyear == 0) {
		return 4, titledist
	} else if (diffyear == 0 && m.year != 0) || titledist < 0.5 {
		return 2, titledist
	} else if titledist < tolerance && (m.year == 0 || diffyear < 2) {
		return 1, titledist
	}
	return 0, titledist
}

type scoredMovie struct {
	score int
	dist  float32
	movie *RottenTomatoMovie
}

type scoredMovieSorter struct {
	movies []scoredMovie
}

func (s *scoredMovieSorter) Len() int { return len(s.movies) }
func (s *scoredMovieSorter) Swap(i, j int) {
	s.movies[i], s.movies[j] = s.movies[j], s.movies[i]
}
func (s *scoredMovieSorter) Less(i, j int) bool {
	if s.movies[i].score != s.movies[j].score {
		return s.movies[i].score > s.movies[j].score
	} else if s.movies[i].dist != s.movies[j].dist {
		return s.movies[i].dist < s.movies[j].dist
	} else {
		return s.movies[i].movie.Year > s.movies[j].movie.Year
	}
}

func selectBest(matches []scoredMovie) []*RottenTomatoMovie {
	var result []*RottenTomatoMovie
	if len(matches) == 1 {
		result = append(result, matches[0].movie)
	} else if len(matches) > 0 {
		s := scoredMovieSorter{movies: matches}
		sort.Sort(&s)

		cutoffdist := matches[0].dist * float32(5-matches[0].score)
		cutoffscore := matches[0].score - 1

		for i, m := range matches {
			if (m.dist > cutoffdist || m.score < cutoffscore) && len(result) < 3 {
				// break if results are getting bad, but only as long as there are
				// only very few great ones.
				break
			}
			//log.Println("->", i, m.dist, m.movie.Title, m.movie.Year)
			result = append(result, matches[i].movie)
		}
	}
	return result
}

func (c *Config) FindMovie(title string, year int) ([]*RottenTomatoMovie, error) {
	q := c.LinkUrl.Query()
	q.Set("q", title)

	query := *c.LinkUrl
	query.RawQuery = q.Encode()
	rawquery := query.String()
	log.Println(rawquery)

	var matches []scoredMovie
	var batch_answers RottenTomatoMovies

	mm := movie_match{title, year}

	for len(matches) < 1 && rawquery != "" {
		fmt.Println(rawquery)
		err := json.FromUrl(rawquery, &batch_answers)
		if err != nil {
			return selectBest(matches), err
		}
		if batch_answers.Error != "" {
			return selectBest(matches), errors.New(batch_answers.Error)
		}
		for _, m := range batch_answers.Movies {
			var s scoredMovie
			s.score, s.dist = mm.scoreMatch(m)
			//log.Println(s.score, s.dist, m.Title, m.Year)
			s.movie = m
			matches = append(matches, s)
		}
		rawquery = batch_answers.Links.Next
	}
	return selectBest(matches), nil
}
