package system

type SysAuthorityBtn struct {
	AuthorityId      int            `gorm:"comment:角色ID"`
	SysMenuID        int            `gorm:"comment:菜单ID"`
	SysBaseMenuBtnID int            `gorm:"comment:菜单按钮ID"`
	SysBaseMenuBtn   SysBaseMenuBtn ` gorm:"comment:按钮详情"`
}
