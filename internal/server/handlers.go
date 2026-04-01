package server
import("encoding/json";"net/http";"strconv";"github.com/stockyard-dev/stockyard-megaphone/internal/store")
func(s *Server)handleListComponents(w http.ResponseWriter,r *http.Request){list,_:=s.db.ListComponents();if list==nil{list=[]store.Component{}};writeJSON(w,200,list)}
func(s *Server)handleCreateComponent(w http.ResponseWriter,r *http.Request){var c store.Component;json.NewDecoder(r.Body).Decode(&c);if c.Name==""{writeError(w,400,"name required");return};if c.Status==""{c.Status="operational"};s.db.CreateComponent(&c);writeJSON(w,201,c)}
func(s *Server)handleUpdateComponent(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);var req struct{Status string `json:"status"`};json.NewDecoder(r.Body).Decode(&req);s.db.UpdateComponent(id,req.Status);writeJSON(w,200,map[string]string{"status":"updated"})}
func(s *Server)handleListIncidents(w http.ResponseWriter,r *http.Request){list,_:=s.db.ListIncidents();if list==nil{list=[]store.Incident{}};writeJSON(w,200,list)}
func(s *Server)handleCreateIncident(w http.ResponseWriter,r *http.Request){var i store.Incident;json.NewDecoder(r.Body).Decode(&i);if i.Title==""{writeError(w,400,"title required");return};if i.Status==""{i.Status="investigating"};s.db.CreateIncident(&i);writeJSON(w,201,i)}
func(s *Server)handleSubscribe(w http.ResponseWriter,r *http.Request){var req struct{Email string `json:"email"`};json.NewDecoder(r.Body).Decode(&req);if req.Email==""{writeError(w,400,"email required");return};s.db.Subscribe(req.Email);writeJSON(w,201,map[string]string{"status":"subscribed"})}
func(s *Server)handleOverview(w http.ResponseWriter,r *http.Request){m,_:=s.db.Stats();writeJSON(w,200,m)}
