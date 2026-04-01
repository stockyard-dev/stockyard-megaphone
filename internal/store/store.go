package store
import("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{*sql.DB}
type Component struct{ID int64 `json:"id"`;Name string `json:"name"`;Status string `json:"status"`;UpdatedAt time.Time `json:"updated_at"`}
type Incident struct{ID int64 `json:"id"`;Title string `json:"title"`;Body string `json:"body"`;Status string `json:"status"`;CreatedAt time.Time `json:"created_at"`;ResolvedAt *string `json:"resolved_at"`}
type Subscriber struct{ID int64 `json:"id"`;Email string `json:"email"`;CreatedAt time.Time `json:"created_at"`}
func Open(d string)(*DB,error){os.MkdirAll(d,0755);dsn:=filepath.Join(d,"megaphone.db")+"?_journal_mode=WAL&_busy_timeout=5000";db,err:=sql.Open("sqlite",dsn);if err!=nil{return nil,fmt.Errorf("open: %w",err)};db.SetMaxOpenConns(1);migrate(db);return &DB{db},nil}
func migrate(db *sql.DB){db.Exec(`CREATE TABLE IF NOT EXISTS components(id INTEGER PRIMARY KEY AUTOINCREMENT,name TEXT NOT NULL,status TEXT DEFAULT 'operational',updated_at DATETIME DEFAULT CURRENT_TIMESTAMP);CREATE TABLE IF NOT EXISTS incidents(id INTEGER PRIMARY KEY AUTOINCREMENT,title TEXT NOT NULL,body TEXT DEFAULT '',status TEXT DEFAULT 'investigating',created_at DATETIME DEFAULT CURRENT_TIMESTAMP,resolved_at TEXT);CREATE TABLE IF NOT EXISTS subscribers(id INTEGER PRIMARY KEY AUTOINCREMENT,email TEXT NOT NULL UNIQUE,created_at DATETIME DEFAULT CURRENT_TIMESTAMP)`)}
func(db *DB)CreateComponent(c *Component)error{res,err:=db.Exec(`INSERT INTO components(name,status)VALUES(?,?)`,c.Name,c.Status);if err!=nil{return err};c.ID,_=res.LastInsertId();return nil}
func(db *DB)ListComponents()([]Component,error){rows,_:=db.Query(`SELECT id,name,status,updated_at FROM components ORDER BY name`);defer rows.Close();var out[]Component;for rows.Next(){var c Component;rows.Scan(&c.ID,&c.Name,&c.Status,&c.UpdatedAt);out=append(out,c)};return out,nil}
func(db *DB)UpdateComponent(id int64,status string){db.Exec(`UPDATE components SET status=?,updated_at=CURRENT_TIMESTAMP WHERE id=?`,status,id)}
func(db *DB)CreateIncident(i *Incident)error{res,err:=db.Exec(`INSERT INTO incidents(title,body,status)VALUES(?,?,?)`,i.Title,i.Body,i.Status);if err!=nil{return err};i.ID,_=res.LastInsertId();return nil}
func(db *DB)ListIncidents()([]Incident,error){rows,_:=db.Query(`SELECT id,title,body,status,created_at,resolved_at FROM incidents ORDER BY created_at DESC LIMIT 50`);defer rows.Close();var out[]Incident;for rows.Next(){var i Incident;rows.Scan(&i.ID,&i.Title,&i.Body,&i.Status,&i.CreatedAt,&i.ResolvedAt);out=append(out,i)};return out,nil}
func(db *DB)Subscribe(email string)error{_,err:=db.Exec(`INSERT OR IGNORE INTO subscribers(email)VALUES(?)`,email);return err}
func(db *DB)Stats()(map[string]interface{},error){var c,s int;db.QueryRow(`SELECT COUNT(*) FROM components`).Scan(&c);db.QueryRow(`SELECT COUNT(*) FROM subscribers`).Scan(&s);return map[string]interface{}{"components":c,"subscribers":s},nil}
