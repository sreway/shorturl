package grpc

import (
	"context"
	"errors"

	"golang.org/x/exp/slog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	entity "github.com/sreway/shorturl/internal/domain/url"
	"github.com/sreway/shorturl/internal/usecases/shortener"
	pb "github.com/sreway/shorturl/proto/shorturl/v1"
)

// newProtobufURL implements create protobuf url type.
func newProtobufURL(url entity.URL) *pb.URL {
	return &pb.URL{
		Id:            url.ID().String(),
		UserID:        url.UserID().String(),
		LongURL:       url.LongURL(),
		ShortURL:      url.ShortURL(),
		CorrelationID: url.CorrelationID(),
		Deleted:       url.Deleted(),
	}
}

// CreateURL implements the RPC method for creating a shortened URL.
func (d *delivery) CreateURL(ctx context.Context, in *pb.AddURLRequest) (*pb.AddURLResponse, error) {
	response := new(pb.AddURLResponse)
	if len(in.UserID) == 0 {
		d.logger.Error("invalid user id", ErrInvalidUserID, slog.String("userID", in.UserID),
			slog.String("handler", "CreateURL"))
		return nil, d.handelErrURL(ErrInvalidUserID)
	}

	url, err := d.shortener.CreateURL(ctx, in.Url, in.UserID)
	if err != nil {
		d.logger.Error("invalid user id", err, slog.String("handler", "CreateURL"))
		return nil, d.handelErrURL(err)
	}
	response.Url = newProtobufURL(url)
	return response, nil
}

// BatchURL implements the RPC method for batch creating shortened URLs.
func (d *delivery) BatchURL(ctx context.Context, in *pb.BatchAddURLRequest) (*pb.BatchAddURLResponse, error) {
	response := new(pb.BatchAddURLResponse)

	if len(in.UserID) == 0 {
		d.logger.Error("invalid user id", ErrInvalidUserID, slog.String("userID", in.UserID),
			slog.String("handler", "BatchURL"))
		return nil, d.handelErrURL(ErrInvalidUserID)
	}

	correlationID := []string{}
	rawURL := []string{}

	for _, item := range in.Urls {
		if item.CorrelationID != "" {
			correlationID = append(correlationID, item.CorrelationID)
		}

		if item.OriginalURL != "" {
			rawURL = append(rawURL, item.OriginalURL)
		}
	}

	if len(correlationID) != len(rawURL) {
		d.logger.Error("slice correlation id length is not equal to the length of raw slicer URLs",
			ErrInvalidRequest, slog.String("handler", "BatchURL"))
		return nil, d.handelErrURL(ErrInvalidRequest)
	}

	urls, err := d.shortener.BatchURL(ctx, correlationID, rawURL, in.UserID)
	if err != nil {
		d.logger.Error("failed batch add urls", err, slog.String("handler", "BatchURL"))
		return nil, d.handelErrURL(err)
	}

	pbURLs := make([]*pb.URL, len(urls))
	for idx, url := range urls {
		pbURLs[idx] = newProtobufURL(url)
	}

	response.Url = pbURLs
	return response, nil
}

// GetURL implements the RPC method for retrieving a shortened URL.
func (d *delivery) GetURL(ctx context.Context, in *pb.GetURLRequest) (*pb.GetURLResponse, error) {
	response := new(pb.GetURLResponse)

	url, err := d.shortener.GetURL(ctx, in.UrlID)
	if err != nil {
		d.logger.Error("failed get url", err, slog.String("handler", "GetURL"))
		return nil, d.handelErrURL(err)
	}
	response.Url = newProtobufURL(url)
	return response, nil
}

// GetUserURLs implements the RPC method for retrieving all URLs belonging to a user.
func (d *delivery) GetUserURLs(ctx context.Context, in *pb.GetUserURLRequest) (*pb.GetUserURLResponse, error) {
	response := new(pb.GetUserURLResponse)

	if len(in.UserID) == 0 {
		d.logger.Error("invalid user id", ErrInvalidUserID, slog.String("userID", in.UserID),
			slog.String("handler", "GetUserURLs"))
		return nil, d.handelErrURL(ErrInvalidUserID)
	}

	urls, err := d.shortener.GetUserURLs(ctx, in.UserID)
	if err != nil {
		d.logger.Error("failed get user urls", err, slog.String("handler", "GetUserURLs"))
		return nil, d.handelErrURL(err)
	}

	pbURLs := make([]*pb.URL, len(urls))
	for idx, url := range urls {
		pbURLs[idx] = newProtobufURL(url)
	}
	response.Url = pbURLs
	return response, nil
}

// DeleteURL implements the RPC method for deleting a shortened URL.
func (d *delivery) DeleteURL(ctx context.Context, in *pb.DeleteURLRequest) (*pb.DeleteURLResponse, error) {
	response := new(pb.DeleteURLResponse)
	if len(in.UserID) == 0 {
		d.logger.Error("invalid user id", ErrInvalidUserID, slog.String("userID", in.UserID),
			slog.String("handler", "DeleteURL"))
		return nil, d.handelErrURL(ErrInvalidUserID)
	}

	err := d.shortener.DeleteURL(ctx, in.UserID, in.UrlID)
	if err != nil {
		d.logger.Error("failed delete urls", err, slog.String("handler", "DeleteURL"))
		return nil, d.handelErrURL(err)
	}

	return response, nil
}

// StorageCheck implements the RPC method for storage health check..
func (d *delivery) StorageCheck(ctx context.Context, _ *pb.StorageCheckRequest) (*pb.StorageCheckResponse, error) {
	response := new(pb.StorageCheckResponse)
	err := d.shortener.StorageCheck(ctx)
	if err != nil {
		d.logger.Error("failed check storage", err, slog.String("handler", "ping"))
		return nil, d.handelErrURL(ErrStorageCheck)
	}
	return response, nil
}

func (d *delivery) handelErrURL(err error) error {
	switch {
	case errors.Is(err, ErrInvalidUserID):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, shortener.ErrDecodeURL):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, shortener.ErrParseURL):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, shortener.ErrParseUUID):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, entity.ErrAlreadyExist):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, entity.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, entity.ErrDeleted):
		return status.Error(codes.Canceled, err.Error())
	case errors.Is(err, ErrStorageCheck):
		return status.Error(codes.Unavailable, err.Error())
	default:
		return status.Error(codes.Unknown, err.Error())
	}
}
