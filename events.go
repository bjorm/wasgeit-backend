package wasgeit

import (
	"time"
)

// Event describes an event taking place in a Venue
type Event struct {
	ID       int64
	Title    string
	DateTime time.Time
	Created  time.Time
	URL      string
	Venue    Venue
}

type Update struct {
	ExistingEv    Event
	UpdatedEv     Event
	ChangedFields []string
}

type ChangeSet struct {
	New     []Event
	Updates []Update
}

// TODO query DB directly instead of loading all events?
func DedupeAndTrackChanges(existingEvents []Event, newEvents []Event, cr Crawler) ChangeSet {
	var cs ChangeSet
	uniquenessVotes := 0

	for _, newEv := range newEvents {
		for _, existingEv := range existingEvents {
			if cr.IsSame(newEv, existingEv) {
				if hasDiff, update := diff(newEv, existingEv); hasDiff {
					cs.Updates = append(cs.Updates, update)
				}
			} else {
				uniquenessVotes++
			}
		}
		if uniquenessVotes == len(existingEvents) {
			cs.New = append(cs.New, newEv)
		}
		uniquenessVotes = 0
	}

	return cs
}

func diff(newEv Event, existingEv Event) (bool, Update) {
	sameTitle := newEv.Title == existingEv.Title
	sameTime := newEv.DateTime.Equal(existingEv.DateTime)

	if sameTitle && sameTime {
		return false, Update{}
	}

	update := Update{ExistingEv: existingEv, UpdatedEv: newEv}
	if !sameTime {
		update.ChangedFields = append(update.ChangedFields, "date")
	}
	if !sameTitle {
		update.ChangedFields = append(update.ChangedFields, "title")
	}

	return true, update
}
