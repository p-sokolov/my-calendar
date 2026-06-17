package calendar

import (
	"strconv"
	"sync"
	"time"

	e "my-calendar/internal/error"
)

type Event struct {
	EventID     int    `json:"id"`
	CreatorID   int    `json:"user_id"`
	Date        string `json:"date"`
	Description string `json:"event"`
}

type Calendar struct {
	mu     sync.RWMutex
	events map[int][]Event
	nextID int
}

const dateFormat = "2006-01-02"

func NewCalendar() *Calendar {
	return &Calendar{
		events: make(map[int][]Event),
		nextID: 1,
	}
}

func (c *Calendar) GetEvents(creatorIdstr string, date string, period string) ([]Event, error) {
	switch period {
	case "":
		return c.GetAllPeriods(creatorIdstr, date)
	case "day":
		return c.GetDaily(creatorIdstr, date)
	case "week":
		return c.GetWeekly(creatorIdstr, date)
	case "month":
		return c.GetMonthly(creatorIdstr, date)
	default:
		return nil, e.ErrUnknownPeriod
	}
}

func (c *Calendar) GetDaily(creatorIdstr string, date string) ([]Event, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

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
			return nil, e.ErrUserIdTypeMismatch
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
	c.mu.RLock()
	defer c.mu.RUnlock()

	var weeklyEvents []Event

	hasCreatorId := creatorIdstr != ""
	hasDate := date != ""

	var creatorId int
	var err error

	if hasCreatorId {
		creatorId, err = strconv.Atoi(creatorIdstr)
		if err != nil {
			return nil, e.ErrUserIdTypeMismatch
		}
	}

	targetDate := time.Now()
	if hasDate {
		targetDate, err = time.Parse(dateFormat, date)
		if err != nil {
			return nil, e.ErrInvalidDateFormat
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
	c.mu.RLock()
	defer c.mu.RUnlock()

	var monthlyEvents []Event

	hasCreatorId := creatorIdstr != ""
	hasDate := date != ""

	var creatorId int
	var err error

	if hasCreatorId {
		creatorId, err = strconv.Atoi(creatorIdstr)
		if err != nil {
			return nil, e.ErrUserIdTypeMismatch
		}
	}

	targetDate := time.Now()
	if hasDate {
		targetDate, err = time.Parse(dateFormat, date)
		if err != nil {
			return nil, e.ErrInvalidDateFormat
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
		for _, event := range c.events {
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

func (c *Calendar) GetAllPeriods(creatorIdstr string, date string) ([]Event, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var allEvents []Event
	var creatorId int
	var err error

	hasCreatorId := creatorIdstr != ""
	hasDate := date != ""

	if hasCreatorId {
		creatorId, err = strconv.Atoi(creatorIdstr)
		if err != nil {
			return nil, e.ErrUserIdTypeMismatch
		}
	}

	switch {
	case hasCreatorId && hasDate:
		for _, e := range c.events[creatorId] {
			if e.Date == date {
				allEvents = append(allEvents, e)
			}
		}
	case hasCreatorId:
		candidates := c.events[creatorId]
		allEvents = append(allEvents, candidates...)
	case hasDate:
		for _, event := range c.events {
			for _, e := range event {
				if e.Date == date {
					allEvents = append(allEvents, e)
				}
			}
		}
	default:
		for _, event := range c.events {
			allEvents = append(allEvents, event...)
		}
	}

	return allEvents, nil
}

func (c *Calendar) CreateEvent(event Event) (Event, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if event.CreatorID <= 0 {
		return Event{}, e.ErrInvalidUserId
	}

	if _, err := time.Parse(dateFormat, event.Date); err != nil {
		return Event{}, e.ErrInvalidDate
	}

	event.EventID = c.nextID
	c.nextID++

	c.events[event.CreatorID] = append(c.events[event.CreatorID], event)
	return event, nil
}

func (c *Calendar) UpdateEvent(event Event) (Event, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, err := time.Parse(dateFormat, event.Date); err != nil {
		return Event{}, e.ErrInvalidDate
	}

	for userID, events := range c.events {
		for i := range events {
			if events[i].EventID == event.EventID {
				old := c.events[userID][i]
				old.Date = event.Date
				old.Description = event.Description

				c.events[userID][i] = old
				return event, nil
			}
		}
	}

	return Event{}, e.ErrEventNotFound
}

func (c *Calendar) DeleteEvent(id int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for userID, events := range c.events {
		for i := range events {
			if events[i].EventID == id {
				c.events[userID] = append(events[:i], events[i+1:]...)
				return nil
			}
		}
	}

	return e.ErrEventNotFound
}
