package helpers

import (
	"flag"
	"log"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var database = flag.Bool("database", false, "run database integration tests")

func TestTokenizeAndDetokenize(t *testing.T) {
	Convey("Given a text it can be tokenized and detokenized", t, func() {
		token, err := Tokenize("example")
		if err != nil {
			t.Fatal(err)
		}

		So(token, ShouldNotBeEmpty)

		detoken, err := GetFromToken(token)
		if err != nil {
			t.Fatal(err)
		}
		So(detoken, ShouldEqual, "example")
	})
}

func TestDetokenizeWrongToken(t *testing.T) {
	Convey("Given an invalid token it should return an error", t, func() {
		_, err := GetFromToken("inventedtoken")
		log.Printf("%v", err)
		So(err, ShouldNotBeNil)
	})
}
