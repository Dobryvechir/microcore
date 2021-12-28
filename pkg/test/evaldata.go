package main

func testEvaluation() {
	initEnvironment()
	testEvaluationSingle("", "2+2","4", KindInteger)
	testEvaluationSingle("", "2+2+2","6", KindInteger)
	testEvaluationSingle("", "20*20*20","8000", KindInteger)
	testEvaluationSingle("S=5", "S+1","51", KindString)
	testEvaluationSingle("S=5","+S+1", "6", KindInteger)
	showResume()
}


