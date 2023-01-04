package main

type SocketPayload struct {
	Message string `json:"message"`
	Sender  string `json:"sender"`
}
