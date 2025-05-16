package decomposition_tools

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type MetaTournamentDate struct {
	DateStart time.Time
	DateEnd   time.Time
	Year      int
}

func GetMetaMatchDate(date string) *MetaTournamentDate {
	m := &MetaTournamentDate{}
	qs := strings.Split(date, ",")

	var year int
	numYear, err := strconv.Atoi(strings.TrimPrefix(qs[1], " "))
	if err == nil {
		year = numYear
	}

	if qs[0] != "" {
		ms := strings.Split(qs[0], " ")
		m.Year = year
		de, err := time.Parse("2/January/2006", fmt.Sprintf("%s/%s/%d", ms[1], ms[0], year))
		if err != nil {
			return nil
		}
		m.DateStart = de
		m.DateEnd = de

		return m
	}

	return nil
}

func GetMetaMatchDateWithTime(date string) *MetaTournamentDate {
	m := &MetaTournamentDate{}
	fs := strings.Split(date, " - ")

	if len(fs) != 2 {
		return GetMetaMatchDate(fs[0])
	}
	qs := strings.Split(fs[0], ",")

	var year int
	numYear, err := strconv.Atoi(strings.TrimPrefix(qs[1], " "))
	if err == nil {
		year = numYear
	}

	if qs[0] != "" {
		ms := strings.Split(qs[0], " ")
		m.Year = year

		de, err := time.Parse("2/January/2006", fmt.Sprintf("%s/%s/%d", ms[1], ms[0], year))
		if err != nil {
			return m
		}

		ts := strings.Split(fs[1], ":")
		hour, err := strconv.Atoi(ts[0])
		if err != nil {
			return m
		}
		minute, err := strconv.Atoi(strings.Replace(ts[1], " ", "", -1))
		if err != nil {
			return m
		}
		nde := de.Add(time.Hour*time.Duration(hour) + time.Minute*time.Duration(minute))

		m.DateStart = nde
		m.DateEnd = nde

		return m
	}

	return nil
}

func GetMetaTournamentDate(date string) *MetaTournamentDate {
	m := &MetaTournamentDate{}
	qs := strings.Split(date, ",")

	var year int
	numYear, err := strconv.Atoi(strings.TrimPrefix(qs[1], " "))
	if err == nil {
		year = numYear
	}
	ds := strings.Split(qs[0], "-")

	if len(ds) == 1 {
		ms := strings.Split(ds[0], " ")
		if len(ms) == 1 {
			return m
		}
		m.Year = year
		de, err := time.Parse("2/Jan/2006", fmt.Sprintf("%s/%s/%d", ms[1], ms[0], year))
		if err != nil {
			return nil
		}
		m.DateStart = de
		m.DateEnd = de

		return m
	}

	ms := strings.Split(ds[0], " ")
	ds2 := strings.Split(strings.TrimPrefix(ds[1], " "), " ")

	if len(ds2) > 1 {
		de, err := time.Parse("2/Jan/2006", fmt.Sprintf("%s/%s/%d", ds2[1], ds2[0], year))
		if err != nil {
			return nil
		}
		m.DateEnd = de
	} else {
		de, err := time.Parse("2/Jan/2006", fmt.Sprintf("%s/%s/%d", ds2[0], ms[0], year))
		if err != nil {
			return nil
		}
		m.DateEnd = de
	}

	dst, err := time.Parse("2/Jan/2006", fmt.Sprintf("%s/%s/%d", ms[1], ms[0], year))
	if err != nil {
		return nil
	}
	m.DateStart = dst
	m.Year = year

	return m
}

func BatchScore(s string) string {
	b := strings.Split(s, "-")

	return fmt.Sprintf("%s:%s", string([]rune(b[2])[1:]), string([]rune(b[3])[1:]))
}
