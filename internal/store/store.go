package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type Incident struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Status string `json:"status"`
	Severity string `json:"severity"`
	AffectedServices string `json:"affected_services"`
	Message string `json:"message"`
	UpdatedBy string `json:"updated_by"`
	ResolvedAt string `json:"resolved_at"`
	CreatedAt string `json:"created_at"`
}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"megaphone.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS incidents(id TEXT PRIMARY KEY,title TEXT NOT NULL,status TEXT DEFAULT 'investigating',severity TEXT DEFAULT 'minor',affected_services TEXT DEFAULT '',message TEXT DEFAULT '',updated_by TEXT DEFAULT '',resolved_at TEXT DEFAULT '',created_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)Create(e *Incident)error{e.ID=genID();e.CreatedAt=now();_,err:=d.db.Exec(`INSERT INTO incidents(id,title,status,severity,affected_services,message,updated_by,resolved_at,created_at)VALUES(?,?,?,?,?,?,?,?,?)`,e.ID,e.Title,e.Status,e.Severity,e.AffectedServices,e.Message,e.UpdatedBy,e.ResolvedAt,e.CreatedAt);return err}
func(d *DB)Get(id string)*Incident{var e Incident;if d.db.QueryRow(`SELECT id,title,status,severity,affected_services,message,updated_by,resolved_at,created_at FROM incidents WHERE id=?`,id).Scan(&e.ID,&e.Title,&e.Status,&e.Severity,&e.AffectedServices,&e.Message,&e.UpdatedBy,&e.ResolvedAt,&e.CreatedAt)!=nil{return nil};return &e}
func(d *DB)List()[]Incident{rows,_:=d.db.Query(`SELECT id,title,status,severity,affected_services,message,updated_by,resolved_at,created_at FROM incidents ORDER BY created_at DESC`);if rows==nil{return nil};defer rows.Close();var o []Incident;for rows.Next(){var e Incident;rows.Scan(&e.ID,&e.Title,&e.Status,&e.Severity,&e.AffectedServices,&e.Message,&e.UpdatedBy,&e.ResolvedAt,&e.CreatedAt);o=append(o,e)};return o}
func(d *DB)Update(e *Incident)error{_,err:=d.db.Exec(`UPDATE incidents SET title=?,status=?,severity=?,affected_services=?,message=?,updated_by=?,resolved_at=? WHERE id=?`,e.Title,e.Status,e.Severity,e.AffectedServices,e.Message,e.UpdatedBy,e.ResolvedAt,e.ID);return err}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM incidents WHERE id=?`,id);return err}
func(d *DB)Count()int{var n int;d.db.QueryRow(`SELECT COUNT(*) FROM incidents`).Scan(&n);return n}

func(d *DB)Search(q string, filters map[string]string)[]Incident{
    where:="1=1"
    args:=[]any{}
    if q!=""{
        where+=" AND (title LIKE ?)"
        args=append(args,"%"+q+"%");
    }
    if v,ok:=filters["status"];ok&&v!=""{where+=" AND status=?";args=append(args,v)}
    if v,ok:=filters["severity"];ok&&v!=""{where+=" AND severity=?";args=append(args,v)}
    rows,_:=d.db.Query(`SELECT id,title,status,severity,affected_services,message,updated_by,resolved_at,created_at FROM incidents WHERE `+where+` ORDER BY created_at DESC`,args...)
    if rows==nil{return nil};defer rows.Close()
    var o []Incident;for rows.Next(){var e Incident;rows.Scan(&e.ID,&e.Title,&e.Status,&e.Severity,&e.AffectedServices,&e.Message,&e.UpdatedBy,&e.ResolvedAt,&e.CreatedAt);o=append(o,e)};return o
}

func(d *DB)Stats()map[string]any{
    m:=map[string]any{"total":d.Count()}
    rows,_:=d.db.Query(`SELECT status,COUNT(*) FROM incidents GROUP BY status`)
    if rows!=nil{defer rows.Close();by:=map[string]int{};for rows.Next(){var s string;var c int;rows.Scan(&s,&c);by[s]=c};m["by_status"]=by}
    return m
}
