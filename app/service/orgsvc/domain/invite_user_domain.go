package orgsvc

import (
	"strconv"

	"github.com/star-table/startable-server/common/core/util/rand"
	"github.com/star-table/startable-server/common/core/util/uuid"
)

func GenInviteCode(currentUserId int64, sourcePlatform string) string {
	return rand.RandomInviteCode(uuid.NewUuid() + strconv.FormatInt(currentUserId, 10) + sourcePlatform)
}
