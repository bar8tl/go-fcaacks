cd c:\rbhome\go\src\bar8tl\p\
md fcaacks
cd fcaacks
copy c:\c_portab\01_rb\_rbprogs\go-fcaacks\pkg\inicopy.go  .
copy c:\c_portab\01_rb\_rbprogs\go-fcaacks\pkg\dbo.go      .
copy c:\c_portab\01_rb\_rbprogs\go-fcaacks\pkg\refresh.go  .
copy c:\c_portab\01_rb\_rbprogs\go-fcaacks\pkg\report.go   .
copy c:\c_portab\01_rb\_rbprogs\go-fcaacks\pkg\settings.go .
go install
pause
