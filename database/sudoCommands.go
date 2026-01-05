package database

type Permission uint64

const (
	PermissionDisconnectUsers Permission = 1 << iota
	PermissionMuteUsers
)

type SudoCommand struct {
	ID          string `gorm:"index"`
	TargetUser  string `gorm:"index"`
	Permissions uint64
}

func HasSudoCommandPermission(actorID, targetID string, p Permission) bool {
	var exists bool
	if err := GormDB.Raw("SELECT COUNT(1) FROM sudo_commands WHERE id = ? AND target_user = ? AND (permissions & ?) = ?", actorID, targetID, uint64(p), uint64(p)).Scan(&exists).Error; err != nil {
		log.Printf("error while getting sudo cmd: %+v", err)
		return false
	}
	return exists
}
