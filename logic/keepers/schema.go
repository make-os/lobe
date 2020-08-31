package keepers

import (
	"github.com/make-os/lobe/storage/common"
	"github.com/make-os/lobe/util"
)

const (
	TagAccount               = "a"
	TagPushKey               = "g"
	TagAddressPushKeyID      = "ag"
	TagRepo                  = "r"
	TagRepoPropVote          = "rpv"
	TagRepoPropEndIndex      = "rei"
	TagNS                    = "ns"
	TagClosedProp            = "cp"
	TagBlockInfo             = "b"
	TagLastRepoSyncherHeight = "rh"
	TagHelmRepo              = "hr"
	TagValidators            = "v"
	TagTx                    = "t"
	TagTrackedRepo           = "tr"
)

// MakeTrackedRepoKey creates a key for accessing a tracked repo.
func MakeTrackedRepoKey(name string) []byte {
	return common.MakePrefix([]byte(TagTrackedRepo), []byte(name))
}

// MakeQueryTrackedRepoKey creates a key for accessing all tracked repo.
func MakeQueryTrackedRepoKey() []byte {
	return common.MakePrefix([]byte(TagTrackedRepo))
}

// MakeAccountKey creates a key for accessing an account
func MakeAccountKey(address string) []byte {
	return common.MakePrefix([]byte(TagAccount), []byte(address))
}

// MakePushKeyKey creates a key for storing push key
func MakePushKeyKey(pushKeyID string) []byte {
	return common.MakePrefix([]byte(TagPushKey), []byte(pushKeyID))
}

// MakeAddrPushKeyIDIndexKey creates a key for address to push key index
func MakeAddrPushKeyIDIndexKey(address, pushKeyID string) []byte {
	return common.MakePrefix([]byte(TagAddressPushKeyID), []byte(address), []byte(pushKeyID))
}

// MakeQueryPushKeyIDsOfAddress creates a key for querying push key ids belonging to an address
func MakeQueryPushKeyIDsOfAddress(address string) []byte {
	return common.MakePrefix([]byte(TagAddressPushKeyID), []byte(address))
}

// MakeRepoKey creates a key for accessing a repository object
func MakeRepoKey(name string) []byte {
	return common.MakePrefix([]byte(TagRepo), []byte(name))
}

// MakeRepoProposalVoteKey creates a key as flag for a repo proposal vote
func MakeRepoProposalVoteKey(repoName, proposalID, voterAddr string) []byte {
	return common.MakePrefix([]byte(TagRepoPropVote), []byte(repoName),
		[]byte(proposalID), []byte(voterAddr))
}

// MakeRepoProposalEndIndexKey creates a key that makes a repo proposal to its
// end height
func MakeRepoProposalEndIndexKey(repoName, proposalID string, endHeight uint64) []byte {
	return common.MakePrefix([]byte(TagRepoPropEndIndex), util.EncodeNumber(endHeight),
		[]byte(repoName), []byte(proposalID))
}

// MakeQueryKeyRepoProposalAtEndHeight creates a key for finding repo proposals
// ending at the given height
func MakeQueryKeyRepoProposalAtEndHeight(endHeight uint64) []byte {
	return common.MakePrefix([]byte(TagRepoPropEndIndex), util.EncodeNumber(endHeight))
}

// MakeClosedProposalKey creates a key for marking a proposal as "closed"
func MakeClosedProposalKey(name, propID string) []byte {
	return common.MakePrefix([]byte(TagClosedProp), []byte(name), []byte(propID))
}

// MakeNamespaceKey creates a key for accessing a namespace
func MakeNamespaceKey(name string) []byte {
	return common.MakePrefix([]byte(TagNS), []byte(name))
}

// MakeKeyBlockInfo creates a key for accessing/storing committed block data.
func MakeKeyBlockInfo(height int64) []byte {
	return common.MakeKey(util.EncodeNumber(uint64(height)), []byte(TagBlockInfo))
}

// MakeKeyRepoSyncerHeight creates a key for accessing last height synch-ed by
// the repo syncer
func MakeKeyRepoSyncerHeight() []byte {
	return common.MakePrefix([]byte(TagLastRepoSyncherHeight))
}

// MakeKeyHelmRepo creates a key for getting/setting the helm repo
func MakeKeyHelmRepo() []byte {
	return common.MakePrefix([]byte(TagHelmRepo))
}

// MakeQueryKeyBlockInfo creates a key for querying committed block data
func MakeQueryKeyBlockInfo() []byte {
	return common.MakePrefix([]byte(TagBlockInfo))
}

// MakeNetMaturityKey creates a key indicating the network's maturity status

// MakeBlockValidatorsKey creates a key for storing validators of blocks
func MakeBlockValidatorsKey(height int64) []byte {
	return common.MakeKey(util.EncodeNumber(uint64(height)), []byte(TagValidators))
}

// MakeQueryKeyBlockValidators creates a key for querying all block validators
func MakeQueryKeyBlockValidators() []byte {
	return common.MakePrefix([]byte(TagValidators))
}

// MakeTxKey creates a key for storing validators of blocks
func MakeTxKey(hash []byte) []byte {
	return common.MakePrefix([]byte(TagTx), hash)
}
