package bla

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
)

// Subject format for actor action subjects.
//
// For actions which are targetting a specific account, the subject format is:
//
//	actor.<actor>.action.<action>.<account>
//
// For actions which are not targetting a specific account, the subject format is:
//
//	actor.<actor>.action.<action>

const ActionSubjectSubscribeForAccount = "actor.%s.action.%s.*"
const ActionSubjectSubscribe = "actor.%s.action.%s"
const ActionSubjectSendForAccount = "actor.%s.action.%s.%s"
const ActionSubjectSend = "actor.%s.action.%s"

const ActionImportSubject = "actor.*.action.*.%s"
const ActionImportLocalSubject = "actor.*.action.*"
const ActionExportSubject = "actor.*.action.*.*"

// const ActorActionSubjectWildcar = "actor.*.action.*"
const ActorSubjectHTTPRender = "actor.%s.http.render"

func RequestSubject[M any, R any](
	nc *nats.Conn,
	msg M,
	subject string,
) (*R, error) {
	slog.Info("sending request", "subject", subject)
	bData, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}
	replyMsg, err := nc.Request(
		subject,
		bData,
		time.Second*10,
	)
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}
	var reply R
	if err := json.Unmarshal(replyMsg.Data, &reply); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	return &reply, nil
}

func AccountFromSubject(subject string) (string, error) {
	subjects := strings.Split(subject, ".")
	if len(subjects) != 5 {
		return "", fmt.Errorf("invalid subject: %s", subject)
	}
	return subjects[4], nil
}

var TemplateFuncMap = template.FuncMap{
	"durationSince": func(t time.Time) string {
		return time.Since(t).Truncate(time.Second).String()
	},
}

func renderTemplate(tmpl string, data any) ([]byte, error) {
	t, err := template.New("tmpl").Funcs(TemplateFuncMap).Parse(tmpl)
	if err != nil {
		return nil, fmt.Errorf("parse template: %w", err)
	}
	buf := bytes.Buffer{}
	if err := t.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("execute template: %w", err)
	}
	return buf.Bytes(), nil
}

type UserInfo struct {
	Sub     string   `json:"sub"`
	Iss     string   `json:"iss"`
	Name    string   `json:"name"`
	Email   string   `json:"email"`
	Groups  []string `json:"groups"`
	Picture string   `json:"picture"`
}

// Types for HTTP subjects.

// Request is a request for an hx-get.
//
// TODO: is the actual Account needed? It requires an extra fetch.
// Would the accounnt ID suffice, and if really needed, the callee
// can get it itself?
type Request struct {
	// UserInfo of the user making the request.
	UserInfo UserInfo
	// Account for which the request is made.
	Account Account

	// RequestURI is the actor-scoped request-target.
	RequestURI string

	// Method specifies the HTTP method (GET, POST, PUT, etc.).
	Method string

	// Params is set by the handler (subscriber) not by the sender.
	Params map[string]string
}

// Reply is a reply for an hx-get.
type Reply struct {
	// TODO: Header and other necessary fields.
	// Header http.Header
	Body []byte
	// Error string
}

// EXPERIMENT: an Actor model for nats.

func SubscribeToSubject[M any, R any](
	nc *nats.Conn,
	subject string,
	cb func(msg *M) (*R, error),
) (*nats.Subscription, error) {
	return nc.QueueSubscribe(
		subject,
		"actor",
		func(msg *nats.Msg) {
			var m M
			if err := json.Unmarshal(msg.Data, &m); err != nil {
				slog.Error("unmarshal msg", "error", err)
				// TODO: handle error
				return
			}
			reply, err := cb(&m)
			if err != nil {
				slog.Error("handling message", "error", err)
				// TODO: handle error
				return
			}
			bReply, err := json.Marshal(reply)
			if err != nil {
				slog.Error("marshalling reply", "error", err)
				// TODO: handle error
				return
			}
			_ = msg.Respond(bReply)
		},
	)
}

func SubscribeToSubjectWithAccount[M any, R any](
	nc *nats.Conn,
	subject string,
	cb func(accountID string, msg *M) (*R, error),
) (*nats.Subscription, error) {
	return nc.QueueSubscribe(
		subject,
		"actor",
		func(msg *nats.Msg) {
			logger := slog.With("subject", msg.Subject)

			respond := func(reply any) {
				bReply, err := json.Marshal(reply)
				if err != nil {
					logger.Error("marshalling nats reply", "error", err)
					err = FromError(err)
					// Try marshalling a new error
					bReply, err = json.Marshal(err)
					if err != nil {
						logger.Error(
							"marshalling second nats reply",
							"error",
							err,
						)
						return
					}
				}
				if err := msg.Respond(bReply); err != nil {
					logger.Error("sending nats response", "error", err)
					return
				}
			}

			var m M
			if err := json.Unmarshal(msg.Data, &m); err != nil {
				apiErr := &Error{
					Status:  http.StatusBadRequest,
					Message: fmt.Sprintf("unmarshalling nats message: %s", err),
				}
				logger.Error("unmarshalling nats message", "error", apiErr)
				respond(FromError(apiErr))
				return
			}
			accountPubKey, err := AccountFromSubject(msg.Subject)
			if err != nil {
				apiErr := &Error{
					Status: http.StatusBadRequest,
					Message: fmt.Sprintf(
						"getting the account from subject: %s",
						err,
					),
				}
				logger.Error(
					"getting the account from subject",
					"error",
					apiErr,
				)
				respond(FromError(apiErr))
				return
			}
			reply, err := cb(accountPubKey, &m)
			if err != nil {
				// Convert err into Errors.
				err = FromError(err)
				respond(err)
				return
			}
			respond(reply)
		},
	)
}

type Actor interface {
	// Subscribe handles initialisation and subscribing the
	// actor to the necessary subjects
	Subscribe() error
	Unsubscribe() error
}

var _ Actor = (*DummyActor)(nil)

type DummyActor struct {
	Conn *nats.Conn
}

type DoMsg struct{}
type DoReply struct{}

func (da *DummyActor) doSomething(msg *DoMsg) (*DoReply, error) {
	return &DoReply{}, nil
}

type ActionFn[M any, R any] func(msg *M) (*R, error)

// type ActionFn[M any, R any] func(
// 	nc *nats.Conn,
// 	actor string,
// 	action string,
// 	cb func(msg *M) (*R, error),
// ) (*nats.Subscription, error)

// Subscribe implements Actor.
func (da *DummyActor) Subscribe() error {
	// How can we store a list of something with type parameters??
	sub, err := SubscribeToSubject[DoMsg, DoReply](
		da.Conn,
		fmt.Sprintf(ActionSubjectSubscribe, "dummy", "something"),
		da.doSomething,
	)
	if err != nil {
		return fmt.Errorf("yo: %w", err)
	}
	// TODO: where to call this?
	defer sub.Unsubscribe()
	return nil
}

// Unsubscribe implements Actor.
func (*DummyActor) Unsubscribe() error {
	panic("unimplemented")
}
