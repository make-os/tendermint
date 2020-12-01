package grpc

import (
	context "context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/tendermint/tendermint/crypto"
	cryptoenc "github.com/tendermint/tendermint/crypto/encoding"
	"github.com/tendermint/tendermint/libs/log"
	privvalproto "github.com/tendermint/tendermint/proto/tendermint/privval"
	"github.com/tendermint/tendermint/types"
)

type SignerServer struct {
	Logger log.Logger

	ChainID string
	PrivVal types.PrivValidator
}

func NewSignerServer(chainID string,
	privVal types.PrivValidator, log log.Logger) *SignerServer {

	return &SignerServer{
		Logger:  log,
		ChainID: chainID,
		PrivVal: privVal,
	}
}

var _ privvalproto.PrivValidatorAPIServer = (*SignerServer)(nil)

// PubKey receives a request for the pubkey
// returns the pubkey on success and error on failure
func (ss *SignerServer) GetPubKey(ctx context.Context, req *privvalproto.PubKeyRequest) (
	*privvalproto.PubKeyResponse, error) {
	var pubKey crypto.PubKey

	pubKey, err := ss.PrivVal.GetPubKey()
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "error getting pubkey: %v", err)
	}

	pk, err := cryptoenc.PubKeyToProto(pubKey)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "error transistioning pubkey to proto: %v", err)
	}

	return &privvalproto.PubKeyResponse{PubKey: pk}, nil
}

// SignVote receives a vote sign requests, attempts to sign it
// returns SignedVoteResponse on success and error on failure
func (ss *SignerServer) SignVote(ctx context.Context, req *privvalproto.SignVoteRequest) (
	*privvalproto.SignedVoteResponse, error) {
	vote := req.Vote

	err := ss.PrivVal.SignVote(req.ChainId, vote)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "error signing vote: %v", err)
	}

	return &privvalproto.SignedVoteResponse{Vote: *vote}, nil
}

// SignProposal receives a proposal sign requests, attempts to sign it
// returns SignedProposalResponse on success and error on failure
func (ss *SignerServer) SignProposal(ctx context.Context, req *privvalproto.SignProposalRequest) (
	*privvalproto.SignedProposalResponse, error) {
	proposal := req.Proposal

	err := ss.PrivVal.SignProposal(req.ChainId, proposal)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "error signing proposal: %v", err)
	}

	return &privvalproto.SignedProposalResponse{Proposal: *proposal}, nil
}
