package bla

import (
	"encoding/json"
	"fmt"
	"strings"

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
const ActorRenderSubject = "actor.%s.render"

func AccountFromSubject(subject string) (string, error) {
	subjects := strings.Split(subject, ".")
	if len(subjects) != 5 {
		return "", fmt.Errorf("invalid subject: %s", subject)
	}
	return subjects[4], nil
}

// EXPERIMENT: an Actor model for nats.

type Actor interface {
	// Subscribe handles initialisation and subscribing the
	// actor to the necessary subjects
	Subscribe() error
	Unsubscribe() error
}

func ActorActionSubscribe[M any, R any](
	nc *nats.Conn,
	actor string,
	action string,
	cb func(msg *M) (*R, error),
) (*nats.Subscription, error) {
	return nc.QueueSubscribe(
		fmt.Sprintf("actor.%s.%s", actor, action),
		action,
		func(msg *nats.Msg) {
			var m M
			if err := json.Unmarshal(msg.Data, &m); err != nil {
				// TODO: handle error
				return
			}
			reply, err := cb(&m)
			if err != nil {
				// TODO: handle error
				return
			}
			bReply, err := json.Marshal(reply)
			if err != nil {
				// TODO: handle error
				return
			}
			_ = msg.Respond(bReply)
		},
	)
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
	sub, err := ActorActionSubscribe[DoMsg, DoReply](da.Conn, "dummy", "something", da.doSomething)
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
