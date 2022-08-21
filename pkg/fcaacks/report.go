// report.go [2017-05-24 BAR8TL]
// Maintains local acks DB and geneates an excel output report
package fcaacks

import lib "bar8tl/p/rblib"
import "code.google.com/p/go-sqlite/go1/sqlite3"
import "encoding/xml"
import "fmt"
import "github.com/tealeg/xlsx"
import "io/ioutil"
import "log"
import "os"
import "path/filepath"
import "strings"
import "strconv"
import "time"

type route struct {
  Issue string `xml:"remitente,attr"`
  Rceiv string `xml:"destinatario,attr"`
}

type docum struct {
  Invoi string `xml:"referenciaProveedor,attr"`
  Serie string `xml:"serie,attr"`
  Folio string `xml:"folioFiscal,attr"`
  Uuid  string `xml:"UUID,attr"`
}

type recei struct {
  Dtime string `xml:"fechahora,attr"`
  Stats string `xml:"estatus,attr"`
}

type Acknw struct {
  Route route    `xml:"ruta"`
  Docum docum    `xml:"documento"`
  Recei recei    `xml:"recepcion"`
  Errds []string `xml:"error"`
}

var ticks int
var RESET_TICKS = func() {
  ticks = 0
}

var TICK = func() {
  if ticks++; ticks%10 == 0 {
    fmt.Print(".")
  }
}

type Drep_tp struct {
  cnnst string
  db    *sqlite3.Conn
}

func NewDrep() *Drep_tp {
  var d Drep_tp
  return &d
}

func (d *Drep_tp) CrtReport(parm lib.Param_tp, s Settings_tp) *Drep_tp {
  s.SetRunVars(parm)
  d.cnnst = s.Cnnst
  var err error
  fmt.Print("Browsing FCA XML acknowledgments")
  RESET_TICKS()
  d.db, err = sqlite3.Open(d.cnnst)
  if err != nil {
    log.Fatalf("Open SQLite database error: %v\n", err)
  }
  ifile := 0
  lsack := d.getLast()
  ilfil := lsack
  files, _ := ioutil.ReadDir(s.Xmldr)
  for _, f := range files {
    filid := f.Name()
    extn := filepath.Ext(filid)
    file := strings.TrimRight(filid, extn)
    ifile, _ = strconv.Atoi(file)
    if extn == s.Xmltp && ifile > lsack {
      d.processAck(s, filid, file, ifile)
      if ifile > ilfil {
        ilfil = ifile
      }
    }
  }
  d.db.Exec(`UPDATE last SET ackno=? where recno="00";`, ilfil)
  d.noteCorrections()
  d.listAcks(s)
  return d
}

func (d *Drep_tp) getLast() (lsack int) {
  cmd, _ := d.db.Query(`select ackno from last where recno="00";`)
  cmd.Scan(&lsack)
  cmd.Close()
  return
}

func (d *Drep_tp) processAck(s Settings_tp, filid, file string, ifile int) int {
  f, err := os.Open(s.Xmldr + filid)
  if err != nil {
    log.Printf("Open XML file %s error: %v\n", filid, err)
    return ifile
  }
  defer f.Close()
  var a Acknw
  xmlv, _ := ioutil.ReadAll(f)
  err = xml.Unmarshal(xmlv, &a)
  if err != nil {
    log.Printf("Unmarshall XML file %s error: %v\n", filid, err)
    return ifile
  }
  d.isrAcks(file, ifile, &a)
  TICK()
  return ifile
}

func (d *Drep_tp) isrAcks(file string, ifile int, a *Acknw) {
  var err1  string
  var err2  string
  var notes string
  if len(a.Errds) >= 2 {
    err1 = a.Errds[0]
    err2 = a.Errds[1]
  }
  if len(a.Errds) == 1 {
    err1 = a.Errds[0]
  }
  args := sqlite3.NamedArgs{
    ":01": ifile,
    ":02": a.Route.Issue,
    ":03": a.Route.Rceiv,
    ":04": a.Docum.Invoi,
    ":05": a.Docum.Serie,
    ":06": a.Docum.Folio,
    ":07": a.Docum.Uuid,
    ":08": a.Recei.Dtime,
    ":09": a.Recei.Stats,
    ":10": err1,
    ":11": err2,
    ":12": notes,
  }
  err := d.db.Exec(
    `insert into acks values(:01,:02,:03,:04,:05,:06,:07,:08,:09,:10,:11,:12)`,
    args)
  if err != nil {
    log.Fatalf("Insert %s sql table error: %v\n", file, err)
  }
}

func (d *Drep_tp) noteCorrections() {
  var err error
  var ackno, serie, folio, stats, rackn string
  var rdb, sdb *sqlite3.Stmt
  for rdb, err = d.db.Query(
    `select ackno,serie,folio,stats from acks where stats<>"00"
    order by ackno;`); err == nil; err = rdb.Next() {
    rdb.Scan(&ackno, &serie, &folio, &stats)
    sdb, err = d.db.Query(
      `select ackno from acks where stats="00" and serie=? and folio=? and
      ackno<>?;`, serie, folio, ackno)
    if err == nil {
      sdb.Scan(&rackn)
      notes := "Corrected with ack# " + rackn
      d.db.Exec(`update acks set notes=? where ackno=?;`, notes, ackno)
    }
  }
  rdb.Close()
}

type layout struct {
  fld string
  wdt float64
  dsc string
}
type xlsCtrl struct {
  ofile *xlsx.File
  sheet *xlsx.Sheet
  row   *xlsx.Row
  cell  *xlsx.Cell
  style []*xlsx.Style
}

func (d *Drep_tp) listAcks(s Settings_tp) {
  c := xlsCtrl{
    ofile: xlsx.NewFile(),
    style: make([]*xlsx.Style, 2),
  }
  c.style[0] = xlsx.NewStyle()
  c.style[0].Font = *xlsx.NewFont(10, "Arial")
  c.style[0].Font.Bold = true
  c.style[1] = xlsx.NewStyle()
  c.style[1].Font = *xlsx.NewFont(10, "Arial")
  c.sheet = c.ofile.AddSheet("Sheet1")
  column := []layout{
    {"ackno", 10.29, "ack#"},
    {"issue", 11.30, "remitente"},
    {"rceiv", 12.00, "destinatario"},
    {"invoi", 13.14, "refProveedor"},
    {"serie", 5.30,  "serie"},
    {"folio", 7.00,  "folio"},
    {"uuid" , 41.00, "uuid"},
    {"dtime", 18.90, "fechahora"},
    {"stats", 7.60,  "estatus"},
    {"errn1", 33.00, "error"},
    {"errn2", 33.00, "error"},
    {"notas", 10.00, "notas"},
  }
  for i := 0; i < len(column); i++ {
    c.sheet.SetColWidth(i, i, column[i].wdt)
  }
  c.sheet.SheetViews = make([]xlsx.SheetView, 1)
  c.sheet.SheetViews[0] = *new(xlsx.SheetView)
  c.row = c.sheet.AddRow()
  for i := 0; i < len(column); i++ {
    c.cell = c.row.AddCell()
    c.cell.SetStyle(c.style[0])
    c.cell.Value = column[i].dsc
  }
  var f [12]string
  var wdate time.Time
  var err error
  var cmd *sqlite3.Stmt
  for cmd, err = d.db.Query(
    `select * from acks order by dtime, ackno;`); err == nil; err = cmd.Next() {
    cmd.Scan(&f[0], &f[1], &f[2], &f[3], &f[4], &f[5], &f[6], &f[7], &f[8],
      &f[9], &f[10], &f[11])
    ws := f[7]
    nyr, _ := strconv.Atoi(ws[0:4])
    nmn, _ := strconv.Atoi(ws[5:7])
    ndy, _ := strconv.Atoi(ws[8:10])
    wdate = time.Date(nyr, time.Month(nmn), ndy, 0, 0, 0, 0, time.UTC)
    if s.Fdate.Before(wdate) && s.Tdate.After(wdate) {
      c.row = c.sheet.AddRow()
      for i := 0; i < len(column); i++ {
        c.cell = c.row.AddCell()
        c.cell.SetStyle(c.style[1])
        c.cell.Value = f[i]
      }
    }
  }
  if err := c.ofile.Save(s.Dbodr + "zdcmex.xlsx"); err != nil {
    log.Fatalf("Save xlsx file error: %v\n", err)
  }
  cmd.Close()
}
