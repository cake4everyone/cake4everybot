package database

// Permission is the bit mask field of permissions a user can have to controll
// the bot.
type Permission uint64

// Constant permission types
const (
	PermissionDisconnectUsers Permission = 1 << iota
	PermissionMuteUsers
)

// SudoCommand is the type which stores an action.
type SudoCommand struct {
	ID          string `gorm:"index"`
	TargetUser  string `gorm:"index"`
	Permissions Permission
}

// IsAllowed return if the given sudo command is allowed.
func (cmd SudoCommand) IsAllowed() bool {
	var exists bool
	if err := GormDB.Raw("SELECT COUNT(1) FROM sudo_commands WHERE id = ? AND target_user = ? AND (permissions & ?) = ?", cmd.ID, cmd.TargetUser, cmd.Permissions, cmd.Permissions).Scan(&exists).Error; err != nil {
		log.Printf("error while getting sudo cmd: %+v", err)
		return false
	}
	return exists
}
