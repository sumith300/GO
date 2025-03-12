package main

import (
	"fmt"
	"regexp"
	"sort"
	"time"
)

type Event struct {
	Name       string
	Date       time.Time
	StartTime  time.Time
	Duration   int
}

type eventStore struct {
	eventsMap map[string]*Event
	eventsList []*Event
}

var (
	dateRegex = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	timeRegex = regexp.MustCompile(`^([01]\d|2[0-3]):([0-5]\d)$`)
)

func (es *eventStore) addEvent(e *Event) error {
	if _, exists := es.eventsMap[e.Name]; exists {
		return fmt.Errorf("event with name '%s' already exists", e.Name)
	}
	es.eventsMap[e.Name] = e
	es.eventsList = append(es.eventsList, e)
	sort.Slice(es.eventsList, func(i, j int) bool {
		if es.eventsList[i].Date.Equal(es.eventsList[j].Date) {
			return es.eventsList[i].StartTime.Before(es.eventsList[j].StartTime)
		}
		return es.eventsList[i].Date.Before(es.eventsList[j].Date)
	})
	return nil
}

func (es *eventStore) hasTimeConflict(newEvent *Event) bool {
    newEnd := newEvent.StartTime.Add(time.Duration(newEvent.Duration) * time.Hour)
    
    for _, existing := range es.eventsList {
        if existing.Date.Equal(newEvent.Date) {
            existingEnd := existing.StartTime.Add(time.Duration(existing.Duration) * time.Hour)
            if newEvent.StartTime.Before(existingEnd) && newEnd.After(existing.StartTime) {
                return true
            }
        }
    }
    return false
}

func validateDateTime(dateStr, timeStr string) (time.Time, time.Time, error) {
	if !dateRegex.MatchString(dateStr) {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid date format, use YYYY-MM-DD")
	}
	if !timeRegex.MatchString(timeStr) {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid time format, use HH:MM 24h format")
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid date: %v", err)
	}

	tm, err := time.Parse("15:04", timeStr)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid time: %v", err)
	}

	combined := time.Date(date.Year(), date.Month(), date.Day(),
	    tm.Hour(), tm.Minute(), 0, 0, time.UTC)
    
	if combined.Before(time.Now().Add(-5 * time.Minute)) {
		return time.Time{}, time.Time{}, fmt.Errorf("date/time must be in the future")
	}

	return date, tm, nil
}

func main() {
	es := &eventStore{
		eventsMap: make(map[string]*Event),
	}

	// Sample menu loop
	for {
		fmt.Println("\nChronoSync Event Manager")
		fmt.Println("1. Add New Event")
		fmt.Println("2. Modify Event")
		fmt.Println("3. Delete Event")
		fmt.Println("4. List All Upcoming Events")
		fmt.Println("5. View Daily Events")
		fmt.Println("6. Display Analytics")
		fmt.Println("7. Exit")

		var choice int
		fmt.Print("Enter choice: ")
		choiceStr, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		choice, err := strconv.Atoi(strings.TrimSpace(choiceStr))
		if err != nil {
		    fmt.Println("Invalid input: Please enter a number between 1-7")
		    continue
		}

		switch choice {
		case 1:
			var name, dateStr, timeStr string
			var duration int

			fmt.Print("Enter event name: ")
			name, _ := bufio.NewReader(os.Stdin).ReadString('\n')
			name = strings.TrimSpace(name)

			fmt.Print("Enter date (YYYY-MM-DD): ")
			fmt.Scanln(&dateStr)

			fmt.Print("Enter start time (HH:MM): ")
			fmt.Scanln(&timeStr)

			fmt.Print("Enter duration in hours: ")
			durStr, _ := bufio.NewReader(os.Stdin).ReadString('\n')
			duration, err := strconv.Atoi(strings.TrimSpace(durStr))
			if err != nil || duration <= 0 {
			    fmt.Println("Invalid duration: must be positive integer")
			    continue
			}

			date, startTime, err := validateDateTime(dateStr, timeStr)
			if err != nil {
				fmt.Println("Validation error:", err)
				continue
			}

			event := &Event{
				Name:      name,
				Date:      date,
				StartTime: startTime,
				Duration:  duration,
			}

			if err := es.addEvent(event); err != nil {
				fmt.Println("Error adding event:", err)
			} else {
				fmt.Println("Event added successfully!")
			}
		case 2:
			var name, newDateStr, newTimeStr string
			var newDuration int

			fmt.Print("Enter event name to modify: ")
			modName, _ := bufio.NewReader(os.Stdin).ReadString('\n')
			modName = strings.TrimSpace(modName)

			event, exists := es.eventsMap[modName] 
			if !exists {
			    fmt.Printf("Error: Event '%s' not found\n", modName)
			    continue
			}

			fmt.Print("Enter new date (YYYY-MM-DD): ")
			newDateStr, _ := bufio.NewReader(os.Stdin).ReadString('\n')
			newDateStr = strings.TrimSpace(newDateStr)

			fmt.Print("Enter new time (HH:MM): ")
			newTimeStr, _ := bufio.NewReader(os.Stdin).ReadString('\n')
			newTimeStr = strings.TrimSpace(newTimeStr)

			fmt.Printf("Current date (%s) - enter new date (YYYY-MM-DD): ", event.Date.Format("2006-01-02"))
			fmt.Scanln(&newDateStr)

			fmt.Printf("Current time (%s) - enter new time (HH:MM): ", event.StartTime.Format("15:04"))
			fmt.Scanln(&newTimeStr)

			fmt.Printf("Current duration (%dh) - enter new duration: ", event.Duration)
			_, err := fmt.Scanln(&newDuration)
			if err != nil || newDuration <= 0 {
				fmt.Println("Invalid duration: must be positive integer")
				continue
			}

			newDate, newTime, err := validateDateTime(newDateStr, newTimeStr)
			if err != nil {
				fmt.Println("Validation error:", err)
				continue
			}

			// Update event details
			event.Date = newDate
			event.StartTime = newTime
			event.Duration = newDuration

			// Re-sort events list
			sort.Slice(es.eventsList, func(i, j int) bool {
				if es.eventsList[i].Date.Equal(es.eventsList[j].Date) {
					return es.eventsList[i].StartTime.Before(es.eventsList[j].StartTime)
				}
				return es.eventsList[i].Date.Before(es.eventsList[j].Date)
			})

			fmt.Println("Event modified successfully!")
		case 3:
			var name string
			fmt.Print("Enter event name to delete: ")
			delName, _ := bufio.NewReader(os.Stdin).ReadString('\n')
			delName = strings.TrimSpace(delName)

			if _, exists := es.eventsMap[delName]; !exists {
				fmt.Printf("Error: Event '%s' not found\n", delName)
				continue
			}

			// Remove from map
			delete(es.eventsMap, name)

			// Remove from slice
			for i, event := range es.eventsList {
				if event.Name == name {
					es.eventsList = append(es.eventsList[:i], es.eventsList[i+1:]...)
					break
				}
			}

			fmt.Println("Event deleted successfully!")
		case 6:
			return
		case 4:
			var dateStr string
			fmt.Print("Enter date to view events (YYYY-MM-DD): ")
			fmt.Scanln(&dateStr)

			if !dateRegex.MatchString(dateStr) {
				fmt.Println("Invalid date format")
				continue
			}

			date, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				fmt.Println("Invalid date:", err)
				continue
			}

			fmt.Printf("\nEvents on %s:\n", dateStr)
			for _, event := range es.eventsList {
				if event.Date.Equal(date) {
					fmt.Printf("- %s: %s (%d hours)\n", 
						event.Name, 
						event.StartTime.Format("15:04"),
						event.Duration)
				}
			}

		case 5:
			if len(es.eventsList) == 0 {
				fmt.Println("No events to analyze")
				continue
			}

			func (es *eventStore) longestEvent() *Event {
			    if len(es.eventsList) == 0 {
			        return nil
			    }
			    longest := es.eventsList[0]
			    for _, event := range es.eventsList {
			        if event.Duration > longest.Duration {
			            longest = event
			        }
			    }
			    return longest
			}
			
			func (es *eventStore) countEventsByDate(date time.Time) int {
			    count := 0
			    for _, event := range es.eventsList {
			        if event.Date.Equal(date) {
			            count++
			        }
			    }
			    return count
			}
			
			func (es *eventStore) averageDuration() float64 {
			    if len(es.eventsList) == 0 {
			        return 0
			    }
			    total := 0
			    for _, event := range es.eventsList {
			        total += event.Duration
			    }
			    return float64(total) / float64(len(es.eventsList))
			}

			fmt.Println("\nEvent Analytics:")
			
			// Longest event
			if longest := es.longestEvent(); longest != nil {
			    fmt.Printf("Longest Event: %s (%d hours)\n", longest.Name, longest.Duration)
			} else {
			    fmt.Println("No events to analyze")
			}
			
			// Event count by date
			var dateStr string
			fmt.Print("Enter date to count events (YYYY-MM-DD): ")
			fmt.Scanln(&dateStr)
			date, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
			    fmt.Println("Invalid date format")
			    break
			}
			count := es.countEventsByDate(date)
			fmt.Printf("Events on %s: %d\n", dateStr, count)
			
			// Average duration
			avg := es.averageDuration()
			fmt.Printf("Average Event Duration: %.2f hours\n", avg)

		default:
			fmt.Println("Invalid choice")
		}
	}
}