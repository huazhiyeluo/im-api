package schema

type ResApply struct {
	Id          uint32
	FromId      uint64
	ToId        uint64
	Type        uint32
	FromName    string
	FromIcon    string
	ToName      string
	ToIcon      string
	Reason      string
	Status      uint32
	OperateTime int64
}
