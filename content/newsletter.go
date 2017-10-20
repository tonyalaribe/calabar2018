package content

import (
	"fmt"
	"net/http"

	"github.com/ponzu-cms/ponzu/management/editor"
	"github.com/ponzu-cms/ponzu/system/item"
)

type Newsletter struct {
	item.Item

	Email string `json:"email"`
}

// MarshalEditor writes a buffer of html to edit a Newsletter within the CMS
// and implements editor.Editable
func (n *Newsletter) MarshalEditor() ([]byte, error) {
	view, err := editor.Form(n,
		// Take note that the first argument to these Input-like functions
		// is the string version of each Newsletter field, and must follow
		// this pattern for auto-decoding and auto-encoding reasons:
		editor.Field{
			View: editor.Input("Email", n, map[string]string{
				"label":       "Email",
				"type":        "text",
				"placeholder": "Enter the Email here",
			}),
		},
	)

	if err != nil {
		return nil, fmt.Errorf("Failed to render Newsletter editor view: %s", err.Error())
	}

	return view, nil
}

func init() {
	item.Types["Newsletter"] = func() interface{} { return new(Newsletter) }
}

// String defines how a Newsletter is printed. Update it using more descriptive
// fields from the Newsletter struct type
func (n *Newsletter) String() string {
	return fmt.Sprintf("Newsletter: %s", n.UUID)
}

func (n *Newsletter) Create(res http.ResponseWriter, req *http.Request) error {
	return nil
}
