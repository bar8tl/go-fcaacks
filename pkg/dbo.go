// dbo.go [2014-10-02 BAR8TL]
// Creates DB and DB tables to locally store FCA Acks
package fcaacks

import "code.google.com/p/go-sqlite/go1/sqlite3"
import "log"
import "strings"

const CNN_SQLITE3 = "file:@?file:locked.sqlite?cache=shared&mode=rwc"

type Ddbo_tp struct {
  cnnst string
}

func NewDdbo(s Settings_tp) *Ddbo_tp {
  var d Ddbo_tp
  d.cnnst = strings.Replace(CNN_SQLITE3, "@", s.Dbodr+s.Dbonm, 1)
  return &d
}

func (d *Ddbo_tp) CrtTables() {
  d.crtAcks().crtLast()
}

func (d *Ddbo_tp) crtAcks() *Ddbo_tp {
  db, _ := sqlite3.Open(d.cnnst)
  defer db.Close()
  db.Exec(`DROP TABLE IF EXISTS acks;`)
  err := db.Exec(`CREATE TABLE acks (ackno INTEGER PRIMARY KEY, issue TEXT,
    rceiv TEXT, invoi TEXT, serie TEXT, folio TEXT, uuid TEXT, dtime TEXT,
    stats TEXT, errn1 TEXT, errn2 TEXT, notes TEXT);`)
  if err != nil {
    log.Fatalf("Error during table ACKS creation: %v\n", err)
  }
  return d
}

func (d *Ddbo_tp) crtLast() *Ddbo_tp {
  db, _ := sqlite3.Open(d.cnnst)
  defer db.Close()
  db.Exec(`DROP TABLE IF EXISTS last;`)
  err := db.Exec(`CREATE TABLE last (recno text primary key, ackno integer);`)
  if err != nil {
    log.Fatalf("Error during table LAST creation: %v\n", err)
  }
  err  = db.Exec(`INSERT INTO last(recno,ackno) VALUES("00",0);`)
  if err != nil {
    log.Fatalf("Error during initializing table LAST: %v\n", err)
  }
  return d
}
