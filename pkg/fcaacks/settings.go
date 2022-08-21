// settings.go [2017-05-24 BAR8TL]
// Define FCA Acks report pgm-level & run-level settings
package fcaacks

import lib "bar8tl/p/rblib"
import "encoding/json"
import "io/ioutil"
import "log"
import "os"
import "strconv"
import "strings"
import "time"

// config: Reads config file and gets program/run parameters
type Program_tp struct {
  Ackdr string `json:"ackDir"`
  Acktp string `json:"ackType"`
  Xmldr string `json:"xmlDir"`
  Xmltp string `json:"xmlType"`
  Dbodr string `json:"dbDir"`
  Dbonm string `json:"dbName"`
}

type Run_tp struct {
  Optcd string `json:"option"`
  Filtr string `json:"rptFilter"`
  Fprm1 string `json:"filtPrm1"`
  Fprm2 string `json:"filtPrm2"`
}

type Config_tp struct {
  Progm Program_tp `json:"program"`
  Run   []Run_tp   `json:"run"`
}

func (c *Config_tp) NewConfig(fname string) {
  f, err := os.Open(fname)
  if err != nil {
    log.Fatalf("File %s opening error: %s\n", fname, err)
  }
  defer f.Close()
  jsonv, _ := ioutil.ReadAll(f)
  err = json.Unmarshal(jsonv, &c)
  if err != nil {
    log.Fatalf("File %s reading error: %s\n", fname, err)
  }
}

// defaults: Reads defaults file and gets program/run defaults settings
type Dflt_tp struct {
  CNNS_SQLIT3   string `json:"CNNS_SQLIT3"`
  ACKS_DIR      string `json:"ACKS_DIR"     `
  ACKS_TYPE     string `json:"ACKS_TYPE"    `
  XML_DIR       string `json:"XML_DIR"      `
  XML_TYPE      string `json:"XML_TYPE"     `
  DB_DIR        string `json:"DB_DIR"       `
  DB_NAME       string `json:"DB_NAME"      `
  REPORT_FILTER string `json:"REPORT_FILTER"`
  FILTER_PARM1  string `json:"FILTER_PARM1" `
  FILTER_PARM2  string `json:"FILTER_PARM2" `
}

type Deflts_tp struct {
  Dflt Dflt_tp `json:"dflt"`
}

func (d *Deflts_tp) NewDeflts(fname string) {
  f, err := os.Open(fname)
  if err != nil {
    log.Fatalf("File %s open error: %s\n", fname, err)
  }
  defer f.Close()
  jsonv, _ := ioutil.ReadAll(f)
  err = json.Unmarshal(jsonv, &d)
  if err != nil {
    log.Fatalf("File %s read error: %s\n", fname, err)
  }
}

// envmnt: Define/Maintain global environment variables
type Envmnt_tp struct {
  Cnnsq string
  Cnnst string
  Ackdr string
  Acktp string
  Xmldr string
  Xmltp string
  Dbodr string
  Dbonm string
  Filtr string
  Fprm1 string
  Fprm2 string
  Found bool
  Dtsys time.Time
  Dtcur time.Time
  Dtnul time.Time
  Fdate time.Time
  Tdate time.Time
}

func (e *Envmnt_tp) NewEnvmnt(s Settings_tp) {
  e.Cnnsq = s.Dflt.CNNS_SQLIT3
  e.Ackdr =
    lib.Ternary_op(len(s.Progm.Ackdr) > 0, s.Progm.Ackdr, s.Dflt.ACKS_DIR)
  e.Acktp =
    lib.Ternary_op(len(s.Progm.Acktp) > 0, s.Progm.Acktp, s.Dflt.ACKS_TYPE)
  e.Xmldr =
    lib.Ternary_op(len(s.Progm.Xmldr) > 0, s.Progm.Xmldr, s.Dflt.XML_DIR)
  e.Xmltp =
    lib.Ternary_op(len(s.Progm.Xmltp) > 0, s.Progm.Xmltp, s.Dflt.XML_TYPE)
  e.Dbodr =
    lib.Ternary_op(len(s.Progm.Dbodr) > 0, s.Progm.Dbodr, s.Dflt.DB_DIR)
  e.Dbonm =
    lib.Ternary_op(len(s.Progm.Dbonm) > 0, s.Progm.Dbonm, s.Dflt.DB_NAME)
  e.Dtsys = time.Now()
  e.Dtcur = time.Now()
  e.Dtnul = time.Date(1901, 1, 1, 0, 0, 0, 0, time.UTC)
  e.Fdate = e.Dtnul
  e.Tdate = e.Dtcur
}

// settings: Container of pgm-level & run-level settings
type Settings_tp struct {
  Config_tp
  lib.Parms_tp
  Deflts_tp
  Envmnt_tp
}

func NewSettings(cfnam, dfnam string) Settings_tp {
  var s Settings_tp
  s.NewParms()
  s.NewConfig(cfnam)
  s.NewDeflts(dfnam)
  s.NewEnvmnt(s)
  return s
}

func (s *Settings_tp) SetRunVars(p lib.Param_tp) {
  s.Found = false
  for _, run := range s.Run {
    if p.Optn == run.Optcd {
      s.Filtr =
        lib.Ternary_op(len(run.Filtr) > 0, run.Filtr, s.Dflt.REPORT_FILTER)
      s.Fprm1 =
        lib.Ternary_op(len(run.Fprm1) > 0, run.Fprm1, s.Dflt.FILTER_PARM1)
      s.Fprm2 =
        lib.Ternary_op(len(run.Fprm2) > 0, run.Fprm2, s.Dflt.FILTER_PARM2)
      s.Found = true
      break
    }
  }
  if p.Optn == "rep" {
    if s.Filtr == "current" {
      s.Tdate = s.Dtsys
      tyear := s.Tdate.Year()
      tmnth := s.Tdate.Month()
      tday  := s.Tdate.Day()
      if s.Fprm1 == "year" {
        s.Fdate = time.Date(tyear, 1, 1, 0, 0, 0, 0, time.UTC)
      } else if s.Fprm1 == "month" {
        s.Fdate = time.Date(tyear, tmnth, 1, 0, 0, 0, 0, time.UTC)
      } else if s.Fprm1 == "day" {
        s.Fdate = time.Date(tyear, tmnth, tday, 0, 0, 0, 0, time.UTC)
      }
    } else if s.Filtr == "past" {
      s.Tdate = s.Dtsys
      dday, _ := strconv.Atoi(s.Fprm1)
      s.Fdate = s.Tdate.AddDate(0, 0, -dday)
    }
  }
  s.Cnnst = strings.Replace(s.Cnnsq, "@", s.Dbodr+s.Dbonm, 1)
}
