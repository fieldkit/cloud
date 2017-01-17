package webserver

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/O-C-R/auth/id"
	"github.com/gorilla/mux"

	"github.com/O-C-R/fieldkit/backend"
	"github.com/O-C-R/fieldkit/config"
	"github.com/O-C-R/fieldkit/data"
	"github.com/O-C-R/fieldkit/data/jsondocument"
	"github.com/O-C-R/fieldkit/data/message"
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

		tokenID := id.ID{}
		if err := tokenID.UnmarshalText([]byte(req.FormValue("token"))); err != nil {
			Error(w, err, 404)
			return
		}

		_, err := c.Backend.AuthTokenByInputIDAndID(inputID, tokenID)
		if err != nil {
			Error(w, err, 500)
			return
		}

		if err == backend.NotFoundError {
			Error(w, err, 404)
			return
		}

		var reader io.Reader
		switch vars["source"] {
		case "direct":
			reader = req.Body
		case "rockblock":
			req.FormValue("data")

			data, err := hex.DecodeString(req.FormValue("data"))
			if err != nil {
				Error(w, err, 400)
				return
			}

			reader = bytes.NewReader(data)
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

		request, err := data.NewRequest(inputID)
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

		messageReader := message.NewFieldkitBinaryReader(bytes.NewReader(requestData))
		messageDocument, err := messageReader.ReadMessage()
		if err != nil {
			Error(w, err, 500)
			return
		}

		messageMessage, err := data.NewMessage(request.ID, inputID)
		if err != nil {
			Error(w, err, 500)
			return
		}

		messageMessage.Data = messageDocument
		if err := c.Backend.AddMessage(messageMessage); err != nil {
			Error(w, err, 500)
			return
		}

		documentData := jsondocument.Null()

		// time, lat, lon, alt, temp, humidity, battery, uptime
		documentDataTime, err := messageDocument.ArrayIndex(1)
		if err != nil {
			Error(w, err, 500)
			return
		}

		documentData.SetDocument("/date", documentDataTime)

		documentDataLat, err := messageDocument.ArrayIndex(2)
		if err != nil {
			Error(w, err, 500)
			return
		}

		documentDataLon, err := messageDocument.ArrayIndex(3)
		if err != nil {
			Error(w, err, 500)
			return
		}

		documentData.SetDocument("/type", jsondocument.String("Feature"))
		documentData.SetDocument("/geometry/type", jsondocument.String("Point"))
		documentData.SetDocument("/geometry/coordinates/0", documentDataLat)
		documentData.SetDocument("/geometry/coordinates/1", documentDataLon)

		document, err := data.NewDocument(messageMessage.ID, request.ID, inputID)
		if err != nil {
			Error(w, err, 500)
			return
		}

		document.Data = documentData
		if err := c.Backend.AddDocument(document); err != nil {
			Error(w, err, 500)
			return
		}

		w.Header().Set("cache-control", "no-store")
		WriteJSON(w, documentData)
	})
}

func DocumentsHandler(c *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		documents, err := c.Backend.DocumentsByProjectSlugAndExpeditionSlug(vars["project"], vars["expedition"])
		if err != nil {
			Error(w, err, 500)
			return
		}

		WriteJSON(w, documents)
	})
}

func DocumentDataHandler(c *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		data, err := c.Backend.DocumentDataByProjectSlugAndExpeditionSlug(vars["project"], vars["expedition"])
		if err != nil {
			Error(w, err, 500)
			return
		}

		WriteJSON(w, data)
	})
}
