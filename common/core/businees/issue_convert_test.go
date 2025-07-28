package businees

import (
	"fmt"
	"testing"
)

func TestLcMemberToUserIdsWithError(t *testing.T) {
	fmt.Println(LcMemberToUserIdsWithError([]string{"U_123", "0"}, true))
	fmt.Println(LcMemberToUserIdsWithError([]string{"U_123", "0", "sdfsf"}, true))
	fmt.Println(LcMemberToUserIdsWithError([]string{"U_123", "U_0"}))

}
