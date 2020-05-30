// package main is a sample for extending applications
// MicroCore Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)

package main

import (
	"github.com/Dobryvechir/microcore/pkg/dvconfig"
	_ "github.com/Dobryvechir/microcore/pkg/dvdbdata"
	_ "github.com/Dobryvechir/microcore/pkg/dvgeolocation"
	_ "github.com/Dobryvechir/microcore/pkg/dvlicense"
	_ "github.com/Dobryvechir/microcore/pkg/dvoc"
	_ "github.com/godror/godror"
	_ "github.com/lib/pq"
)

func main() {
	dvconfig.ServerStart()
}
