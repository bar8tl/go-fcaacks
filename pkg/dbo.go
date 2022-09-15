// dbo.go - Creates DB and DB tables to locally store FCA Acks
// [2014-10-02 BAR8TL]
// [2022-09-14 BAR8TL CRQ001:Replace sqlite pkg go1 by mattn

package fcaacks

// import "code.google.com/p/go-sqlite/go1/sqlite3" [CRQ001:remove]
import "database/sql"                  // [CRQ001:add]
import _ "github.com/mattn/go-sqlite3" // [CRQ001:add]
import "log"
// import "strings" [CRQ001:remove]

// [CRQ001:remove]
// const CNN_SQLITE3 = "file:@?file:locked.sqlite?cache=shared&mode=rwc"

type Ddbo_tp struct {
  cnnst string
}

func NewDdbo(s Settings_tp) *Ddbo_tp {
  var d Ddbo_tp
// [CRQ001:remove]
// d.cnnst = strings.Replace(CNN_SQLITE3, "@", s.Dbodr+s.Dbonm, 1)
  d.cnnst = s.Dbodr+s.Dbonm // [CRQ001:add]
  return &d
}

func (d *Ddbo_tp) CrtTables() {
  d.crtAcks().crtLast()
}

func (d *Ddbo_tp) crtAcks() *Ddbo_tp {
// db, _ := sqlite3.Open(d.cnnst) [CRQ001:remove]
  db, _ := sql.Open("sqlite3", d.cnnst) // [CRQ001:add]
  defer db.Close()
  db.Exec(`DROP TABLE IF EXISTS acks;`)
// [CRQ001:remove]
// err := db.Exec(`CREATE TABLE acks (ackno INTEGER PRIMARY KEY, issue TEXT,
// [CRQ001:add]
  _, err := db.Exec(`CREATE TABLE acks (ackno INTEGER PRIMARY KEY, issue TEXT,
    rceiv TEXT, invoi TEXT, serie TEXT, folio TEXT, uuid TEXT, dtime TEXT,
    stats TEXT, errn1 TEXT, errn2 TEXT, notes TEXT);`)
  if err != nil {
    log.Fatalf("Error during table ACKS creation: %v\n", err)
  }
  return d
}

func (d *Ddbo_tp) crtLast() *Ddbo_tp {
//  db, _ := sqlite3.Open(d.cnnst) [CRQ001:remove]
  db, _ := sql.Open("sqlite3", d.cnnst) // [CRQ001:add]
  defer db.Close()
  db.Exec(`DROP TABLE IF EXISTS last;`)
// [CRQ001:remove]
//  err := db.Exec(`CREATE TABLE last (recno text primary key, ackno integer);`)
// [CRQ001:add]
  _, err := db.Exec(`CREATE TABLE last (recno text primary key, ackno integer);`)
  if err != nil {
    log.Fatalf("Error during table LAST creation: %v\n", err)
  }
// [CRQ001:remove]
// err  = db.Exec(`INSERT INTO last(recno,ackno) VALUES("00",0);`)
// [CRQ001:add]
  _, err  = db.Exec(`INSERT INTO last(recno,ackno) VALUES("00",0);`)
  if err != nil {
    log.Fatalf("Error during initializing table LAST: %v\n", err)
  }
  return d
}
