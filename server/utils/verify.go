package utils

var (
	IdVerify               = Rules{"ID": []string{NotEmpty()}}
	ApiVerify              = Rules{"Path": {NotEmpty()}, "Description": {NotEmpty()}, "ApiGroup": {NotEmpty()}, "Method": {NotEmpty()}}
	MenuVerify             = Rules{"Path": {NotEmpty()}, "ParentId": {NotEmpty()}, "Name": {NotEmpty()}, "Component": {NotEmpty()}, "Sort": {Ge("0")}}
	MenuMetaVerify         = Rules{"Title": {NotEmpty()}}
	LoginVerify            = Rules{"CaptchaId": {NotEmpty()}, "UserName": {NotEmpty()}, "Password": {NotEmpty()}}
	RegisterVerify         = Rules{"UserName": {NotEmpty()}, "NickName": {NotEmpty()}, "Password": {NotEmpty()}, "AuthorityId": {NotEmpty()}}
	PageInfoVerify         = Rules{"Page": {NotEmpty()}, "PageSize": {NotEmpty()}}
	StateInfoVerify        = Rules{"Page": {NotEmpty()}, "PageSize": {NotEmpty()}}
	CustomerVerify         = Rules{"CustomerName": {NotEmpty()}, "CustomerPhoneData": {NotEmpty()}}
	AutoCodeVerify         = Rules{"Abbreviation": {NotEmpty()}, "StructName": {NotEmpty()}, "PackageName": {NotEmpty()}, "Fields": {NotEmpty()}}
	AutoPackageVerify      = Rules{"PackageName": {NotEmpty()}}
	AuthorityVerify        = Rules{"AuthorityId": {NotEmpty()}, "AuthorityName": {NotEmpty()}}
	AuthorityIdVerify      = Rules{"AuthorityId": {NotEmpty()}}
	UserIdVerify           = Rules{"UserId": {NotEmpty()}}
	OldAuthorityVerify     = Rules{"OldAuthorityId": {NotEmpty()}}
	ChangePasswordVerify   = Rules{"Password": {NotEmpty()}, "NewPassword": {NotEmpty()}}
	SetUserAuthorityVerify = Rules{"AuthorityId": {NotEmpty()}}
	ComboVerify            = Rules{"ComboName": {NotEmpty()}, "ComboType": {NotEmpty()}, "ComboPrice": {Gt("0")}, "Amount": {Ge("0")}}
	MemberVerify           = Rules{"CardId": {NotEmpty()}, "Telephone": {Eq("11")}, "UserName": {NotEmpty()}}
	ConsumeVerify          = Rules{"CardID": {NotEmpty()}, "ComboId": {NotEmpty()}, "Deadline": {NotEmpty()}, "State": {NotEmpty()}, "Number": {NotEmpty()}}
	CardVerify             = Rules{"OnlyId": {NotEmpty()}}
	WxRegisterVerify       = Rules{"OpenID": {NotEmpty()}, "Code": {NotEmpty()}}
)
