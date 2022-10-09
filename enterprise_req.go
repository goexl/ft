package ft

// EnterpriseReq 企业信息上传请求
type EnterpriseReq struct {
	// 地方平台编码
	PlatformId int64 `json:"platformId"`
	// 统一社会信用代码
	Code string `json:"uniscId"`
	// 企业名称
	Name string `json:"enterpriseName"`
	// 是否实名认证
	RealNameCertification string `json:"realNameCertification"`
	// 实名认证人姓名
	OperatorName string `json:"operatorName"`
	// 实名认证人身份
	OperatorIdentity string `json:"operatorIdentity"`
	// 实名认证人身份证号
	OperatorIdCardNo string `json:"operatorIdCardNo"`
	// 是否通过平台获得贷款
	Loaned string `json:"loaned"`
	// 营业执照住所
	Address string `json:"address"`
	// 企业所属行业
	Industry string `json:"industry"`
	// 企业所在省
	Province string `json:"province"`
	// 企业所在市
	City string `json:"city"`
	// 企业所在区
	Area string `json:"area"`
	// 注册资本，单位万元
	RegisteredCapital int64 `json:"registeredCapital"`
	// 经营范围
	BusinessScope string `json:"businessScope"`
	// 经营期限类型
	OperatingTimeLimitType string `json:"operatingTimeLimitType"`
	// 营业期限开始日期
	OperatingTimeLimitDateBegin string `json:"operatingTimeLimitDateBegin"`
	// 营业期限结束日期
	OperatingTimeLimitDateEnd string `json:"operatingTimeLimitDateEnd"`
	// 核准日期
	ApprovalDate string `json:"approvalDate"`
	// 入驻时间
	SettlingTime string `json:"settlingTime"`
	// 外部系统编号
	ExternalSystemId string `json:"externalSystemId"`
}
