/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

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
	testEvaluationSingle("I=3", "I+=I*100+I", "306", KindANY)

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
	testEvaluationSingle("", "A=4;for(I=1;I<2;I++){A=A+I};A", "5", KindInteger)
	testEvaluationSingle("", "A={'a':'b','c':'d'};B='';for(xx in A){B+=xx};B", "ac", KindString)
	testEvaluationSingle("", "A={'a':'b','c':'d'};B='';for(xx of A){B+=xx};B", "bd", KindString)
	testEvaluationSingle("I=3", "B=0;if (I<4) {B=5};B", "5", KindInteger)
	testEvaluationSingle("I=3", "B=0;if (I<3) {B=4} else {B=5};B", "5", KindInteger)
	testEvaluationSingle("", "A={'a':'b','c':'d'};B='';for(xx in A){B+=xx;if (xx=='a') {break}};B", "a", KindString)
	testEvaluationSingle("", "A={'a':'b','c':'d'};B='';for(xx in A){if (xx=='a') {continue};B+=xx;};B", "c", KindString)
	testEvaluationSingle("I=3", "B=0;if (I==3) B=4; else if (I<4) B=5; else B=6;B", "4", KindInteger)
	testEvaluationSingle("I=0", "B=0;if (I==3) B=4; else if (I<4) B=5; else B=6;B", "5", KindInteger)
	testEvaluationSingle("I=4", "B=0;if (I==3) B=4; else if (I<4) B=5; else B=6;B", "6", KindInteger)
	testEvaluationSingle("", "B={'a':{'b':'c'},'d':'e'};B.a=5;B['a']", "5", KindInteger)
	testEvaluationSingle("", "B={'a':{'b':'c'},'d':'e'};B.a.b={'c':7};B.a.b.c", "7", KindInteger)
	testEvaluationSingle("", "B={'a':{'b':'c'},'d':'e'};B.f=5;B['f']", "5", KindInteger)
	testEvaluationSingle("", "B={'a':{'b':'c'},'d':'e'};B.a.f={'c':7};B.a.f.c", "7", KindInteger)
	testEvaluationSingle("", "B=['a',{'b':'c'},'d','e'];B.2=5;B['2']", "5", KindInteger)
	testEvaluationSingle("", "B={'a':['b','c'],'d':'e'};B.a.1={'c':7};B.a.1.c", "7", KindInteger)
	testEvaluationSingle("", "B=['a',{'b':'c'},'d','e'];B[2]=5;B['2']", "5", KindInteger)
	testEvaluationSingle("I=1", "B={'a':['b','c'],'d':'e'};B.a[I]={'c':7};B.a[I].c", "7", KindInteger)
	testEvaluationSingle("", "B=['a',{'b':4},'d','e'];B[1].b++;B[1].b", "5", KindInteger)
	testEvaluationSingle("", "B=['a',{'b':4},'d','e'];++B[1].b;B[1].b", "5", KindInteger)
	testEvaluationSingle("", "B=['a',{'b':4},'d','e'];B[1].b+=2;B[1].b", "6", KindInteger)
	testEvaluationSingle("", "B=['a',7,'d','e'];B[1]--;--B[1];++B[1]", "6", KindInteger)
	testEvaluationSingle("I=5", "B=7;B+=I;B", "12", KindInteger)
	testEvaluationSingle("I=5", "B=7;B*=I;B", "35", KindInteger)
	testEvaluationSingle("I=5", "B=7;B-=I;B", "2", KindInteger)
	testEvaluationSingle("I=5", "B=35;B/=I;B", "7", KindInteger)
	testEvaluationSingle("I=5", "B=7;B%=I;B", "2", KindInteger)
	testEvaluationSingle("I=3", "B=7;B**=I;B", "343", KindInteger)
	testEvaluationSingle("I=3", "B=7;B<<=I;B", "56", KindInteger)
	testEvaluationSingle("I=3", "B=27;B>>=I;B", "3", KindInteger)
	testEvaluationSingle("I=3", "B=-27;B>>=I;B", "-4", KindInteger)
	testEvaluationSingle("I=3", "B=27;B>>>=I;B", "3", KindInteger)
	testEvaluationSingle("I=5", "B=3;B&=I;B", "1", KindInteger)
	testEvaluationSingle("I=5", "B=7;B^=I;B", "2", KindInteger)
	testEvaluationSingle("I=13", "B=7;B|=I;B", "15", KindInteger)
	testEvaluationSingle("I=5", "B=7;B&&=I;B", "5", KindInteger)
	testEvaluationSingle("I=5", "B=0;B&&=I;B", "0", KindInteger)
	testEvaluationSingle("I=5", "B=7;B||=I;B", "7", KindInteger)
	testEvaluationSingle("I=5", "B='';B||=I;B", "5", KindInteger)
	testEvaluationSingle("I=5", "B=0;B^=I;B", "5", KindInteger)
	testEvaluationSingle("I=5", "B=7;B??=I;B", "7", KindInteger)
	testEvaluationSingle("I=5", "B=0;B??=I;B", "0", KindInteger)
	testEvaluationSingle("I=5", "B=null;B??=I;B", "5", KindInteger)
	testEvaluationSingle("", "B=['a',{'b':4},'d','e'];delete B[1].b;B[1].b", "null", KindANY)
	testEvaluationSingle("", "B=['a',{'b':4},'d','e'];B.push(5,7);B[4]", "5", KindANY)
	testEvaluationSingle("", "B=['a',{'b':4},'d','e'];B.push(5,7);B[4]", "5", KindANY)
	testEvaluationSingle("", "B=['a',{'b':4},'d','e'];B.push(5,7);B[4]", "5", KindANY)
	testEvaluationSingle("", "B=['a',{'b':4},'d','e'];C=B.slice();C.length", "4", KindANY)
	testEvaluationSingle("", "B=['a',{'b':4},'d','e'];C=B.slice(2,-1);C.length", "1", KindANY)
	testEvaluationSingle("", "B=['a',{'b':4},'d','e'];C=B.slice(2,-1);C[0]", "d", KindANY)
	testEvaluationSingle("", "B=['a',{'b':4},'d','e'];C=B.fill('ab');C[0]", "ab", KindANY)
	testEvaluationSingle("", "B=[20,{'b':4},'d','e'];C=B.slice(1);C.concat(B,B,B).length", "15", KindANY)
	testEvaluationSingle("", "B=[20,{'b':4},'d','e'];C=B.slice(1);C.concat(B,B,B);C[4]", "null", KindANY)
	testEvaluationSingle("", "B=[20,{'b':4},'d','e'];C=B.slice(1);D=C.concat(B,B,B);D[3]", "20", KindANY)
	testEvaluationSingle("", "B=[20,{'b':4},'d','e'];B.copyWithin(1,0,1);B[1]", "20", KindANY)
	testEvaluationSingle("", "B=[20,5,'d','e'].entries();C=B.next().value[1];D=B.next().value[1];C+D", "25", KindANY)
	testEvaluationSingle("", "B=[1,5,3,4].every((val,ind)=>{val==ind+1})", "false", KindANY)
	testEvaluationSingle("", "B=[1,2,3,4].every((val,ind)=>{val==ind+1})", "true", KindANY)
	testEvaluationSingle("", "B=[1,2,3,4].fill(25,2,3);B[1]+B[2]+B[3]", "31", KindANY)
	testEvaluationSingle("", "B=[1,2,3,4].fill(25,2);B[1]+B[2]+B[3]", "52", KindANY)
	testEvaluationSingle("", "B=[1,2,3,4].fill(5);B[0]+B[1]+B[2]+B[3]", "20", KindANY)
	testEvaluationSingle("", "B=[11,22,3,4].find((val,ind)=>{val==ind+1})", "3", KindANY)
	testEvaluationSingle("", "B=[11,22,33,44].find((val,ind)=>{val==ind+1})", "null", KindANY)
	testEvaluationSingle("", "B=[11,22,3,4].findIndex((val,ind)=>{val==ind+1})", "2", KindANY)
	testEvaluationSingle("", "B=[11,22,33,44].findIndex((val,ind)=>{val==ind+1})", "-1", KindANY)
	testEvaluationSingle("", "B=[11,22,3,4].findLast((val,ind)=>{val==ind+1})", "4", KindANY)
	testEvaluationSingle("", "B=[11,22,33,44].findLast((val,ind)=>{val==ind+1})", "null", KindANY)
	testEvaluationSingle("", "B=[11,22,3,4].findLastIndex((val,ind)=>{val==ind+1})", "4", KindANY)
	testEvaluationSingle("", "B=[11,22,33,44].findLastIndex((val,ind)=>{val==ind+1})", "-1", KindANY)
	testEvaluationSingle("", "B=[1,null,2,[3,4,[5,6]]].flat();B.length", "5", KindANY)
	testEvaluationSingle("", "B=[1,null,2,[3,4,[5,6]]].flat(2);B.length", "6", KindANY)
	testEvaluationSingle("", "B=[1,null,2,[3,4,[5,6]]].flat(2);B[5]", "6", KindANY)
	testEvaluationSingle("", "B=[1,2,[3,4,5,6]].flatMap((val,ind)=>{if (Array.isArray(val)){val[0]=ind}else{val-ind} });B.length*1000+B[0]*100+B[1]*10+B[2]", "3112", KindANY)
	testEvaluationSingle("", "B=[1,2,[3,4,5,6]].flatMap((val,ind)=>{if (Array.isArray(val)){val[0]=ind;return val;}else{return val-ind} });B.length*1000+B[1]*100+B[2]*10+B[3]", "6124", KindANY)
	testEvaluationSingle("", "B=3;[1,2,3,4].forEach((val,ind)=>{B+=val*100+ind});B", "1009", KindANY)
	testEvaluationSingle("", "B=Array.from('good');B[1]=='o' && B[3]=='d'", "true", KindANY)
	testEvaluationSingle("", "B=Array.from([1,2,3,4],(val,ind)=>{val-ind+10});B[1]+B[3]", "22", KindANY)
	testEvaluationSingle("", "B=[1,2,3,4].includes(3)", "true", KindANY)
	testEvaluationSingle("", "B=[1,2,3,4].includes(-1)", "false", KindANY)
	testEvaluationSingle("", "B=[1,2,3,3].indexOf(3)", "2", KindANY)
	testEvaluationSingle("", "B=[1,2,3,3].indexOf(4)", "-1", KindANY)
	testEvaluationSingle("", "B=[1,2,3,3].lastIndexOf(3)", "3", KindANY)
	testEvaluationSingle("", "B=[1,2,3,3].lastIndexOf(4)", "-1", KindANY)
	testEvaluationSingle("", "B=[1,2,3,4];Array.isArray(B)", "true", KindANY)
	testEvaluationSingle("", "B=[1,2,3,4].length;Array.isArray(B)", "false", KindANY)
	testEvaluationSingle("", "B=['1a','2','3'].join()", "1a,2,3", KindANY)
	testEvaluationSingle("", "B=['1a','2','3'].join('')", "1a23", KindANY)
	testEvaluationSingle("", "B=['1a','2','3'].join('-')", "1a-2-3", KindANY)
	testEvaluationSingle("", "B=[1,2,3,4].keys()[3]", "3", KindANY)
	testEvaluationSingle("", "B=[1,2,3,4].map((val,ind)=>{val+ind+100});B[3]", "107", KindANY)
	testEvaluationSingle("", "B=Array.of(1,2,3,4);B[2]", "3", KindANY)
	testEvaluationSingle("", "B=[1,2,3,4];B.pop();B.pop();B[B.length-1]", "2", KindANY)
	testEvaluationSingle("", "B=[1,2,3,4];B.pop()", "4", KindANY)
	testEvaluationSingle("", "B=[[0, 1], [2, 3], [4, 5]];C=B.reduceRight((accumulator, currentValue) => {accumulator.concat(currentValue)},[]);C[0]", "4", KindANY)
	testEvaluationSingle("", "B=[1,2,3,4].revert();B[0]", "4", KindANY)
	testEvaluationSingle("", "B=[1,2,3,4];B.shift()", "1", KindANY)
	testEvaluationSingle("", "B=[1,2,3,4];B.shift();B.length", "3", KindANY)
	testEvaluationSingle("", "B=[1,2,3,4];B.some((val,ind)=>{val==3})", "true", KindANY)
	testEvaluationSingle("", "B=[1,2,3,4];B.some((val,ind)=>{val==5})", "false", KindANY)
	testEvaluationSingle("", "B=[1,4,3,2];B.sort();B[1]", "2", KindANY)
	testEvaluationSingle("", "B=[1,4,3,2];B.sort((a,b)=>{b-a});B[1]", "3", KindANY)
	testEvaluationSingle("", "B=[1,2,3,4];B.splice(1,2);B[1]", "4", KindANY)
	testEvaluationSingle("", "B=[1,2,3,4];B.splice(1,2,20,30,50);B.length", "5", KindANY)
	testEvaluationSingle("", "B=[1,2,3,4];B.splice(1,2,20,30,50);B[3]", "50", KindANY)
	testEvaluationSingle("", "B=[1,2];B.toLocaleString()", "1,2", KindANY)
	testEvaluationSingle("", "B=[1,2];B.toString()", "1,2", KindANY)
	testEvaluationSingle("", "B=[1,2];B.unshift(4,5);B.length", "4", KindANY)
	testEvaluationSingle("", "B=[1,2];B.unshift(4,5);B[0]", "4", KindANY)
	testEvaluationSingle("", "B=[11,22].values();C=B.next().value;D=B.next().value;C+D", "33", KindANY)
	testEvaluationSingle("", "B=[11,22].values();C=B.next().value;D=B.next().value;B.next().done", "true", KindANY)
	testEvaluationSingle("", "B=[11,12,33].filter((x)=>{x%11==0});B.length", "2", KindANY)
	testEvaluationSingle("", "'Diet'.charAt(2)", "e", KindString)
	testEvaluationSingle("", "'Diet'.charCodeAt(2)", "101", KindInteger)
	testEvaluationSingle("", "'Diet'.codePointAt(2)", "101", KindInteger)
	testEvaluationSingle("", "B=Array.from('ще').length", "2", KindANY)
	testEvaluationSingle("", "B=Array.from('ще');B[0]=='щ' && B[1]=='е'", "true", KindANY)
	testEvaluationSingle("", "B=Array(7);B.length", "7", KindANY)
	testEvaluationSingle("", "B=Array(7,3);B.length", "2", KindANY)
	testEvaluationSingle("", "B=Array(7,3);B[0]+B[1]", "10", KindANY)
	testEvaluationSingle("", "B='Good';B.concat(' ','Evening')", "Good Evening", KindANY)
	testEvaluationSingle("", "B='Good';B.endsWith('oo')", "false", KindANY)
	testEvaluationSingle("", "B='Good';B.endsWith('od')", "true", KindANY)
	testEvaluationSingle("", "String.fromCharCode('D'.charCodeAt(0),'v'.charCodeAt(0))", "Dv", KindANY)
	testEvaluationSingle("", "'ж'.codePointAt(0)", "1078", KindANY)
	testEvaluationSingle("", "String.fromCodePoint('у'.codePointAt(0),'ж'.codePointAt(0))", "уж", KindANY)
	testEvaluationSingle("", "'Добрий'.includes('бр')", "true", KindANY)
	testEvaluationSingle("", "'Добрий'.includes('br')", "false", KindANY)
	testEvaluationSingle("", "'Dobrzy'.indexOf('rzy')", "3", KindANY)
	testEvaluationSingle("", "'Dobrzy'.indexOf('o',3)", "-1", KindANY)
	testEvaluationSingle("", "'Dobrzo'.lastIndexOf('o')", "5", KindANY)
	testEvaluationSingle("", "'Dobrzo'.lastIndexOf('o',3)", "1", KindANY)
	testEvaluationSingle("", "'Dobrzo'.lastIndexOf('z',3)", "-1", KindANY)
	testEvaluationSingle("", "'Dobrzo'.localeCompare('Dobry')", "1", KindANY)
	testEvaluationSingle("", "V=new RegExp('\\\\w','g');V.dotAll", "false", KindANY)
	testEvaluationSingle("", "V=new RegExp('\\\\w','gs');V.dotAll", "true", KindANY)
	testEvaluationSingle("", "V=new RegExp('\\\\w','g');V.flags", "g", KindANY)
	testEvaluationSingle("", "V=new RegExp('\\\\w');V.global", "false", KindANY)
	testEvaluationSingle("", "V=new RegExp('\\\\w','gs');V.global", "true", KindANY)
	testEvaluationSingle("", "V=new RegExp('\\\\w','g');V.sticky", "false", KindANY)
	testEvaluationSingle("", "V=new RegExp('\\\\w','gys');V.sticky", "true", KindANY)
	testEvaluationSingle("", "V=new RegExp('\\\\w','g');V.hasIndices", "false", KindANY)
	testEvaluationSingle("", "V=new RegExp('\\\\w','gysd');V.hasIndices", "true", KindANY)
	testEvaluationSingle("", "V=new RegExp('\\\\w','g');V.multiline", "false", KindANY)
	testEvaluationSingle("", "V=new RegExp('\\\\w','gysium');V.multiline", "true", KindANY)
	testEvaluationSingle("", "V=new RegExp('\\\\w','g');V.ignoreCase", "false", KindANY)
	testEvaluationSingle("", "V=new RegExp('\\\\w','gysi');V.ignoreCase", "true", KindANY)
	testEvaluationSingle("", "V=new RegExp('\\\\w','g');V.unicode", "false", KindANY)
	testEvaluationSingle("", "V=new RegExp('\\\\w','gysiu');V.unicode", "true", KindANY)
	testEvaluationSingle("", "V=new RegExp('[w]','gysi');V.source", "[w]", KindANY)
	testEvaluationSingle("", "V=new RegExp('oo','g');T='good good day';V.test(T)", "true", KindANY)
	testEvaluationSingle("", "V=new RegExp('foo','g');T='good good day';V.test(T)", "false", KindANY)
	testEvaluationSingle("", "V=new RegExp('oo','g');T='good good day';V.test(T);V.lastIndex", "3", KindANY)
	testEvaluationSingle("", "V=new RegExp('oo','g');T='good good day';V.test(T);V.test(T);V.lastIndex", "7", KindANY)
	testEvaluationSingle("", "V=new RegExp('oo','g');T='good good day';V.test(T);V.test(T)", "true", KindANY)
	testEvaluationSingle("", "V=new RegExp('oo','g');T='good good day';V.test(T);V.test(T);V.test(T)", "false", KindANY)
	testEvaluationSingle("", "V=new RegExp('ech','y');T='#vechir#';(V.test(T)?1000:0)+V.lastIndex", "0", KindANY)
	testEvaluationSingle("", "V=new RegExp('ech','y');T='#vechir#';V.lastIndex=2;(V.test(T)?1000:0)+V.lastIndex", "1005", KindANY)
	testEvaluationSingle("", "V=new RegExp('ech','y');T='#vechir#';V.lastIndex=5;(V.test(T)?1000:0)+V.lastIndex", "0", KindANY)
	testEvaluationSingle("", "V=new RegExp('ech','y');T='#vechir#';V.lastIndex=1;(V.test(T)?1000:0)+V.lastIndex", "0", KindANY)
	testEvaluationSingle("", "V=new RegExp('^ech','my');T='#\\nechir#';V.lastIndex=2;(V.test(T)?1000:0)+V.lastIndex", "1005", KindANY)
	testEvaluationSingle("", "V=new RegExp('^ech','my');T='#vechir#';V.lastIndex=2;(V.test(T)?1000:0)+V.lastIndex", "0", KindANY)
	testEvaluationSingle("", "V=new RegExp('^ech','y');T='#\\nechir#';V.lastIndex=2;(V.test(T)?1000:0)+V.lastIndex", "0", KindANY)
	testEvaluationSingle("", "V=new RegExp('^ech','y');T='#vechir#';V.lastIndex=2;(V.test(T)?1000:0)+V.lastIndex", "0", KindANY)

	proveErrors()
	showResume()
}
