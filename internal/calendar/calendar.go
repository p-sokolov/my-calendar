package calendar

import (
	"fmt"
	"strconv"
	"time"
)

type Event struct {
	// EventID 	int		`json:"-"`
	CreatorID 	int		`json:"user_id"`
	Date 		string 	`json:"date"`
	Description string 	`json:"event"`
}

type Calendar struct {
	events map[int][]Event
}

const dateFormat = "2006-01-02"

func NewCalendar() *Calendar {
	return &Calendar{ events: make(map[int][]Event)}
}

func (c *Calendar) GetDaily(creatorIdstr string, date string) ([]Event, error) {

	// if creatorId in null and date is null then we should search all daily event based on time.Now
	// if creatorId!=nil and date is null then we got map[creatorId] and thats it
	// if creatorId is null and date!=nil then we should shrink all event by creatorId and filter by time.Now (today)
	// if we know all params then do map[creatorId] and filter it by trget date
	
	var dailyEvents []Event
	var creatorId int
	var err error
	
	hasCreatorId := creatorIdstr != ""
	hasDate := date != ""

	if hasCreatorId {
		creatorId, err = strconv.Atoi(creatorIdstr)
		if err != nil {        
	        return nil, fmt.Errorf("user_id must be number")
	    }
	}	

	today := time.Now().Format(dateFormat)
	
	switch {
	case hasCreatorId && hasDate:
		for _, e := range c.events[creatorId] {
			if e.Date == date {
				dailyEvents = append(dailyEvents, e)
			}
		}
	case hasCreatorId:
		candidates := c.events[creatorId]
		for _, e := range candidates {
			if e.Date == today {
				dailyEvents = append(dailyEvents, e)
			}
		}
	case hasDate:
		for _, event := range c.events {
			for _, e := range event {
				if e.Date == date {
					dailyEvents = append(dailyEvents, e)
				}
			}
		}
	default:
		for _, event := range c.events {
			for _, e := range event {
				if e.Date == today {
					dailyEvents = append(dailyEvents, e)
				}
			}
		}
	}
	
	return dailyEvents, nil
}

func (c *Calendar) GetWeekly(creatorIdstr string, date string) ([]Event, error) {
	var weeklyEvents []Event
	
	hasCreatorId := creatorIdstr != ""
	hasDate := date != ""

	var creatorId int
	var err error

	if hasCreatorId {
		creatorId, err = strconv.Atoi(creatorIdstr)
		if err != nil {        
	        return nil, fmt.Errorf("user_id must be number")
	    }
	}

	targetDate := time.Now()
	if hasDate {
		targetDate, err = time.Parse(dateFormat, date)
		if err != nil {
        	return nil, fmt.Errorf("date must be in YYYY-MM-DD format")
		}
	}

	weekday := int(targetDate.Weekday())
	if weekday == 0 {
	    weekday = 7
	}
	
	startOfWeek := time.Date(
		targetDate.Year(),
		targetDate.Month(),
		targetDate.Day()-(weekday-1),
		0, 0, 0, 0,
		targetDate.Location(),
	)
	endOfWeek := startOfWeek.AddDate(0, 0, 7)

	switch {
	case hasCreatorId:
		for _, e := range c.events[creatorId] {
			curDate, _ := time.Parse(dateFormat, e.Date)

			if !curDate.Before(startOfWeek) && curDate.Before(endOfWeek) {
				weeklyEvents = append(weeklyEvents, e)
			}
		}

	default:
		for _, events := range c.events {
			for _, e := range events {
				curDate, _ := time.Parse(dateFormat, e.Date)

				if !curDate.Before(startOfWeek) && curDate.Before(endOfWeek) {
					weeklyEvents = append(weeklyEvents, e)
				}
			}
		}
	}
	
	return weeklyEvents, nil
}

func (c *Calendar) GetMonthly(creatorIdstr string, date string) ([]Event, error) {
	var monthlyEvents []Event

	hasCreatorId := creatorIdstr != ""
	hasDate := date != ""

	var creatorId int
	var err error

	if hasCreatorId {
		creatorId, err = strconv.Atoi(creatorIdstr)
		if err != nil {        
	        return nil, fmt.Errorf("user_id must be number")
	    }
	}

	targetDate := time.Now()
	if hasDate {
		targetDate, err = time.Parse(dateFormat, date)
		if err != nil {
        	return nil, fmt.Errorf("date must be in YYYY-MM-DD format")
		}
	}

	tagetYear, targetMonth, _ := targetDate.Date()

	switch {
	case hasCreatorId:
		for _, e := range c.events[creatorId] {
			curDate, _ := time.Parse(dateFormat, e.Date)
			curYear, curMonth, _ := curDate.Date()

			if curYear == tagetYear && curMonth == targetMonth {
				monthlyEvents = append(monthlyEvents, e)
			}
		}
	default:
		for _, event := range c.events{
			for _, e := range event {				
				curDate, _ := time.Parse(dateFormat, e.Date)
				curYear, curMonth, _ := curDate.Date()

				if curYear == tagetYear && curMonth == targetMonth {
					monthlyEvents = append(monthlyEvents, e)
				}
			}
		}		
	}

	return monthlyEvents, nil
}

func (c *Calendar) CreateEvent(event Event) error {	
	if event.CreatorID <= 0 {
		return fmt.Errorf("invalid user_id]")
    }

    if _, err := time.Parse(dateFormat, event.Date); err != nil {
		return fmt.Errorf("invalid date")
    }

    c.events[event.CreatorID] = append(c.events[event.CreatorID], event)  
    return nil
}