package consts

const (
	ScenePolicyStatusEnabled = 1

	ACLTypePublicRead = 1
	ACLTypePrivate    = 2

	AssetRoleOriginal = 1
	AssetRoleDerived  = 2

	DerivedKindUnspecified = 0

	RiskStatusPending = 1
	RiskStatusPassed  = 2
	RiskStatusRejected = 3

	ProcessStatusPending       = 1
	ProcessStatusProcessing    = 2
	ProcessStatusSuccess       = 3
	ProcessStatusFailed        = 4
	ProcessStatusPartialSuccess = 5

	BindingActive   = 1
	BindingInactive = 0

	DefaultIssueReadTTLSeconds = 300
	MaxIssueReadTTLSeconds     = 3600
)
