package main

func testEvaluation() {
	initEnvironment()
	testEvaluationSingle("", "2+2", "4", KindInteger)
	testEvaluationSingle("", "2+2+2", "6", KindInteger)
	testEvaluationSingle("", "20*20*20", "8000", KindInteger)
	testEvaluationSingle("", "6 / 3", "2", KindInteger)
	testEvaluationSingle("", "7/2", "3.5", KindFloat)
	testEvaluationSingle("", "2+3*10+4", "36", KindInteger)
	testEvaluationSingle("", "3-2", "1", KindInteger)
	testEvaluationSingle("", "2 - 3", "-1", KindInteger)
	testEvaluationSingle("", "-1-3", "-4", KindInteger)
	testEvaluationSingle("", "-1 - 2 * 3", "-7", KindInteger)
	testEvaluationSingle("", "(-1 - 2)*3", "-9", KindInteger)
	testEvaluationSingle("", "5%2", "1", KindInteger)
	testEvaluationSingle("", "-10 % 3", "-1", KindInteger)
	testEvaluationSingle("", "2**3", "8", KindInteger)
	testEvaluationSingle("", "2 ** 10", "1024", KindInteger)
	testEvaluationSingle("", "3<<5", "96", KindInteger)
	testEvaluationSingle("", "21>>2", "5", KindInteger)
	testEvaluationSingle("", "21>>>2", "5", KindInteger)
	testEvaluationSingle("", "-21>>>2", "4611686018427387898", KindInteger)
	testEvaluationSingle("", "-21>>2", "-6", KindInteger)
	testEvaluationSingle("", "7 & 4", "4", KindInteger)
	testEvaluationSingle("", "5 | 6", "7", KindInteger)
	testEvaluationSingle("", "5 ^ 6", "3", KindInteger)
	testEvaluationSingle("", "~5", "-6", KindInteger)
	testEvaluationSingle("", "~-6", "5", KindInteger)
	testEvaluationSingle("", "5 == 6", "false", KindBoolean)
	testEvaluationSingle("", "5 == 5", "true", KindBoolean)
	testEvaluationSingle("", "5 != 6", "true", KindBoolean)
	testEvaluationSingle("", "false && true", "false", KindBoolean)
	testEvaluationSingle("", "false || true", "true", KindBoolean)
	testEvaluationSingle("", "'5' == '6' && '6'=='6'", "false", KindBoolean)
	testEvaluationSingle("", "'AB'+'CD'", "ABCD", KindString)
	testEvaluationSingle("", "'AB'+'CD'=='ABCD' && 6==6", "true", KindBoolean)
	testEvaluationSingle("", "2<3", "true", KindBoolean)
	testEvaluationSingle("", "5>3", "true", KindBoolean)
	testEvaluationSingle("", "5 **2<=25.1", "true", KindBoolean)
	testEvaluationSingle("", "5 **2>=24.9", "true", KindBoolean)
	testEvaluationSingle("", "5 **(1+1)<=25.1", "true", KindBoolean)
	testEvaluationSingle("", "5 **(1+1)>=24.9", "true", KindBoolean)
	testEvaluationSingle("", "!1", "false", KindBoolean)
	testEvaluationSingle("", "'CAT' && 'DOG'", "DOG", KindString)
	testEvaluationSingle("", "'CAT' || 'DOG'", "CAT", KindString)
	testEvaluationSingle("", "0 && 'DOG'", "0", KindString)
	testEvaluationSingle("", "1+2,3+4", "7", KindInteger)
	testEvaluationSingle("", "1+2;3+4", "7", KindInteger)
	testEvaluationSingle("", "1+2,3+4;5+6;7+8,9+10", "19", KindInteger)
	testEvaluationSingle("", "1>2 ? 'A' : 'B'", "B", KindString)
	testEvaluationSingle("", "1<2 ? 'A' : 'B'", "A", KindString)

	testEvaluationSingle("S=5", "S+1", "51", KindString)
	testEvaluationSingle("S=5", "+S+1", "6", KindInteger)
	testEvaluationSingle("S=5", "S-1", "4", KindInteger)
	testEvaluationSingle("S=5", "-S*2", "-10", KindInteger)
	testEvaluationSingle("S=5", "-S-1", "-6", KindInteger)
	testEvaluationSingle("S=5", "--S+6", "10", KindInteger)
	testEvaluationSingle("S=5", "++S-2", "4", KindInteger)
	testEvaluationSingle("S=5", "S++ + +4", "9", KindInteger)
	testEvaluationSingle("S=5", "S++ + S", "11", KindInteger)
	testEvaluationSingle("S=5", "++S + S", "12", KindInteger)
	testEvaluationSingle("S=5", "S-- + +4", "9", KindInteger)
	testEvaluationSingle("S=5", "S-- + S", "9", KindInteger)
	testEvaluationSingle("S=5;T=4;U=3", "S-- + T-- + U--", "12", KindInteger)
	testEvaluationSingle("S=5", "S--,S+10", "14", KindInteger)
	testEvaluationSingle("S=5", "I=15,+S+I", "20", KindInteger)
	testEvaluationSingle("S=5;I=3", "I+=12, +S+I", "20", KindInteger)
	testEvaluationSingle("S=ABCD", "S[2]", "C", KindString)
	testEvaluationSingle("S=ABCD", "S[1:3]", "BC", KindString)
	testEvaluationSingle("V={\"AB\":\"DOG\",\"CD\":\"CAT\"}", "V['AB']", "DOG", KindString)
	testEvaluationSingle("V={\"AB\":\"DOG\",\"CD\":\"CAT\"}", "'AB' IN V", "true", KindBoolean)
	testEvaluationSingle("V={\"AB\":\"DOG\",\"CD\":\"CAT\"}", "'BC' IN V", "false", KindBoolean)
	testEvaluationSingle("v={\"AB\":\"DOG\",\"CD\":\"CAT\"}", "v.AB + \" barks at \"+ v.CD", "DOG barks at CAT", KindString)
	testEvaluationSingle("i=2", "i<10 ? (i<7? 1: 2) : 3", "1", KindInteger)
	testEvaluationSingle("i=8", "i<10 ? (i<7? 1: 2) : 3", "2", KindInteger)
	testEvaluationSingle("i=12", "i<10 ? (i<7? 1: 2) : 3", "3", KindInteger)
	testEvaluationSingle("i=2", "i<10 ? i<7? 1: 2 : 3", "1", KindInteger)
	testEvaluationSingle("i=8", "i<10 ? i<7? 1: 2 : 3", "2", KindInteger)
	testEvaluationSingle("i=12", "i<10 ? i<7? 1: 2 : 3", "3", KindInteger)
	showResume()
}
