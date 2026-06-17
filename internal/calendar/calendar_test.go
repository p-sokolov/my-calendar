package calendar

import "testing"

func TestCreateEvent(t *testing.T) {
	c := NewCalendar()

	event := Event{
		CreatorID:  1,
		Date:       "2026-06-08",
		Description: "Daily standup",
	}

	created, err := c.CreateEvent(event)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if created.EventID != 1 {
		t.Errorf("expected id=1, got %d", created.EventID)
	}

	if len(c.events[1]) != 1 {
		t.Errorf("expected 1 event, got %d", len(c.events[1]))
	}
}

func TestCreateEvent_InvalidDate(t *testing.T) {
	c := NewCalendar()

	event := Event{
		CreatorID: 1,
		Date:      "2026506-08",
	}

	_, err := c.CreateEvent(event)

	if err == nil {
		t.Fatal("expected error")
	}

	if err.Error() != "invalid date" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteEvent(t *testing.T) {
	c := NewCalendar()

	event, _ := c.CreateEvent(Event{
		CreatorID: 1,
		Date:      "2026-06-08",
	})

	err := c.DeleteEvent(event.EventID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(c.events[1]) != 0 {
		t.Fatal("event was not deleted")
	}
}

func TestCreateEventValidation(t *testing.T) {
	tests := []struct {
		name      string
		event     Event
		wantError bool
	}{
		{
			name: "valid",
			event: Event{
				CreatorID: 1,
				Date:      "2026-06-08",
			},
		},
		{
			name: "invalid user",
			event: Event{
				CreatorID: 0,
				Date:      "2026-06-08",
			},
			wantError: true,
		},
		{
			name: "invalid date",
			event: Event{
				CreatorID: 1,
				Date:      "abc",
			},
			wantError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := NewCalendar()

			_, err := c.CreateEvent(tc.event)

			if tc.wantError && err == nil {
				t.Fatal("expected error")
			}

			if !tc.wantError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}