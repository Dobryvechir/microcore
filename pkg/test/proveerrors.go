/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package main

func proveErrors() {
	checkErrorPref("if x=5 {x=7}","Expected ( but found x")
	checkErrorPref("for x=5 {x=7}","Expected ( but found x")
	checkErrorPref("if (x=5) x=7 else x=3","'else' should be used only inside 'if' declaration at else")
}
