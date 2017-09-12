package content

import (
	"fmt"
	"net/http"

	"github.com/bosssauce/reference"
	"github.com/ponzu-cms/ponzu/management/editor"
	"github.com/ponzu-cms/ponzu/system/item"
)

type Bookings struct {
	item.Item

	FullName       string `json:"full_name"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	RegistrationId string `json:"registration_id"`
	ModeOfPayment  string `json:"mode_of_payment"`
	Room           string `json:"room"`
	Hotel          string `json:"hotel"`
}

// MarshalEditor writes a buffer of html to edit a Bookings within the CMS
// and implements editor.Editable
func (b *Bookings) MarshalEditor() ([]byte, error) {
	view, err := editor.Form(b,
		// Take note that the first argument to these Input-like functions
		// is the string version of each Bookings field, and must follow
		// this pattern for auto-decoding and auto-encoding reasons:
		editor.Field{
			View: editor.Input("FullName", b, map[string]string{
				"label":       "Full Name",
				"type":        "text",
				"placeholder": "Enter the Full Name here",
			}),
		},
		editor.Field{
			View: editor.Input("Email", b, map[string]string{
				"label":       "Email",
				"type":        "text",
				"placeholder": "Enter the Email here",
			}),
		},
		editor.Field{
			View: editor.Input("Phone", b, map[string]string{
				"label":       "Phone",
				"type":        "text",
				"placeholder": "Enter the Phone here",
			}),
		},
		editor.Field{
			View: editor.Input("RegistrationId", b, map[string]string{
				"label":       "RegistrationId",
				"type":        "text",
				"placeholder": "Enter the RegistrationId here",
			}),
		},
		editor.Field{
			View: editor.Input("ModeOfPayment", b, map[string]string{
				"label":       "ModeOfPayment",
				"type":        "text",
				"placeholder": "Enter the ModeOfPayment here",
			}),
		},
		editor.Field{
			View: reference.Select("Room", b, map[string]string{
				"label": "Room",
			}, "Room",
				`{{ .type  }} `),
		},
		editor.Field{
			View: reference.Select("Hotel", b, map[string]string{
				"label": "Hotel",
			}, "Hotel",
				`{{ .name  }} `),
		},
	)

	if err != nil {
		return nil, fmt.Errorf("Failed to render Bookings editor view: %s", err.Error())
	}

	return view, nil
}

func init() {
	item.Types["Bookings"] = func() interface{} { return new(Bookings) }
}

func (b *Bookings) AutoApprove(res http.ResponseWriter, req *http.Request) error {
	return nil
}

func (b *Bookings) Create(res http.ResponseWriter, req *http.Request) error {
	return nil
}

func (b *Bookings) String() string {
	return b.Email
}
