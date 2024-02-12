package shared_test

import (
	"testing"

	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/vm/opcode"
	"github.com/git-town/git-town/v12/src/vm/shared"
	"github.com/shoenig/test/must"
)

func TestIsCheckout(t *testing.T) {
	t.Parallel()
	branch := gitdomain.NewLocalBranchName("foo")
	tests := map[shared.Opcode]bool{
		&opcode.Checkout{Branch: branch}:         true,  // Checkout is (obviously) a checkout opcode
		&opcode.CheckoutIfExists{Branch: branch}: true,  // CheckoutIfExists is also a checkout opcode
		&opcode.AbortMerge{}:                     false, // any other opcode doesn't match
	}
	for give, want := range tests {
		have := shared.IsCheckoutOpcode(give)
		must.Eq(t, want, have)
	}
}
