// fcaacks.go [2014-10-02 BAR8TL]
// Start FCA XML Invoice Acknowledgments retrieve and report processes
package main

import "bar8tl/p/fcaacks"

func main() {
  s := fcaacks.NewSettings("_config.json", "_deflts.json")
  for _, parm := range s.Cmdpr {
           if parm.Optn == "dbo" {
      dbo := fcaacks.NewDdbo(s)
      dbo.CrtTables()
    } else if parm.Optn == "ini" {
      ini := fcaacks.NewDini()
      ini.IniCopy(s)
    } else if parm.Optn == "ref" {
      ref := fcaacks.NewDref(s)
      ref.RefreshFiles(parm, s)
    } else if parm.Optn == "rep" {
      rep := fcaacks.NewDrep()
      rep.CrtReport(parm, s)
    }
  }
}
