package controllers

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/buro9/microcosm/web/api"
	"github.com/buro9/microcosm/web/bag"
	"github.com/buro9/microcosm/web/errors"
	"github.com/buro9/microcosm/web/templates/funcs"
	"github.com/go-playground/form/v4"
	"github.com/go-playground/validator/v10"
)

// CommentCreate will receive a POST containing a comment and create it.
//
// Note that the comment can be created against any item that permission exists
// to create it against. Meaning that the comment can be attached to an existing
// Conversation, Event, Huddle, etc that the user has permission to comment on.
// A comment can also be created against their own profile, and this then becomes
// the profile page.
//
// A comment here is just the raw HTML part and metadata about what it is attached
// to, i.e. attached to a conversation with id = 47
//
// Attachments are handled separately in that they are uploaded async ahead of the
// comment itself.
func CommentCreate(w http.ResponseWriter, r *http.Request) {
	// Verify user is signed in
	user := bag.GetProfile(r.Context())
	if user == nil {
		errors.Render(w, r, http.StatusForbidden, fmt.Errorf(`Need to be signed in to create comments`))
		return
	}

	// r.ParseForm() translates all form values into r.Form url.Values
	if err := r.ParseForm(); err != nil {
		errors.Render(w, r, http.StatusInternalServerError, err)
		return
	}

	// Decoder translates r.Form into the struct
	decoder = form.NewDecoder()
	var commentForm api.CommentForm
	if err := decoder.Decode(&commentForm, r.Form); err != nil {
		errors.Render(w, r, http.StatusBadRequest, err)
		return
	}

	// Validate applies the struct validation
	validate = validator.New()
	if err := validate.Struct(commentForm); err != nil {
		// We could format this better, so that the UI could do something...
		// at the moment it's just cold hard validation errors sent to the
		// browser.
		//
		// We can get hold of the errors like this:
		// validationErrors := err.(validator.ValidationErrors)
		errors.Render(w, r, http.StatusBadRequest, err)
		return
	}

	comment, status, err := commentForm.Post(r.Context())
	if err != nil {
		errors.Render(w, r, status, err)
		return
	}

	location := funcs.Api2ui(
		funcs.LinkFromLinks(comment.Meta.Links, "commentPage").Href,
	)

	u, err := url.Parse(location)
	if err != nil {
		errors.Render(w, r, http.StatusInternalServerError, fmt.Errorf(`Couldn't parse %s as a URL`, location))
		return
	}
	u.Fragment = fmt.Sprintf(`comment%d`, comment.ID)

	// Redirect to location.
	w.Header().Set("Location", u.String())
	w.WriteHeader(http.StatusSeeOther)
}
