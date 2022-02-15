package casbin

//CasbinModel 权限结构
type casbinModel struct {
	ID    int    `gorm:"primaryKey;autoIncrement"`
	Ptype string `gorm:"size:68;uniqueIndex:unique_index"`
	V0    string `gorm:"size:68;uniqueIndex:unique_index"`
	V1    string `gorm:"size:68;uniqueIndex:unique_index"`
	V2    string `gorm:"size:68;uniqueIndex:unique_index"`
	V3    string `gorm:"size:68;uniqueIndex:unique_index"`
	V4    string `gorm:"size:68;uniqueIndex:unique_index"`
	V5    string `gorm:"size:68;uniqueIndex:unique_index"`
}
