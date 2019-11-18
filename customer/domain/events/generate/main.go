// +build generator

//go:generate go run main.go

package main

import (
	"bytes"
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"unicode"
	"unicode/utf8"

	"github.com/joncalhoun/pipe"
)

const mode = "format" // "format", "noformat", "stdout"

type Field struct {
	FieldName string
	DataType  string
}

type Event struct {
	EventType    string
	EventFactory string
	Fields       []Field
}

var events = []Event{
	{
		EventType:    "Registered",
		EventFactory: "ItWasRegistered",
		Fields: []Field{
			{FieldName: "customerID", DataType: "*values.CustomerID"},
			{FieldName: "confirmableEmailAddress", DataType: "*values.ConfirmableEmailAddress"},
			{FieldName: "personName", DataType: "*values.PersonName"},
		},
	},
	{
		EventType:    "EmailAddressConfirmed",
		EventFactory: "EmailAddressWasConfirmed",
		Fields: []Field{
			{FieldName: "customerID", DataType: "*values.CustomerID"},
			{FieldName: "emailAddress", DataType: "*values.EmailAddress"},
		},
	},
	{
		EventType:    "EmailAddressConfirmationFailed",
		EventFactory: "EmailAddressConfirmationHasFailed",
		Fields: []Field{
			{FieldName: "customerID", DataType: "*values.CustomerID"},
			{FieldName: "confirmationHash", DataType: "*values.ConfirmationHash"},
		},
	},
	{
		EventType:    "EmailAddressChanged",
		EventFactory: "EmailAddressWasChanged",
		Fields: []Field{
			{FieldName: "customerID", DataType: "*values.CustomerID"},
			{FieldName: "confirmableEmailAddress", DataType: "*values.ConfirmableEmailAddress"},
		},
	},
}

type Config struct {
	RelativeOutputPath string
	Events             []Event
}

var config = Config{
	RelativeOutputPath: "..",
}

func main() {
	generateEvents()
	generateTestsForEvents()
}

func generateEvents() {
	var err error

	methodName := func(input string) string {
		parts := strings.Split(input, ".")
		return parts[1]
	}

	lcFirst := func(s string) string {
		if s == "" {
			return ""
		}
		r, n := utf8.DecodeRuneInString(s)
		return string(unicode.ToLower(r)) + s[n:]
	}

	tick := func() string { return "`" }

	for _, event := range events {
		t := template.New(event.EventType)

		t = t.Funcs(
			template.FuncMap{
				"title":      strings.Title,
				"methodName": methodName,
				"eventName":  t.Name,
				"lcFirst":    lcFirst,
				"tick":       tick,
			},
		)

		t, err = t.Parse(eventTemplate)
		die(err)

		switch mode {
		case "format":
			outFile, err := os.Create(config.RelativeOutputPath + "/" + strings.Title(event.EventType) + ".go")
			die(err)

			rc, wc, _ := pipe.Commands(
				exec.Command("gofmt"),
			)

			err = t.Execute(wc, event)
			die(err)

			err = wc.Close()
			die(err)

			_, err = io.Copy(outFile, rc)
			die(err)
		case "noformat":
			outFile, err := os.Create(config.RelativeOutputPath + "/" + strings.Title(event.EventType) + ".go")
			die(err)

			err = t.Execute(outFile, event)
			die(err)
		case "stdout":
			err = t.Execute(os.Stdout, event)
			die(err)
		}
	}
}

func generateTestsForEvents() {
	var err error

	methodName := func(input string) string {
		parts := strings.Split(input, ".")
		return parts[1]
	}

	lcFirst := func(s string) string {
		if s == "" {
			return ""
		}
		r, n := utf8.DecodeRuneInString(s)
		return string(unicode.ToLower(r)) + s[n:]
	}

	valueFactoryForTestTemplates := map[string]string{
		"customerID":              customerIDFactoryForTestTemplate,
		"emailAddress":            emailAddressFactoryForTestTemplate,
		"confirmationHash":        confirmationHashFactoryForTestTemplate,
		"confirmableEmailAddress": confirmableEmailAddressFactoryForTestTemplate,
		"personName":              personNameFactoryForTestTemplate,
	}

	valueFactoryForTest := func(templateName string) string {
		t := template.New("valueFactoryForTest")

		tpl, found := valueFactoryForTestTemplates[templateName]
		if !found {
			die(errors.New("could not find valueFactoryForTest template for: " + templateName))
		}

		t, err = t.Parse(tpl)
		die(err)

		buf := bytes.NewBuffer([]byte{})
		err = t.Execute(buf, false)
		die(err)

		return buf.String()
	}

	for _, event := range events {
		t := template.New(event.EventType)

		t = t.Funcs(
			template.FuncMap{
				"title":               strings.Title,
				"methodName":          methodName,
				"eventName":           t.Name,
				"lcFirst":             lcFirst,
				"valueFactoryForTest": valueFactoryForTest,
			},
		)

		t, err = t.Parse(testTemplate)
		die(err)

		switch mode {
		case "format":
			outFile, err := os.Create(config.RelativeOutputPath + "/" + strings.Title(event.EventType) + "_test.go")
			die(err)

			rc, wc, _ := pipe.Commands(
				exec.Command("gofmt"),
			)

			err = t.Execute(wc, event)
			die(err)

			err = wc.Close()
			die(err)

			_, err = io.Copy(outFile, rc)
			die(err)
		case "noformat":
			outFile, err := os.Create(config.RelativeOutputPath + "/" + strings.Title(event.EventType) + "_test.go")
			die(err)

			err = t.Execute(outFile, event)
			die(err)
		case "stdout":
			err = t.Execute(os.Stdout, event)
			die(err)
		}
	}
}

func die(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var eventTemplate = `
{{$eventVar := lcFirst eventName}}
// Code generated by generate/main.go. DO NOT EDIT.

package events

import (
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
	"reflect"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	jsoniter "github.com/json-iterator/go"
)

const (
	{{lcFirst eventName}}AggregateName       = "Customer"
	{{eventName}}MetaTimestampFormat = time.RFC3339Nano
)

type {{eventName}} struct {
	{{range .Fields}}{{.FieldName}} {{.DataType}}
	{{end -}}
	meta *Meta
}

/*** Factory Methods ***/

func {{.EventFactory}}(
	{{range .Fields}}{{.FieldName}} {{.DataType}},
	{{end -}}
	streamVersion uint,
) *{{eventName}} {

	{{$eventVar}} := &{{eventName}}{
		{{range .Fields}}{{.FieldName}}: {{.FieldName}},
		{{end -}}
	}

	eventType := reflect.TypeOf({{$eventVar}}).String()
	eventTypeParts := strings.Split(eventType, ".")
	eventName := eventTypeParts[len(eventTypeParts)-1]
	eventName = strings.Title(eventName)
	fullEventName := {{$eventVar}}AggregateName + eventName

	{{$eventVar}}.meta = &Meta{
		identifier:    customerID.String(),
		eventName:     fullEventName,
		occurredAt:    time.Now().Format({{eventName}}MetaTimestampFormat),
		streamVersion: streamVersion,
	}

	return {{$eventVar}}
}

/*** Getter Methods ***/

{{range .Fields}}
func ({{$eventVar}} *{{eventName}}) {{methodName .DataType}}() {{.DataType}} {
	return {{$eventVar}}.{{.FieldName}}
}
{{end}}

/*** Implement shared.DomainEvent ***/

func ({{$eventVar}} *{{eventName}}) Identifier() string {
	return {{$eventVar}}.meta.identifier
}

func ({{$eventVar}} *{{eventName}}) EventName() string {
	return {{$eventVar}}.meta.eventName
}

func ({{$eventVar}} *{{eventName}}) OccurredAt() string {
	return {{$eventVar}}.meta.occurredAt
}

func ({{$eventVar}} *{{eventName}}) StreamVersion() uint {
	return {{$eventVar}}.meta.streamVersion
}

/*** Implement json.Marshaler ***/

func ({{$eventVar}} *{{eventName}}) MarshalJSON() ([]byte, error) {
	data := &struct {
		{{range .Fields}}{{methodName .DataType}} {{.DataType}} {{tick}}json:"{{.FieldName}}"{{tick}}
		{{end -}}
		Meta *Meta {{tick}}json:"meta"{{tick}}
	}{
		{{range .Fields}}{{methodName .DataType}}: {{$eventVar}}.{{.FieldName}},
		{{end -}}
		Meta: {{$eventVar}}.meta,
	}

	return jsoniter.Marshal(data)
}

/*** Implement json.Unmarshaler ***/

func ({{$eventVar}} *{{eventName}}) UnmarshalJSON(data []byte) error {
	unmarshaledData := &struct {
		{{range .Fields}}{{methodName .DataType}} {{.DataType}} {{tick}}json:"{{.FieldName}}"{{tick}}
		{{end -}}
		Meta *Meta {{tick}}json:"meta"{{tick}}
	}{}

	if err := jsoniter.Unmarshal(data, unmarshaledData); err != nil {
		return errors.Wrap(errors.Mark(err, shared.ErrUnmarshalingFailed), "{{$eventVar}}.UnmarshalJSON")
	}

	{{range .Fields}}{{$eventVar}}.{{.FieldName}} = unmarshaledData.{{methodName .DataType}}
	{{end -}}

	{{$eventVar}}.meta = unmarshaledData.Meta

	return nil
}
`

var customerIDFactoryForTestTemplate = `customerID := values.GenerateCustomerID()`

var emailAddressFactoryForTestTemplate = `emailAddress, err := values.EmailAddressFrom("foo@bar.com")
	So(err, ShouldBeNil)`

var confirmationHashFactoryForTestTemplate = `confirmationHash := values.GenerateConfirmationHash("secret_hash")`

var confirmableEmailAddressFactoryForTestTemplate = `emailAddress, err := values.EmailAddressFrom("foo@bar.com")
	So(err, ShouldBeNil)
	confirmableEmailAddress := emailAddress.ToConfirmable()`

var personNameFactoryForTestTemplate = `personName, err := values.PersonNameFrom("John", "Doe")
	So(err, ShouldBeNil)`

var testTemplate = `
{{$eventVar := lcFirst eventName}}
// Code generated by generate/main.go. DO NOT EDIT.

package events_test

import (
	"go-iddd/customer/domain/events"
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
	"testing"
	"time"

	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func Test{{.EventFactory}}(t *testing.T) {
	Convey("Given valid parameters as input", t, func() {
		{{range .Fields}}{{valueFactoryForTest .FieldName}}
		{{end}}

		Convey("When a new {{eventName}} event is created", func() {
			streamVersion := uint(1)
			{{$eventVar}} := events.{{.EventFactory}}({{range .Fields}}{{.FieldName}}, {{end}} streamVersion)

			Convey("It should succeed", func() {
				So({{$eventVar}}, ShouldNotBeNil)
				So({{$eventVar}}, ShouldImplement, (*shared.DomainEvent)(nil))
				So({{$eventVar}}, ShouldHaveSameTypeAs, (*events.{{eventName}})(nil))
			})
		})
	})
}

func Test{{eventName}}ExposesExpectedValues(t *testing.T) {
	Convey("Given a {{eventName}} event", t, func() {
		{{range .Fields}}{{valueFactoryForTest .FieldName}}
		{{end -}}
		streamVersion := uint(1)

		beforeItOccurred := time.Now()
		{{$eventVar}} := events.{{.EventFactory}}({{range .Fields}}{{.FieldName}}, {{end}} streamVersion)
		afterItOccurred := time.Now()

		Convey("It should expose the expected values", func() {
			{{range .Fields}}So({{$eventVar}}.{{methodName .DataType}}(), ShouldResemble, {{.FieldName}})
			{{end -}}
			So({{$eventVar}}.Identifier(), ShouldEqual, customerID.String())
			So({{$eventVar}}.EventName(), ShouldEqual, "Customer{{eventName}}")
			itOccurred, err := time.Parse(events.{{eventName}}MetaTimestampFormat, {{$eventVar}}.OccurredAt())
			So(err, ShouldBeNil)
			So(beforeItOccurred, ShouldHappenBefore, itOccurred)
			So(afterItOccurred, ShouldHappenAfter, itOccurred)
			So({{$eventVar}}.StreamVersion(), ShouldEqual, streamVersion)
		})
	})
}

func Test{{eventName}}MarshalJSON(t *testing.T) {
	Convey("Given a {{eventName}} event", t, func() {
		{{range .Fields}}{{valueFactoryForTest .FieldName}}
		{{end -}}
		streamVersion := uint(1)

		{{$eventVar}} := events.{{.EventFactory}}({{range .Fields}}{{.FieldName}}, {{end}} streamVersion)

		Convey("When it is marshaled to json", func() {
			data, err := {{$eventVar}}.MarshalJSON()

			Convey("It should create the expected json", func() {
				So(err, ShouldBeNil)
				So(string(data), ShouldStartWith, "{")
				So(string(data), ShouldEndWith, "}")
			})
		})
	})
}

func Test{{eventName}}UnmarshalJSON(t *testing.T) {
	Convey("Given a {{eventName}} event marshaled to json", t, func() {
		{{range .Fields}}{{valueFactoryForTest .FieldName}}
		{{end -}}
		streamVersion := uint(1)

		{{$eventVar}} := events.{{.EventFactory}}({{range .Fields}}{{.FieldName}}, {{end}} streamVersion)

		data, err := {{$eventVar}}.MarshalJSON()
		So(err, ShouldBeNil)

		Convey("When it is unmarshaled", func() {
			unmarshaled := &events.{{eventName}}{}
			err := unmarshaled.UnmarshalJSON(data)

			Convey("It should be equal to the original {{eventName}} event", func() {
				So(err, ShouldBeNil)
				So({{$eventVar}}, ShouldResemble, unmarshaled)
			})
		})
	})

	Convey("Given invalid json", t, func() {
		data := []byte("666")

		Convey("When it is unmarshaled to {{eventName}} event", func() {
			unmarshaled := &events.{{eventName}}{}
			err := unmarshaled.UnmarshalJSON(data)

			Convey("It should fail", func() {
				So(err, ShouldNotBeNil)
				So(errors.Is(err, shared.ErrUnmarshalingFailed), ShouldBeTrue)
			})
		})
	})
}
`
