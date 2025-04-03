package internal

import (
	"context"
	"database/sql"

	"untether/services/plaid/client"
	pb "untether/services/plaid/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PlaidService struct {
	pb.UnimplementedPlaidServiceServer
	client client.PlaidClient
	db     *sql.DB
}

func NewPlaidService(client client.PlaidClient, db *sql.DB) *PlaidService {
	return &PlaidService{
		client: client,
		db:     db,
	}
}

func (s *PlaidService) CreateLinkToken(ctx context.Context, req *pb.CreateLinkTokenRequest) (*pb.CreateLinkTokenResponse, error) {
	if req.UserId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user_id is required")
	}

	token, err := s.client.CreateLinkToken(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create link token: %v", err)
	}

	return &pb.CreateLinkTokenResponse{
		LinkToken: token,
	}, nil
}

func (s *PlaidService) ExchangePublicToken(ctx context.Context, req *pb.ExchangePublicTokenRequest) (*pb.ExchangePublicTokenResponse, error) {
	if req.PublicToken == "" {
		return nil, status.Errorf(codes.InvalidArgument, "public_token is required")
	}

	accessToken, err := s.client.ExchangePublicToken(ctx, req.PublicToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to exchange public token: %v", err)
	}

	return &pb.ExchangePublicTokenResponse{
		AccessToken: accessToken,
	}, nil
}

func (s *PlaidService) GetAccounts(ctx context.Context, req *pb.GetAccountsRequest) (*pb.GetAccountsResponse, error) {
	if req.AccessToken == "" {
		return nil, status.Errorf(codes.InvalidArgument, "access_token is required")
	}

	accounts, err := s.client.GetAccounts(ctx, req.AccessToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get accounts: %v", err)
	}

	// Note: We can't store accounts in the database here because we don't have the user_id
	// This would need to be handled by the client service that has the user context

	pbAccounts := make([]*pb.BankAccount, len(accounts))
	for i, account := range accounts {
		pbAccounts[i] = &pb.BankAccount{
			AccountId: account.AccountID,
			Name:      account.Name,
			Type:      account.Type,
			Subtype:   account.Subtype,
			Mask:      account.Mask,
		}
	}

	return &pb.GetAccountsResponse{
		Accounts: pbAccounts,
	}, nil
}

func (s *PlaidService) GetBalance(ctx context.Context, req *pb.GetBalanceRequest) (*pb.GetBalanceResponse, error) {
	if req.AccessToken == "" {
		return nil, status.Errorf(codes.InvalidArgument, "access_token is required")
	}
	if req.AccountId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "account_id is required")
	}

	balance, err := s.client.GetBalance(ctx, req.AccessToken, req.AccountId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get balance: %v", err)
	}

	return &pb.GetBalanceResponse{
		Balance: balance,
	}, nil
}
