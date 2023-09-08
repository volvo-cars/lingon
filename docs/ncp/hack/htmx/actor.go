package htmx

import (
	_ "embed"
	"ncp/bla"
	"net/http"

	"github.com/go-chi/chi/v5"
)

//go:embed actor.html
var actorHTML string

// func schemaRender(nc *nats.Conn, req bla.Request) (*bla.Reply, error) {
// 	switch {
// 	case req.RequestURI == "/":
// 		reply, err := bla.SendAccountGetMsg(nc, bla.AccountGetMsg{
// 			ID: req.Account.ID,
// 		})
// 		if err != nil {
// 			return nil, fmt.Errorf("getting account: %w", err)
// 		}
// 		schemaListReply, err := bla.SendSchemaListForAccountMsg(
// 			nc,
// 			req.Account.ID,
// 			bla.SchemaListMsg{},
// 		)
// 		if err != nil {
// 			return nil, fmt.Errorf("getting schema list: %w", err)
// 		}
// 		type Data struct {
// 			UserInfo bla.UserInfo
// 			Account  bla.Account
// 			Schemas  []bla.Schema
// 		}
// 		data := Data{
// 			UserInfo: req.UserInfo,
// 			Account:  reply.Account,
// 			Schemas:  schemaListReply.Schemas,
// 		}

// 		body, err := renderTemplate(schemasHTML, data)
// 		if err != nil {
// 			return nil, fmt.Errorf("rendering template: %w", err)
// 		}
// 		return &bla.Reply{
// 			Body: body,
// 		}, nil
// 	default:
// 		// TODO: change to case something instead of "default"
// 		ss := strings.Split(req.RequestURI, "/")
// 		schemaID := ss[len(ss)-1]
// 		schemaReply, err := bla.SendSchemaGetForAccountMsg(
// 			nc,
// 			req.Account.ID,
// 			bla.SchemaGetMsg{
// 				Key: schemaID,
// 			},
// 		)
// 		if err != nil {
// 			return nil, fmt.Errorf("getting schema: %w", err)
// 		}
// 		schema := schemaReply.Schema
// 		consumerInfo, err := bla.ConsumerState(
// 			nc,
// 			req.Account.ID,
// 			schema.Spec.Name,
// 			strings.ReplaceAll(schema.Spec.Version, ".", "_"),
// 		)
// 		if err != nil {
// 			return nil, fmt.Errorf("getting consumer state: %w", err)
// 		}
// 		type Data struct {
// 			UserInfo     bla.UserInfo
// 			Account      bla.Account
// 			Schema       bla.Schema
// 			ConsumerInfo *api.ConsumerInfo
// 		}
// 		d := Data{
// 			UserInfo:     req.UserInfo,
// 			Account:      req.Account,
// 			Schema:       schemaReply.Schema,
// 			ConsumerInfo: consumerInfo,
// 		}
// 		body, err := renderTemplate(schemaHTML, d)
// 		if err != nil {
// 			return nil, fmt.Errorf("rendering template: %w", err)
// 		}
// 		return &bla.Reply{
// 			Body: body,
// 		}, nil
// 	}
// 	return nil, fmt.Errorf("page not found")
// }

func (s *Server) serveActor(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := r.Context().Value(AuthContext).(bla.UserInfo)
	if !ok {
		http.Error(w, "no auth context", http.StatusUnauthorized)
		return
	}
	accountID := chi.URLParam(r, "accountID")
	actor := chi.URLParam(r, "actor")
	subPath := chi.URLParam(r, "*")

	reply, err := bla.SendAccountGetMsg(s.nc, bla.AccountGetMsg{
		ID: accountID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	type Data struct {
		UserInfo bla.UserInfo
		Account  bla.Account
		Actors   []ActorPage
		Actor    string
		Path     string
	}
	data := Data{
		UserInfo: userInfo,
		Account:  reply.Account,
		Actor:    actor,
		Path:     subPath,
		Actors: []ActorPage{
			{
				Name: "Schemas",
				Link: "schema",
			},
		},
	}
	renderPage(w, accountTmpl, actorHTML, data)
}
