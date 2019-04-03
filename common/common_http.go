// Package common outlines common helper methods and types.
package common

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
)

/* BEGIN EXPORTED METHODS */

// GetCtxValue fetches the value at a given key in a fasthttp context.
func GetCtxValue(ctx *fasthttp.RequestCtx, key string) []byte {
	jsonMap := make(map[string]*json.RawMessage) // Init JSON map buffer

	json.Unmarshal(ctx.PostBody(), &jsonMap) // Unmarshal

	if formData := ctx.FormValue(key); formData != nil { // Check has form data
		return formData // Return form data
	} else if userValue := ctx.UserValue(key); userValue != nil { // Check has user value
		return []byte(userValue.(string)) // Return user value
	} else if jsonValue := jsonMap["key"]; jsonValue != nil { // Check has JSON value
		bytes, _ := jsonValue.MarshalJSON() // Marshal

		return bytes // Return JSON value
	}

	return nil // No value
}

/* END EXPORTED METHODS */
