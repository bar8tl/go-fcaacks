// inicopy.go [2014-10-02 BAR8TL]
// Inital copy of FCA ack files into a local folder
package fcaacks

import "bytes"
import "fmt"
import "io/ioutil"
import "log"
import "os/exec"
import "path/filepath"
import "strings"

type Dini_tp struct {
}

func NewDini() *Dini_tp {
  var c Dini_tp
  return &c
}

func (c Dini_tp) IniCopy(s Settings_tp) {
  ffold, fextn := s.Ackdr, s.Acktp
  tfold, textn := s.Xmldr, s.Xmltp
  c.copyFiles(ffold+"archive\\", tfold, fextn, textn)
  c.copyFiles(ffold, tfold, fextn, textn)
}

func (c Dini_tp) copyFiles(ffold, tfold, fextn, textn string) {
  files, _ := ioutil.ReadDir(ffold)
  for _, f := range files {
    filid := f.Name()
    extn := filepath.Ext(filid)
    file := strings.TrimRight(filid, extn)
    tk := strings.Split(file, "_")
    filn := tk[2]
    if extn == fextn {
      log.Println("> copying file", filid)
      cmd := exec.Command("cmd", "/c", "copy "+ffold+filid+" "+tfold+filn+textn)
      var stderr bytes.Buffer
      cmd.Stderr = &stderr
      err := cmd.Run()
      if err != nil {
        log.Fatal(fmt.Sprint(err) + ": " + stderr.String())
      }
    }
  }
}
