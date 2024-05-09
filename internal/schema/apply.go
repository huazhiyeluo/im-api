package schema

type ResApply struct {
	Id          uint32 `json:"id"`
	FromId      uint64 `json:"fromId"`
	ToId        uint64 `json:"toId"`
	Type        uint32 `json:"type"`
	FromName    string `json:"fromName"`
	FromIcon    string `json:"fromIcon"`
	ToName      string `json:"toName"`
	ToIcon      string `json:"toIcon"`
	Reason      string `json:"reason"`
	Status      uint32 `json:"status"`
	OperateTime int64  `json:"operateTime"`
}
