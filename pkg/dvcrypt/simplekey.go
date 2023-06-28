/***********************************************************************
MicroCore
Copyright 2020 - 2023 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvcrypt

/***

import (
	"crypto/rand"
	"log"
)


var btoaMap []byte // = getBtoaMap()
var atobMap []byte // = getAtobMap()

const headSize = 7

func GenerateHashKey(base int) []int {
	n := base << 1
	rn := make([]byte, base+n)
	_, err := rand.Read(rn)
	if err != nil {
		log.Printf("hash key cannot be secured %v", err)
	}
	key := make([]int, base<<1)
	for i := base - 1; i >= 0; i-- {
		key[i+base] = (rn[i+n] ^ i) & 255
		key[i] = i
	}
	for i := base - 1; i >= 0; i-- {
		pos = i << 1
		p = (rn[i] | (rn[i+1] << 8)) % (i + 1)
		v := key[p]
		u := key[i]
		key[p] = u
		key[i] = v
	}
	return key
}

func CheckHashKey(key []int, base int) bool {
	n := len(key)
	if n != (base << 1) {
		return false
	}
	pos := make([int]int, base)
	for i := 0; i < base; i++ {
		if !(key[i+base] >= 0 && key[i+base] <= 255) {
			return false
		}
		p := key[i]
		if !(p >= 0 && p < base) {
			return false
		}
		if _, ok := pos[p]; ok {
			return false
		}
		pos[p] = 1
	}
	return true
}

func GenerateId(size int) []byte {
	key := make([]byte, size)
	_, err := rand.Read(key)
	if err != nil {
		log.Printf("error in creating id %v", err)
	}
	for i := 0; i < size; i++ {
		b := (key[i] ^ (i * 31)) & 31
		if b < 26 {
			key[i] = b + 65
		} else {
			key[i] = b - 26 + 48
		}
	}
	return key
}

func encodeByHashKeyForString(data string, key []int, id []byte) string {
	codes := []byte(data)
	res := encodeByHashKey(codes, key, id)
	resData := convertBtoa(res)
	return resData
}

func encodeByHashKeyForUint8List(data []byte,key []int, id []byte) []byte {
  List<int> codes = utf8Dec.convert(data).codeUnits
  List<int> res = encodeByHashKey(codes, key, id)
  String resData = convertBtoa(res)
  return Uint8List.fromList(resData.codeUnits);
}

func encodeByHashKey(codes []byte, key []int, id []byte) []byte {
  int base = key.length >> 1;
  int algo = findBestAlgo(codes);
  int size = calculateLengthForAlgo(codes, algo);
  int idLen = id.length;
  int minSize = size + idLen + _headSize;
  int padding = calculatePadding(minSize, base);
  int leftPadding = padding == 0 ? 0 : _random.nextInt(padding);
  int rightPadding = padding - leftPadding;
  int totSize = minSize + padding;
  List<int> resp = List.filled(totSize, 0);
  resp[0] = leftPadding & 255;
  resp[1] = leftPadding >> 8;
  resp[2] = rightPadding & 255;
  resp[3] = rightPadding >> 8;
  resp[6] = algo;
  providePadding(resp, _headSize, leftPadding);
  int pos = _headSize + leftPadding;
  provideListCopy(resp, pos, id, idLen);
  pos += idLen;
  provideAlgoData(resp, pos, algo, codes);
  pos += size;
  providePadding(resp, pos, rightPadding);
  int sum = (-calculateSumIntList(resp)) & 0xffff;
  resp[4] = sum & 255;
  resp[5] = sum >> 8;
  sum = calculateSumIntList(resp);
  if (sum != 0) {
    throw Exception("Expected zero sum, but it is $sum");
  }
  List<int> res = encodingByHash(resp, key);
  return res;
}

func calculateSumIntList(data []int) int {
  int res = 0;
  int n = data.length;
  for (var i = 0; i < n; i++) {
    res += data[i];
    i++;
    res += data[i] << 8;
  }
  return res & 0xffff;
}

func encodingByHash(src []byte, key []int) []byte {
  int base = key.length >> 1;
  int n = src.length;
  List<int> dst = List.filled(n, 0);
  for (int p = 0; p < n; p += base) {
    for (int i = 0; i < base; i++) {
      dst[key[i] + p] = src[i + p] ^ key[i + base];
    }
  }
  return dst;
}

func getBtoaMap() []byte {
  btoa := make([]byte, 64);
  for i := 0; i < 26; i++ {
    btoa[i] = byte(i + 65);
    btoa[i + 26] = byte(i + 97);
  }
  for i := 0; i < 10; i++ {
    btoa[i + 52] = byte(i + 48);
  }
  btoa[62] = '_';
  btoa[63] = '.';
  return btoa;
}

func getAtobMap() []byte {
  atob:=make([]byte, 256)
  for i:=0;i<256;i++ {
     atob[i] = 255;
  }
  for i := 0; i < 26; i++ {
    atob[i + 65] = i;
    atob[i + 97] = i + 26;
  }
  for i := 0; i < 10; i++ {
    atob[i + 48] = i + 52;
  }
  atob['_'] = 62;
  atob['.'] = 63;
  return atob;
}

func convertBtoa(src []byte) string {
  List<String> conv = getBtoaMap();
  StringBuffer res = StringBuffer();
  int n = src.length;
  int rest = 0;
  int bits = 0;
  for (int i = 0; i < n; i++) {
    rest |= (src[i] << bits);
    res.write(conv[rest & 63]);
    bits += 2;
    rest >>= 6;
    if (bits == 6) {
      res.write(conv[rest]);
      rest = 0;
      bits = 0;
    }
  }
  if (bits > 0) {
    res.write(conv[rest]);
  }
  return res.toString();
}

func provideAlgoData(dst []byte, dstPos int, algo int, codes []int) {
  var n = codes.length;
  if (algo == 8) {
    provideListCopy(dst, dstPos, codes, n);
  } else {
    provideHalfListCopy(dst, dstPos, codes, n);
  }
}

func providePadding(buf []byte, pos int, len int) {
  int endPos = pos + len;
  for (var i = pos; i < endPos; i++) {
    buf[i] = _random.nextInt(256);
  }
}

func provideListCopy(dst []byte, dstPos int, src []int, srcSize int) {
  for (var i = 0; i < srcSize; i++) {
    dst[dstPos + i] = src[i] & 255;
  }
}

func provideHalfListCopy(dst []byte, dstPos int, src []int, srcSize int) {
  for (var i = 0; i < srcSize; i++) {
    int k = src[i];
    dst[dstPos++] = k & 255;
    dst[dstPos++] = (k >> 8) & 255;
  }
}

func findBestAlgo(codes []int) int {
  int algo = 8;
  int n = codes.length;
  for (var i = 0; i < n; i++) {
    if (codes[i] > 255) {
      return 16;
    }
  }
  return algo;
}

func calculateLengthForAlgo(codes []int, algo int) {
  int n = codes.length;
  if (algo == 16) {
    n = n * 2;
  }
  return n;
}

func calculatePadding(minSize int, base int) int {
  int rest = minSize % base;
  return rest == 0 ? 0 : base - rest;
}

func getDecodeKeyByEncodeKey(key []int) []int {
  int n = key.length;
  List<int> res = List.filled(n, 0);
  int m = n ~/ 2;
  for (var i = 0; i < m; i++) {
    res[key[i]] = i;
    res[key[i] + m] = key[i + m];
  }
  return res;
}

func convertAtob(data string, unitSize int) []int {
  List<int> conv = getAtobMap();
  int n = data.length;
  int m = n * 3 ~/ 4;
  int extra = m % unitSize;
  if (extra != 0) {
    if (extra > 1) {
      throw Exception('Extra length is too much {extra}');
    }
    m -= extra;
  }
  List<int> res = List.filled(m, 0);
  int rest = 0;
  int bits = 0;
  int pos = 0;
  List<int> codes = data.codeUnits;
  for (var i = 0; i < n; i++) {
    int cd = codes[i];
    if (cd < 32 || cd >= 128 || conv[cd] < 0) {
      throw Exception("unexpected character {cd}");
    }
    rest |= conv[cd] << bits;
    bits += 6;
    if (bits >= 8) {
      bits -= 8;
      res[pos++] = rest & 0xff;
      rest >>= 8;
      if (pos == m) {
        break;
      }
    }
  }
  return res;
}

func decodeByHashKeyForUList(data []byte, key []int, id []byte) []byte {
  String str = String.fromCharCodes(data);
  String res = decodeByHashKey(str, key, id);
  Uint8List lst = utf8Enc.convert(res);
  return lst;
}

func decodeByHashKey(data string, key []int, id []byte) (string, error) {
  int unitSize = key.length >> 1;
  List<int> resp = convertAtob(data, unitSize);
  if (resp.isEmpty) {
    throw Exception('Zero length');
  }
  List<int> res = encodingByHash(resp, key);
  if (res.length < 3 + id.length) {
    throw Exception('Too short message');
  }
  int leftPadding = res[0] | (res[1] << 8);
  int rightPadding = res[2] | (res[3] << 8);
  int sum = calculateSumIntList(res);
  if (sum != 0) {
    throw Exception('sum check error');
  }
  int algo = res[6];
  int posBegin = _headSize + leftPadding;
  int posEnd = res.length - rightPadding;
  if (posEnd < posBegin + id.length) {
    throw Exception('too short data');
  }
  if (!checkExactId(res, posBegin, id)) {
    throw Exception('Bad id');
  }
  String resData = restoreStringByAlgo(res, posBegin + id.length, posEnd, algo);
  return resData;
}

func checkExactId(res []byte, pos int, id []byte) bool {
  int n = id.length;
  for (var i = 0; i < n; i++) {
    if (res[pos + i] != (id[i] & 255)) {
      return false;
    }
  }
  return true;
}

func restoreStringByAlgo(src []int, posStart int, posEnd int, algo int) string {
  List<int> res = src;
  if (algo > 8) {
    int n = (posEnd - posStart) >> 1;
    res = List.filled(n, 0);
    for (var i = 0; i < n; i++) {
      int pos = posStart + (i << 1);
      res[i] = src[pos] | (src[pos + 1] << 8);
    }
    posStart = 0;
    posEnd = n;
  }
  return String.fromCharCodes(res, posStart, posEnd);
}

***/
