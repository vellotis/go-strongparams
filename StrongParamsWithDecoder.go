package strongparams

import (
	"github.com/gorilla/schema"
	"github.com/pkg/errors"
)

var defaultDecoder = func() *schema.Decoder {
	decoder := schema.NewDecoder()
	decoder.SetAliasTag("params") // Use `params` tags instead of `schema`
	return decoder
}()

type StrongParamsWithDecoder struct {
	decoder *schema.Decoder
}

// WithDecoder declares the decoder to be used for Params mechanism.
//   paramsWithDecoder, err := WithDecoder(schema.NewDecoder())
//   // handle error
//   queryParams := paramsWithDecoder.Params(request).Query()
// Returns an error if the passed decoder is a `nil` value.
func WithDecoder(decoder *schema.Decoder) (*StrongParamsWithDecoder, error) {
	if decoder == nil {
		return nil, errors.New("`decoder` parameter cannot be `nil`")
	}

	return &StrongParamsWithDecoder{
		decoder: decoder,
	}, nil
}

// WithDecoderSafe is an equivalent of WithDecoder but instead of returning an error when the passed decoder is a `nil`
// value it panics with the same error.
//   WithDecoder(schema.NewDecoder()).Params(request).Query()
func WithDecoderSafe(decoder *schema.Decoder) *StrongParamsWithDecoder {
	paramsWithDecoder, err := WithDecoder(decoder)
	if err != nil {
		panic(errors.New(err.Error()))
	}
	return paramsWithDecoder
}

// Params declares the *http.Request to be used for the strong parameters mechanism.
//
// **NOTE** The method panics if the passed request parameter is nil.
func (this *StrongParamsWithDecoder) Params() *StrongParams {
	return Params().WithDecoder(this.decoder)
}
