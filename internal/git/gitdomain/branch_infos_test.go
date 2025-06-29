package gitdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestBranchInfos(t *testing.T) {
	t.Parallel()

	t.Run("BranchIsActiveInAnotherWorktree", func(t *testing.T) {
		t.Parallel()
		t.Run("branch is active in another worktree", func(t *testing.T) {
			t.Parallel()
			branchInfos := gitdomain.BranchInfos{
				{
					LocalName:  gitdomain.NewLocalBranchNameOption("branch-1"),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
				},
				{
					LocalName:  gitdomain.NewLocalBranchNameOption("branch-2"),
					SyncStatus: gitdomain.SyncStatusOtherWorktree,
				},
			}
			have := branchInfos.BranchIsActiveInAnotherWorktree("branch-2")
			must.True(t, have)
		})
		t.Run("branch is local but not active in another worktree", func(t *testing.T) {
			t.Parallel()
			branchInfos := gitdomain.BranchInfos{
				{
					LocalName:  gitdomain.NewLocalBranchNameOption("branch-1"),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
				},
			}
			have := branchInfos.BranchIsActiveInAnotherWorktree("branch-1")
			must.False(t, have)
		})
		t.Run("branch is remote", func(t *testing.T) {
			t.Parallel()
			branchInfos := gitdomain.BranchInfos{
				{
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
				},
			}
			have := branchInfos.BranchIsActiveInAnotherWorktree("branch-1")
			must.False(t, have)
		})
		t.Run("branch doesn't exist", func(t *testing.T) {
			t.Parallel()
			branchInfos := gitdomain.BranchInfos{
				{
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
				},
			}
			have := branchInfos.BranchIsActiveInAnotherWorktree("zonk")
			must.False(t, have)
		})
	})

	t.Run("FindByRemote", func(t *testing.T) {
		t.Parallel()
		t.Run("has a local branch with matching tracking branch", func(t *testing.T) {
			t.Parallel()
			branch := gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchNameOption("two"),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusUpToDate,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/two")),
				RemoteSHA:  Some(gitdomain.NewSHA("111111")),
			}
			bs := gitdomain.BranchInfos{branch}
			have, has := bs.FindByRemoteName("origin/two").Get()
			must.True(t, has)
			must.Eq(t, &branch, have)
		})
		t.Run("has a local branch with the given name", func(t *testing.T) {
			t.Parallel()
			bs := gitdomain.BranchInfos{gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchNameOption("kg/one"),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			}}
			have := bs.FindByRemoteName("kg/one")
			must.True(t, have.IsNone())
		})
	})

	t.Run("FindLocalOrRemote", func(t *testing.T) {
		t.Parallel()
		t.Run("has local name", func(t *testing.T) {
			t.Parallel()
			branch1 := gitdomain.NewLocalBranchName("branch-1")
			branch1info := gitdomain.BranchInfo{
				LocalName:  Some(branch1),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			}
			bis := gitdomain.BranchInfos{
				branch1info,
			}
			have := bis.FindLocalOrRemote(branch1, gitdomain.RemoteOrigin)
			must.Eq(t, MutableSome(&branch1info), have)
		})
		t.Run("has remote name", func(t *testing.T) {
			t.Parallel()
			branch1info := gitdomain.BranchInfo{
				LocalName:  None[gitdomain.LocalBranchName](),
				LocalSHA:   None[gitdomain.SHA](),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
				RemoteSHA:  Some(gitdomain.NewSHA("111111")),
			}
			bis := gitdomain.BranchInfos{
				branch1info,
			}
			have := bis.FindLocalOrRemote(gitdomain.NewLocalBranchName("branch-1"), gitdomain.RemoteOrigin)
			must.Eq(t, MutableSome(&branch1info), have)
		})
		t.Run("no match", func(t *testing.T) {
			t.Parallel()
			branch1info := gitdomain.BranchInfo{
				LocalName:  None[gitdomain.LocalBranchName](),
				LocalSHA:   None[gitdomain.SHA](),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
				RemoteSHA:  Some(gitdomain.NewSHA("111111")),
			}
			bis := gitdomain.BranchInfos{
				branch1info,
			}
			have := bis.FindLocalOrRemote(gitdomain.NewLocalBranchName("zonk"), gitdomain.RemoteOrigin)
			must.Eq(t, MutableNone[gitdomain.BranchInfo](), have)
		})
	})

	t.Run("FindMatchingRecord", func(t *testing.T) {
		t.Parallel()
		t.Run("has matching local name", func(t *testing.T) {
			t.Parallel()
			bis := gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchNameOption("branch-1"),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
			}
			give := gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchNameOption("branch-1"),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			}
			have := bis.FindMatchingRecord(give)
			want := MutableSome(&bis[0])
			must.Eq(t, want, have)
		})
		t.Run("has matching remote name", func(t *testing.T) {
			t.Parallel()
			bis := gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
			}
			give := gitdomain.BranchInfo{
				LocalName:  None[gitdomain.LocalBranchName](),
				LocalSHA:   None[gitdomain.SHA](),
				SyncStatus: gitdomain.SyncStatusRemoteOnly,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/branch-1")),
				RemoteSHA:  Some(gitdomain.NewSHA("111111")),
			}
			have := bis.FindMatchingRecord(give)
			want := MutableSome(&bis[0])
			must.Eq(t, want, have)
		})
	})

	t.Run("HasLocalBranch", func(t *testing.T) {
		t.Parallel()
		t.Run("has a matching local branch", func(t *testing.T) {
			t.Parallel()
			bs := gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchNameOption("one"),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
			}
			must.True(t, bs.HasLocalBranch(gitdomain.NewLocalBranchName("one")))
		})
		t.Run("has a remote branch with that name", func(t *testing.T) {
			t.Parallel()
			bs := gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/one")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
			}
			must.False(t, bs.HasLocalBranch(gitdomain.NewLocalBranchName("one")))
		})
		t.Run("has a local branch with a matching tracking branch", func(t *testing.T) {
			t.Parallel()
			bs := gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchNameOption("two"),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/one")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
			}
			must.False(t, bs.HasLocalBranch(gitdomain.NewLocalBranchName("one")))
		})
	})

	t.Run("HasMatchingRemoteBranchFor", func(t *testing.T) {
		t.Parallel()
		t.Run("has a local branch with a matching remote", func(t *testing.T) {
			t.Parallel()
			bs := gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchNameOption("two"),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusUpToDate,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/one")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
			}
			must.True(t, bs.HasMatchingTrackingBranchFor(gitdomain.NewLocalBranchName("one"), gitdomain.RemoteOrigin))
		})
		t.Run("has a remote-only branch with that name", func(t *testing.T) {
			t.Parallel()
			bs := gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(gitdomain.NewRemoteBranchName("origin/one")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
			}
			must.True(t, bs.HasMatchingTrackingBranchFor(gitdomain.NewLocalBranchName("one"), gitdomain.RemoteOrigin))
		})
		t.Run("has a local branch with a matching name", func(t *testing.T) {
			t.Parallel()
			bs := gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchNameOption("one"),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
			}
			must.False(t, bs.HasMatchingTrackingBranchFor(gitdomain.NewLocalBranchName("one"), gitdomain.RemoteOrigin))
		})
	})

	t.Run("LocalBranches", func(t *testing.T) {
		t.Parallel()
		bs := gitdomain.BranchInfos{
			gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchNameOption("up-to-date"),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusUpToDate,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/up-to-date")),
				RemoteSHA:  Some(gitdomain.NewSHA("111111")),
			},
			gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchNameOption("ahead"),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusNotInSync,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/ahead")),
				RemoteSHA:  Some(gitdomain.NewSHA("222222")),
			},
			gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchNameOption("behind"),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusNotInSync,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/behind")),
				RemoteSHA:  Some(gitdomain.NewSHA("222222")),
			},
			gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchNameOption("local-only"),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			},
			gitdomain.BranchInfo{
				LocalName:  None[gitdomain.LocalBranchName](),
				LocalSHA:   None[gitdomain.SHA](),
				SyncStatus: gitdomain.SyncStatusRemoteOnly,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/remote-only")),
				RemoteSHA:  Some(gitdomain.NewSHA("111111")),
			},
			gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchNameOption("deleted-at-remote"),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusDeletedAtRemote,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			},
		}
		have := bs.LocalBranches().Names()
		want := gitdomain.NewLocalBranchNames("up-to-date", "ahead", "behind", "local-only", "deleted-at-remote")
		must.Eq(t, want, have)
	})

	t.Run("LocalBranchesWithDeletedTrackingBranches", func(t *testing.T) {
		t.Parallel()
		bs := gitdomain.BranchInfos{
			gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchNameOption("up-to-date"),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusUpToDate,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/up-to-date")),
				RemoteSHA:  Some(gitdomain.NewSHA("111111")),
			},
			gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchNameOption("ahead"),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusNotInSync,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/ahead")),
				RemoteSHA:  Some(gitdomain.NewSHA("222222")),
			},
			gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchNameOption("behind"),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusNotInSync,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/behind")),
				RemoteSHA:  Some(gitdomain.NewSHA("222222")),
			},
			gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchNameOption("local-only"),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			},
			gitdomain.BranchInfo{
				LocalName:  None[gitdomain.LocalBranchName](),
				LocalSHA:   None[gitdomain.SHA](),
				SyncStatus: gitdomain.SyncStatusRemoteOnly,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/remote-only")),
				RemoteSHA:  Some(gitdomain.NewSHA("111111")),
			},
			gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchNameOption("deleted-at-remote"),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusDeletedAtRemote,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			},
		}
		have := bs.LocalBranchesWithDeletedTrackingBranches().Names()
		want := gitdomain.NewLocalBranchNames("deleted-at-remote")
		must.Eq(t, want, have)
	})

	t.Run("LookupLocalBranch", func(t *testing.T) {
		t.Parallel()
		t.Run("local branch with matching name", func(t *testing.T) {
			t.Parallel()
			branchOne := gitdomain.NewLocalBranchName("one")
			branchInfos := gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  Some(branchOne),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
			}
			branchInfo, hasBranchInfo := branchInfos.FindByLocalName(branchOne).Get()
			must.True(t, hasBranchInfo)
			must.EqOp(t, branchOne, branchInfo.LocalName.GetOrPanic())
		})
		t.Run("remote branch with matching name", func(t *testing.T) {
			t.Parallel()
			branchInfos := gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  None[gitdomain.LocalBranchName](),
					LocalSHA:   None[gitdomain.SHA](),
					SyncStatus: gitdomain.SyncStatusRemoteOnly,
					RemoteName: Some(gitdomain.NewRemoteBranchName("kg/one")),
					RemoteSHA:  Some(gitdomain.NewSHA("111111")),
				},
			}
			have := branchInfos.FindByLocalName(gitdomain.NewLocalBranchName("kg/one"))
			must.True(t, have.IsNone())
		})
	})

	t.Run("Names", func(t *testing.T) {
		t.Parallel()
		bs := gitdomain.BranchInfos{
			gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchNameOption("one"),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			},
			gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchNameOption("two"),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			},
			gitdomain.BranchInfo{
				LocalName:  None[gitdomain.LocalBranchName](),
				LocalSHA:   None[gitdomain.SHA](),
				SyncStatus: gitdomain.SyncStatusRemoteOnly,
				RemoteName: Some(gitdomain.NewRemoteBranchName("origin/three")),
				RemoteSHA:  Some(gitdomain.NewSHA("111111")),
			},
		}
		have := bs.Names()
		want := gitdomain.NewLocalBranchNames("one", "two")
		must.Eq(t, want, have)
	})

	t.Run("Remove", func(t *testing.T) {
		t.Parallel()
		t.Run("contains the removed element", func(t *testing.T) {
			t.Parallel()
			bs := gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchNameOption("one"),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
				gitdomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchNameOption("two"),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
			}
			have := bs.Remove(gitdomain.NewLocalBranchName("two"))
			want := gitdomain.BranchInfos{
				gitdomain.BranchInfo{
					LocalName:  gitdomain.NewLocalBranchNameOption("one"),
					LocalSHA:   Some(gitdomain.NewSHA("111111")),
					SyncStatus: gitdomain.SyncStatusLocalOnly,
					RemoteName: None[gitdomain.RemoteBranchName](),
					RemoteSHA:  None[gitdomain.SHA](),
				},
			}
			must.Eq(t, want, have)
		})
	})

	t.Run("Select", func(t *testing.T) {
		t.Parallel()
		bs := gitdomain.BranchInfos{
			gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchNameOption("one"),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			},
			gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchNameOption("two"),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			},
			gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchNameOption("three"),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			},
			gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchNameOption("four"),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			},
		}
		have, nonExisting := bs.Select(gitdomain.RemoteOrigin, "one", "three")
		want := gitdomain.BranchInfos{
			gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchNameOption("one"),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			},
			gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchNameOption("three"),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			},
		}
		must.Eq(t, want, have)
		must.Eq(t, nonExisting, gitdomain.LocalBranchNames(nil))
	})

	t.Run("does not contain the removed element", func(t *testing.T) {
		t.Parallel()
		bs := gitdomain.BranchInfos{
			gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchNameOption("one"),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			},
			gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchNameOption("two"),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			},
		}
		have := bs.Remove(gitdomain.NewLocalBranchName("zonk"))
		want := gitdomain.BranchInfos{
			gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchNameOption("one"),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			},
			gitdomain.BranchInfo{
				LocalName:  gitdomain.NewLocalBranchNameOption("two"),
				LocalSHA:   Some(gitdomain.NewSHA("111111")),
				SyncStatus: gitdomain.SyncStatusLocalOnly,
				RemoteName: None[gitdomain.RemoteBranchName](),
				RemoteSHA:  None[gitdomain.SHA](),
			},
		}
		must.Eq(t, want, have)
	})
}
