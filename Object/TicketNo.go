package Object

import (
	"encoding/json"
	"fmt"
	"github.com/Deansquirrel/goZlDianzqOfferTicket/global"
	"io/ioutil"
	"net/http"
)

type TicketNo struct {
}

func GetTktNoMulti(num int) ([]string, error) {
	rUrl := fmt.Sprintf("%s/Api/Number/GetTktNo_Multi?workerId=%d&nums=%d", global.Config.TotalConfig.SnoServer, global.Config.TotalConfig.SnoWorkerId, num)
	resp, err := http.Get(rUrl)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			global.MyLog(err.Error())
		}
	}()
	body, err := ioutil.ReadAll(resp.Body)
	var r []string
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func GetSno(prefix string) (string, error) {
	rUrl := fmt.Sprintf("%s/Api/Number/GetSno?workerId=%d&&prefix=%s", global.Config.TotalConfig.SnoServer, global.Config.TotalConfig.SnoWorkerId, prefix)
	resp, err := http.Get(rUrl)
	if err != nil {
		return "", err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			global.MyLog(err.Error())
		}
	}()
	sno, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(sno), nil
}
