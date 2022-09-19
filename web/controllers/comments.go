package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type FormValues struct {
	Markdown  string
	ID        string
	InReplyTo string
	ItemId    string
	ItemType  string
}

// LogoutPost will remove the session cookie, thus logging the user out
func CommentsPost(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	values := new(FormValues)
	values.Markdown = r.Form.Get("markdown")
	values.ID = r.Form.Get("id")
	values.InReplyTo = r.Form.Get("inReplyTo")
	values.ItemId = r.Form.Get("itemId")
	values.ItemType = r.Form.Get("itemType")

	b, err := json.Marshal(values)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(
		[]byte(
			fmt.Sprintf(
				"Placeholder for post form handler, form values: %s",
				b,
			),
		),
	)
}
