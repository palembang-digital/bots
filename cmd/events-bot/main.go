package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/palembang-digital/bots/pkg/configs"
	"github.com/palembang-digital/bots/pkg/telegram"
	"github.com/palembang-digital/website/pkg/models"
)

var (
	todayEventMesagge = `Event Palembang Digital Hari Ini 🥳
{{range .Events}}
<b><a href="{{.RegistrationURL}}">{{.Title}}</a></b>
Jam {{scheduledStartTime .ScheduledStart}}
{{end}}
Yuk join! 🤗`
)

func main() {
	var cfg = configs.Get()

	// Get events from Patal API
	log.Println("Getting all events from Patal API")
	allEvents, err := getPatalEvents()
	if err != nil {
		log.Panic(err)
	}
	log.Printf("All Patal events: %+v", allEvents)

	todayEvents := []models.Event{}
	for _, event := range allEvents {
		if dateEqual(time.Now(), *event.ScheduledStart) {
			todayEvents = append(todayEvents, event)
		}
	}
	log.Printf("Today events: %+v", todayEvents)

	if len(todayEvents) > 0 {
		log.Println("Authenticating Telegram bot...")
		bot, err := telegram.New(cfg.TelegramToken, cfg.TelegramChatID, cfg.TelegramDebug)
		if err != nil {
			log.Panic(err)
		}

		tmpl, err := template.New("today-events").Funcs(template.FuncMap{
			"scheduledStartTime": func(scheduledStart time.Time) string {
				loc, err := time.LoadLocation("Asia/Jakarta")
				if err != nil {
					return scheduledStart.Format("15:04") + " GMT+00"
				}
				return scheduledStart.In(loc).Format("15:04") + " WIB"
			},
		}).
			Parse(todayEventMesagge)
		if err != nil {
			log.Panic(err)
		}

		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, map[string]interface{}{
			"Date":   time.Now().Format("02/01/2006"),
			"Events": todayEvents,
		}); err != nil {
			log.Panic(err)
		}

		message := buf.String()

		log.Println("Send message", message)
		if err := bot.Send(cfg.TelegramChatID, message); err != nil {
			log.Panic(err)
		}
	}

	log.Println("Job completed. KTHXBYE! 👋")
}

// getPatalEvents gets the list of Patal Events from Patal API
func getPatalEvents() ([]models.Event, error) {
	resp, err := http.Get("https://palembangdigital.org/api/v1/events")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	events := []models.Event{}
	if err := json.Unmarshal(body, &events); err != nil {
		return nil, err
	}

	return events, nil
}

func dateEqual(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
