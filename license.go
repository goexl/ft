package ft

import (
	"encoding/hex"
	"os"
	"strings"

	"github.com/emmansun/gmsm/sm3"
)

func (c *Client) Upload(license string, req *LicenseUploadReq, opts ...option) (rsp *LicenseUploadRsp, err error) {
	if data, re := os.ReadFile(license); nil != re {
		err = re
	} else {
		sm := sm3.New()
		sm.Write(data)
		req.HashCode = strings.ToUpper(hex.EncodeToString(sm.Sum(nil)))
	}
	if nil != err {
		return
	}

	rsp = new(LicenseUploadRsp)
	err = c.sendfile(`/api/creditInquiry/uploadLicense`, license, req, rsp, opts...)

	return
}
