package main

import (
	"MiniProject/hospitalService"
	"MiniProject/userService"
	"net/http"
)

func main() {
	users := userService.NewUserHandler()
	hospitals := hospitalService.NewHospitalHandler()
	http.HandleFunc("/users/", users.Default)
	http.HandleFunc("/hospitals/", hospitals.Default)
	// http.HandleFunc("/hospitals/{id}/createSlot", hospitals.Default)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
