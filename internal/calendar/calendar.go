package calendar

import (
	"time"
	"strconv"
)

type Event struct {
	EventID 	int
	CreatorID 	int
	Date 		string
	Description string
}

type Calendar struct {
	events map[int]Event
}

func NewCalendar() *Calendar {
	return &Calendar{ events: make(map[int]Event)}
}

func (c *Calendar) GetDaily(creatorId string, date string) ([]Event, error) {	
	var dailyEvents []Event

	creatorId_, err := strconv.Atoi(creatorId)
	if err != nil {        
        return nil, err
    }
    
	for _, e := range c.events {
		if date == e.Date || creatorId == e.CreatorID {
			dailyEvents = append(dailyEvents, e)
		}
	}

	return dailyEvents
}

func (c *Calendar) GetMonthly(creatorId int, month time.Month, year int) []Event {	
	var weeklyEvents []Event
	
	for _, e := range c.events {

		curDate, _ := time.Parse("2001-01-01", e.Date)
		curYear, curMonth, _ := curDate.Date()
		
		if (curYear == year && curMonth == month) || creatorId == e.CreatorID {
			weeklyEvents = append(weeklyEvents, e)
		}
	}

	return weeklyEvents
}
