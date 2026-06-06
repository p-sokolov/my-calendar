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

	// if creatorId in null and date is null then we should search all daily event based on time.Now
	// if creatorId!=nil and date is null then we got map[creatorId] and thats it
	// if creatorId is null and date!=nil then we should shrink all event by time.Now (today) and filter by target date
	// if we know all params then do map[creatorId] and filter it by trget date
	
	var dailyEvents []Event

	if creatorId != "" {
		creatorId_, err := strconv.Atoi(creatorId)
		if err != nil {        
	        return nil, err
	    }
	}    
	for _, e := range c.events {
		if date == e.Date || creatorId == e.CreatorID {
			dailyEvents = append(dailyEvents, e)
		}
	}

	return dailyEvents, nil
}

func (c *Calendar) GetMonthly(creatorId string, date string) ([]Event, error) {
	
	var weeklyEvents []Event
	targetYear, targetMonth, _ := date.Date()

   	if creatorId != "" {
		creatorId_, err := strconv.Atoi(creatorId)
		if err != nil {        
	        return nil, err
	    }
	}   
	for _, e := range c.events {
		curDate, _ := time.Parse("2001-01-01", e.Date)
		curYear, curMonth, _ := curDate.Date()
		
		if (curYear == targetYear && curMonth == targetMonth) || creatorId_ == e.CreatorID {
			weeklyEvents = append(weeklyEvents, e)
		}
	}

	return weeklyEvents
}
