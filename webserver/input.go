package webserver

import (
	"crypto/sha1"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/O-C-R/auth/id"
	"github.com/gorilla/mux"

	"github.com/O-C-R/fieldkit/backend"
	"github.com/O-C-R/fieldkit/config"
	"github.com/O-C-R/fieldkit/data"
)

func InputsHandler(c *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		expeditions, err := c.Backend.InputsByProjectSlugAndExpeditionSlug(vars["project"], vars["expedition"])
		if err != nil {
			Error(w, err, 500)
			return
		}

		WriteJSON(w, expeditions)
	})
}

func InputHandler(c *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)

		inputID := id.ID{}
		if err := inputID.UnmarshalText([]byte(vars["id"])); err != nil {
			Error(w, err, 404)
			return
		}

		input, err := c.Backend.InputByID(inputID)
		if err == backend.NotFoundError {
			Error(w, err, 404)
			return
		}

		if err != nil {
			Error(w, err, 500)
			return
		}

		WriteJSON(w, input)
	})
}

func InputAddHandler(c *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		expedition, err := c.Backend.ExpeditionByProjectSlugAndSlug(vars["project"], vars["expedition"])
		if err == backend.NotFoundError {
			Error(w, err, 404)
			return
		}

		if err != nil {
			Error(w, err, 500)
			return
		}

		errs := data.NewErrors()

		// Validate the input name.
		name, slug := data.Name(req.FormValue("name"))
		if name == "" || slug == "" {
			errs.Error("name", "invalid Name")
		} else {
			slugInUse, err := c.Backend.InputSlugInUse(expedition.ID, slug)
			if err != nil {
				Error(w, err, 500)
				return
			}

			if slugInUse {
				errs.Error("name", "Name in use")
			}
		}

		if len(errs) > 0 {
			WriteJSONStatusCode(w, errs, 400)
			return
		}

		input, err := data.NewInput(expedition.ID, name, slug)
		if err != nil {
			Error(w, err, 500)
			return
		}

		if err := c.Backend.AddInput(input); err != nil {
			Error(w, err, 500)
			return
		}

		WriteJSON(w, input)
	})
}

func InputRequestHandler(c *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)

		inputID := id.ID{}
		if err := inputID.UnmarshalText([]byte(vars["id"])); err != nil {
			Error(w, err, 404)
			return
		}

		input, err := c.Backend.InputByID(inputID)
		if err == backend.NotFoundError {
			Error(w, err, 404)
			return
		}

		if err != nil {
			Error(w, err, 500)
			return
		}

		var reader io.Reader
		switch vars["source"] {
		case "direct":
			reader = req.Body
		}

		hash := sha1.New()
		reader = io.TeeReader(reader, hash)
		requestData, err := ioutil.ReadAll(reader)
		if err != nil {
			Error(w, err, 500)
			return
		}

		checksum := id.ID{}
		if err := checksum.UnmarshalBinary(hash.Sum(nil)); err != nil {
			Error(w, err, 500)
			return
		}

		request, err := data.NewRequest(input.ID)
		if err != nil {
			Error(w, err, 500)
			return
		}

		request.Format = vars["format"]
		request.Data = requestData
		request.Checksum = checksum
		if err := c.Backend.AddRequest(request); err != nil {
			Error(w, err, 500)
			return
		}

		WriteJSON(w, request)
	})
}
