// refresh.go - Keep FCA acks local folder updated with FCA acks
// [2017-05-24 BAR8TL]
// [2022-09-14 BAR8TL CRQ001:Replace sqlite pkg go1 by mattn

package fcaacks

import lib "bar8tl/p/rblib"
import "bytes"
// import "code.google.com/p/go-sqlite/go1/sqlite3" [CRQ001:remove]
import "database/sql"                  // [CRQ001:add]
import _ "github.com/mattn/go-sqlite3" // [CRQ001:add]
import "fmt"
import "io/ioutil"
import "log"
import "os/exec"
import "path/filepath"
import "strconv"
import "strings"

type Dref_tp struct {
  cnnst string
}

func NewDref(s Settings_tp) *Dref_tp {
  var d Dref_tp
  return &d
}

func (d *Dref_tp) RefreshFiles(parm lib.Param_tp, s Settings_tp) {
  s.SetRunVars(parm)
// d.cnnst = s.Cnnst [CRQ001:remove]
  d.cnnst = s.Dbodr+s.Dbonm // [CRQ001:add]

  files, err := ioutil.ReadDir(s.Ackdr)
  if err != nil {
    log.Fatalf("Directory reading error: %v\n", err)
  }
// db, err := sqlite3.Open(d.cnnst) [CRQ001:remove]
  db, _ := sql.Open("sqlite3", d.cnnst) // [CRQ001:add]
  if err != nil {
    log.Fatalf("Open SQLite database error: %v\n", err)
  }
  for _, f := range files {
    if f.IsDir() {
      continue
    }
    filid := f.Name()
    extn := filepath.Ext(filid)
    file := strings.TrimRight(filid, extn)
    tk := strings.Split(file, "_")
    filn := tk[2]
    ifile, _ := strconv.Atoi(filn)
    if strings.ToLower(extn) == strings.ToLower(s.Acktp) {
      _, err := db.Query(`select ackno from acks where ackno=?;`, ifile)
      if err != nil {
        log.Println("> copying file", filid)
        cmd := exec.Command("cmd", "/c", "copy "+s.Ackdr+filid+" "+s.Xmldr+filn+
          s.Xmltp)
        var stderr bytes.Buffer
        cmd.Stderr = &stderr
        err := cmd.Run()
        if err != nil {
          log.Fatal(fmt.Sprint(err) + ": " + stderr.String())
        }
      }
    }
  }
  db.Close()
}
