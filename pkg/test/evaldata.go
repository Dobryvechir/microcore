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
	testEvaluationSingle("", "{}", "{}", KindDvVariable)
	testEvaluationSingle("S=Good-Word", "{'a':'good',S,b:true}", "{\n  \"a\": \"good\",\n  \"S\": \"Good-Word\",\n  \"b\": true\n}", KindANY)
	testEvaluationSingle("S=Good-Word", "K={'a':'good',S,b:true};K['a']", "good", KindString)
	testEvaluationSingle("S=Good-Word", "K=['a','good',S,5,true];K[3]", "5", KindInteger)
	testEvaluationSingle("S=Good-Word", "K=['a','good',S,5,true];K[2:4]", "[\n  \"Good-Word\",\n  5\n]", KindInteger)

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

	testEvaluationSingle("S=Good-Word", "S.length", "9", KindInteger)
	testEvaluationSingle("S=Good-Word", "v=[S,S,3,4,1];v.length", "5", KindInteger)
	testEvaluationSingle("S=Good-Word", "v=[];v.length", "0", KindInteger)
	testEvaluationSingle("", "window.encodeURIComponent('+r%e@t com!$&\\'()*+,;=:[]<>\"-_.~/?@')", "%2Br%25e%40t%20com%21%24%26%27%28%29%2A%2B%2C%3B%3D%3A%5B%5D%3C%3E%22-_.~%2F%3F%40", KindString)
	testEvaluationSingle("", "encodeURIComponent('/r?i&%e@t com')", "%2Fr%3Fi%26%25e%40t%20com", KindString)
	testEvaluationSingle("", "encodeURI('/r?i&%e@t com')", "/r?i&%25e@t%20com", KindString)
	testEvaluationSingle("", "decodeURIComponent('%2Fr%3Fi%26%25e%40t%20com')", "/r?i&%e@t com", KindString)
	testEvaluationSingle("", "decodeURI('/r?i&%25e@t%20com')", "/r?i&%e@t com", KindString)
	testEvaluationSingle("", "v=Math.generateUUID();v.length", "36", KindInteger)
	testEvaluationSingle("", "v=Math.generateUUID();Math.validUUID(v)", "true", KindBoolean)
	testEvaluationSingle("S=e1234567-8978-0987-5678-0987654321Ab", "Math.validUUID(S)", "true", KindBoolean)
	testEvaluationSingle("S=h1234567-8978-0987-5678-0987654321Ab", "Math.validUUID(S)", "false", KindBoolean)
	testEvaluationSingle("", "Math.abs(3)", "3", KindFloat)
	testEvaluationSingle("", "Math.abs(-3)", "3", KindFloat)
	testEvaluationSingle("", "Math.acos(0.5)", "1.0471975511965976", KindFloat)
	testEvaluationSingle("", "Math.acosh(2)", "1.3169578969248166", KindFloat)
	testEvaluationSingle("", "Math.asin(1)", "1.5707963267948966", KindFloat)
	testEvaluationSingle("", "Math.asinh(1)", "0.881373587019543", KindFloat)
	testEvaluationSingle("", "Math.atan(1)", "0.7853981633974483", KindFloat)
	testEvaluationSingle("", "Math.atanh(0.5)", "0.5493061443340548", KindFloat)
	testEvaluationSingle("", "Math.cbrt(2)", "1.2599210498948732", KindFloat)
	testEvaluationSingle("", "Math.ceil(7.5)", "8", KindFloat)
	testEvaluationSingle("", "Math.clz32(100)", "25", KindInteger)
	testEvaluationSingle("", "Math.cos(1)", "0.5403023058681398", KindFloat)
	testEvaluationSingle("", "Math.cosh(1)", "1.5430806348152437", KindFloat)
	testEvaluationSingle("", "Math.exp(1)", "2.718281828459045", KindFloat)
	testEvaluationSingle("", "Math.expm1(1)", "1.718281828459045", KindFloat)
	testEvaluationSingle("", "Math.floor(45.95)", "45", KindFloat)
	testEvaluationSingle("", "Math.log(10)", "2.302585092994046", KindFloat)
	testEvaluationSingle("", "Math.log1p(1)", "0.6931471805599453", KindFloat)
	testEvaluationSingle("", "Math.log10(2)", "0.3010299956639812", KindFloat)
	testEvaluationSingle("", "Math.log2(3)", "1.5849625007211563", KindFloat)
	testEvaluationSingle("", "Math.round(20.5)", "21", KindFloat)
	testEvaluationSingle("", "Math.sign(-5)", "-1", KindInteger)
	testEvaluationSingle("", "Math.sin(1)", "0.8414709848078965", KindFloat)
	testEvaluationSingle("", "Math.sinh(1)", "1.1752011936438014", KindFloat)
	testEvaluationSingle("", "Math.sqrt(9)", "3", KindFloat)
	testEvaluationSingle("", "Math.tan(1)", "1.557407724654902", KindFloat)
	testEvaluationSingle("", "Math.tanh(1)", "0.7615941559557649", KindFloat)
	testEvaluationSingle("", "Math.trunc(17.834768)", "17", KindFloat)
	testEvaluationSingle("", "Math.E", "2.718281828459045", KindFloat)
	testEvaluationSingle("", "Math.LN10", "2.302585092994046", KindFloat)
	testEvaluationSingle("", "Math.LN2", "0.6931471805599453", KindFloat)
	testEvaluationSingle("", "Math.LOG10E", "0.4342944819032518", KindFloat)
	testEvaluationSingle("", "Math.LOG2E", "1.4426950408889634", KindFloat)
	testEvaluationSingle("", "Math.PI", "3.14159265359", KindFloat)
	testEvaluationSingle("", "Math.SQRT1_2", "0.7071067811865476", KindFloat)
	testEvaluationSingle("", "Math.SQRT2", "1.4142135623730951", KindFloat)
	testEvaluationSingle("S=Good evening", "S.split(' ')[0]", "Good", KindString)
	testEvaluationSingle("S=Good evening", "S.split(' ')[1]", "evening", KindString)
	testEvaluationSingle("S=Good evening", "var L=S+'var'", "Good eveningvar", KindString)
	testEvaluationSingle("S=Good evening", "const L=S+'const'", "Good eveningconst", KindString)
	testEvaluationSingle("S=Good evening", "let L=S+'let'", "Good eveninglet", KindString)
	testEvaluationSingle("S=Good evening", "var D=[];D[0]", "null", KindString)
	testEvaluationSingle("", "{'p':45}['p']", "45", KindInteger)
	testEvaluationSingle("S=Good evening", "var D=[S];D[0]", "Good evening", KindString)
	testEvaluationSingle("S=Good evening", "var D=[S,S,S];D[0]", "Good evening", KindString)
	testEvaluationSingle("S=Good evening", "var D=[S,S,S];D[0][3]", "d", KindString)
	testEvaluationSingle("S=Good evening", "var D=[{'p':45},S,S];D[0]['p']", "45", KindInteger)
	testEvaluationSingle("S=Good evening", "var D=[];D.0", "null", KindString)
	testEvaluationSingle("S=Good evening", "var D=[S];D.0", "Good evening", KindString)
	testEvaluationSingle("S=Good evening", "var D=[S,S,S];D.0", "Good evening", KindString)
	testEvaluationSingle("S=Good evening", "var D=[S,S,S];D.0.3", "d", KindString)
	testEvaluationSingle("S=Good evening", "var D=[{'p':45},S,S];D.0.p", "45", KindInteger)
	testEvaluationSingle("S=Good", "D=null;(D && D.0 && D.0.value || {})[S] || []", "[]", KindString)
	testEvaluationSingle("S=Good", "D=[];(D.0 && D.0.value || {})[S] || []", "[]", KindString)
	testEvaluationSingle("S=Good", "D=[{'p':S}];(D.0 && D.0.value || {})[S] || []", "[]", KindString)
	testEvaluationSingle("I=3", "D=[{'value':'Good'}];(D.0 && D.0.value || {})[I] || []", "d", KindString)
	testEvaluationSingle("", "return", "null", KindString)
	testEvaluationSingle("I=3", "return;'ABC'", "null", KindString)
	testEvaluationSingle("I=3", "G=9;return G;'ABC'", "9", KindString)
	testEvaluationSingle("", "return 10", "10", KindInteger)
	testEvaluationSingle("I=0", "D=[{'value':'Good'}];return D[I].value + ' Evening';'ABC'", "Good Evening", KindString)
	testEvaluationSingle("I=3", "D=[{'value':'Good'}];return (2+I)*(I+5);'ABC'", "40", KindInteger)
	testEvaluationSingle("I=3", "D=[{'value':'Good'}];return;'ABC'", "null", KindString)
	testEvaluationSingle("I=3", "D=[{'value':'Good'}];return 24;'ABC'", "24", KindInteger)
	testEvaluationSingle("I=3", "D=[{'value':'Good'}];return D[0].value[I];'ABC'", "d", KindString)
	testEvaluationSingle("", "A=()=>{return 10};A()", "10", KindInteger)
	testEvaluationSingle("", "A=(first,second)=>{return first*second+1};A(5,4)", "21", KindInteger)
	testEvaluationSingle("I=3", "A=(val)=>{return val*val+I};A(3)+A(4)", "31", KindInteger)
	testEvaluationSingle("", "S=2;val=5;A=(val)=>{S=200;return val*val};A(val-4)+val+S", "206", KindInteger)
	testEvaluationSingle("", "S=[1,2,3,4];initial=5;S.reduce((previous,current)=>{previous+current},initial)", "15", KindInteger)
	showResume()
}
