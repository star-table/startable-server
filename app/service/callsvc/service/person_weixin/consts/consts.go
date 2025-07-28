package callsvc

const (
	EventTypeContact          = "change_contact"
	EventChangeTypeCreateUser = "create_user"
	EventChangeTypeUpdateUser = "update_user"
	EventChangeTypeDeleteUser = "delete_user"
	EventChangeTypeCreateDept = "create_party"
	EventChangeTypeUpdateDept = "update_party"
	EventChangeTypeDeleteDept = "delete_party"
)

const (
	InfoTypePaySuccess   = "license_pay_success"
	InfoTypePayRefund    = "license_refund"
	InfoTypeRegisterCorp = "register_corp"
)

const (
	EventChangeAppAdmin = "change_app_admin"
)

const (
	EventOpenOrder        = "open_order"
	EventChangeOrder      = "change_order"
	EventPayForAppSuccess = "pay_for_app_success"
	EventRefund           = "refund"
	EventChangeEdition    = "change_editon"
)
