package hospitalService

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type Hospital struct {
	Id      int    `json:"Id"`
	Name    string `json:"Name"`
	Email   string `json:"Email"`
	City    string `json:"City"`
	Contact string `json:"Contact"`
	Address string `json:"Address"`
}

type hospitalHandler struct {
	details map[string]Hospital
	slots   slotHandler
	booking bookingHandler
}

type Slot struct {
	Slot_id   int    `json: "Slot_id"`
	Hosp_id   int    `json: "Hosp_id"`
	Available int    `json: "Availabale"`
	Date      string `json: "Date"`
	TimeSlot  string `json: "TimeSlot"`
}

type slotHandler struct {
	sync.Mutex
	details map[string]Slot
}

type Booking struct {
	Booking_id int `json: "Booking_id"`
	Slot_id    int `json: "Slot_id"`
	User_id    int `json: "User_id"`
	Hosp_id    int `json: "Hosp_id"`
}

type bookingHandler struct {
	sync.Mutex
	details map[string]Booking
}

func NewHospitalHandler() *hospitalHandler {
	return &hospitalHandler{
		details: map[string]Hospital{
			"1": {
				Id:      1,
				Name:    "National Hospital",
				Email:   "national@gmail.com",
				City:    "Jabalpur",
				Contact: "8223821911",
				Address: "4th Bridge",
			},
			"2": {
				Id:      2,
				Name:    "City Hospital",
				Email:   "city@gmail.com",
				City:    "Pune",
				Contact: "9344821911",
				Address: "Hinjewadi",
			},

			"3": {
				Id:      3,
				Name:    "Apollo Hospital",
				Email:   "apollo@gmail.com",
				City:    "Mumbai",
				Contact: "9823492312",
				Address: "Marine Drives",
			},
		},
		slots: slotHandler{
			details: map[string]Slot{},
		},
		booking: bookingHandler{
			details: map[string]Booking{},
		},
	}
}

func (h *hospitalHandler) Default(rw http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.String(), "/")
	if len(parts) > 3 {
		id, err := strconv.Atoi(parts[2])
		if err != nil {
			panic(err)
		}
		if r.Method == "GET" {
			h.viewMySlots(rw, r, id)
			return
		}
		h.createSlot(rw, r)
		return
	}
	switch r.Method {
	case "GET":
		h.getSlots(rw, r)
	case "POST":
		h.bookSlot(rw, r)
	}
}

func (h *hospitalHandler) getSlots(rw http.ResponseWriter, r *http.Request) {
	slots := make([]Slot, 0)

	i := 0
	h.slots.Lock()
	for _, slot := range h.slots.details {
		if slot.Available > 0 {
			// slots[i] = slot
			slots = append(slots, slot)
			i++
		}
	}
	h.slots.Unlock()

	data, err := json.Marshal(slots)
	if err != nil {
		panic(err)
	}

	rw.Header().Add("content-type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(data)
}

func (h *hospitalHandler) viewMySlots(rw http.ResponseWriter, r *http.Request, id int) {
	bookings := make([]Booking, 0)

	i := 0
	h.booking.Lock()
	for _, booking := range h.booking.details {
		if booking.Hosp_id == id {
			// bookings[i] = booking
			bookings = append(bookings, booking)
			i++
		}
	}
	h.booking.Unlock()

	data, err := json.Marshal(bookings)
	if err != nil {
		panic(err)
	}

	rw.Header().Add("content-type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(data)
}

func (h *hospitalHandler) createSlot(rw http.ResponseWriter, r *http.Request) {

	if r.Body == http.NoBody {
		rw.WriteHeader(http.StatusBadRequest)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var slot Slot
	json.Unmarshal(body, &slot)

	id := len(h.slots.details) + 1

	slot.Slot_id = id

	h.slots.Lock()
	h.slots.details[fmt.Sprint(id)] = slot
	h.slots.Unlock()

	rw.WriteHeader(http.StatusCreated)
}

func (h *hospitalHandler) bookSlot(rw http.ResponseWriter, r *http.Request) {

	if r.Body == http.NoBody {
		rw.WriteHeader(http.StatusBadRequest)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var bookings Booking
	json.Unmarshal(body, &bookings)

	id := len(h.booking.details) + 1
	bookings.Booking_id = id

	slot := h.slots.details[fmt.Sprint(bookings.Slot_id)]
	slot.Available--

	h.booking.Lock()
	h.slots.details[fmt.Sprint(bookings.Slot_id)] = slot
	h.booking.details[fmt.Sprint(id)] = bookings
	h.booking.Unlock()

	rw.WriteHeader(http.StatusAccepted)
}
