package sign

import (
	"fmt"
	"sort"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/util/gconv"
)

func GetSignMd5(Data interface{}, KeyPri, Key string, Pf ...string) (string, string) {
	var pf string
	if len(Pf) > 0 {
		pf = Pf[0]
	}
	st := gconv.Map(Data)
	var newArr []string
	var outInfo string
	for k := range st {

		newArr = append(newArr, k)
	}
	sort.Strings(newArr)

	for _, v := range newArr {

		nst := gvar.New(st[v])
		if nst.IsNil() {
			continue
		}
		if pf != "" && nst.IsFloat() {
			st[v] = fmt.Sprintf(pf, st[v])
		}
		outInfo += fmt.Sprintf("%s=%v&", v, st[v])

	}

	tt := fmt.Sprintf("%s%s=%s", outInfo, KeyPri, Key)
	signMd5, _ := gmd5.Encrypt(tt)
	return tt, signMd5
}
