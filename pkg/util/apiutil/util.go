package apiutil

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/tinyzimmer/kvdi/pkg/apis/kvdi/v1alpha1"
	v1 "github.com/tinyzimmer/kvdi/pkg/apis/meta/v1"
	"github.com/tinyzimmer/kvdi/pkg/util/errors"
)

// WriteOrLogError will write the provided content to the response writer, or
// log any error. It assumes the content is valid JSON.
func WriteOrLogError(out []byte, w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	if _, err := w.Write(append(out, []byte("\n")...)); err != nil {
		fmt.Println("Failed to write API response:", string(out), "error", err)
	}
}

// ReturnAPIErrors returns a BadRequest status code with a json encoded list
// of errors.
func ReturnAPIErrors(errs []error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	out := make([]string, 0)
	for _, err := range errs {
		out = append(out, err.Error())
	}
	jout, err := json.Marshal(map[string][]string{
		"errors": out,
	})
	if err != nil {
		fmt.Println("Failed to marshal errors to json:", err)
		jout = []byte(`{"error": "Multiple errors happened while processing the request"}`)
	}
	WriteOrLogError(jout, w)
}

// ReturnAPIError returns a BadRequest status code with a json encoded error
// message.
func ReturnAPIError(err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	WriteOrLogError(errors.ToAPIError(err).JSON(), w)
}

// ReturnAPINotFound returns a NotFound status code with a json encoded error
// message.
func ReturnAPINotFound(err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	WriteOrLogError(errors.ToAPIError(err).JSON(), w)
}

// ReturnAPIForbidden returns a Forbidden status code with a json encoded error
// message. If the denial happened due to an error, it logs the error server side.
func ReturnAPIForbidden(err error, msg string, w http.ResponseWriter) {
	if err != nil {
		fmt.Println("Forbidden request due to:", err.Error())
	}
	w.WriteHeader(http.StatusForbidden)
	WriteOrLogError(errors.ToAPIError(fmt.Errorf("Forbidden: %s", msg)).JSON(), w)
}

// WriteJSON encodes the provided interface to JSON and writes it to the response
// stream.
func WriteJSON(i interface{}, w http.ResponseWriter) {
	out, err := json.MarshalIndent(i, "", "    ")
	if err != nil {
		ReturnAPIError(err, w)
		return
	}
	WriteOrLogError(out, w)
}

// UnmarshalRequest will read the body of the given request and decode it into
// the given interface.
func UnmarshalRequest(r *http.Request, in interface{}) error {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, in)
}

// WriteOK write a simple boolean okay response.
func WriteOK(w http.ResponseWriter) {
	WriteJSON(map[string]bool{
		"ok": true,
	}, w)
}

// FilterUserRolesByNames returns a list of UserRoles matching the provided names
// and clusterw
func FilterUserRolesByNames(roles []v1alpha1.VDIRole, names []string) []*v1.VDIUserRole {
	userRoles := make([]*v1.VDIUserRole, 0)
	for _, name := range names {
		for _, role := range roles {
			if role.GetName() == name {
				userRoles = append(userRoles, role.ToUserRole())
			}
		}
	}
	return userRoles
}
