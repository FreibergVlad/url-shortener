package utils

import (
	"context"

	"google.golang.org/grpc/metadata"
)

const userIdMetadataKey = "x-user-id"

func UserIDFromIncomingContext(ctx context.Context) string {
	values := metadata.ValueFromIncomingContext(ctx, userIdMetadataKey)
	if values == nil || len(values) < 1 {
		return ""
	}
	return values[0]
}

func OutgoingContextWithUserID(ctx context.Context, userId string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, userIdMetadataKey, userId)
}

func IncomingContextWithUserID(ctx context.Context, userId string) context.Context {
	return metadata.NewIncomingContext(ctx, UserIdMetadata(userId))
}

func UserIdMetadata(userId string) metadata.MD {
	return metadata.Pairs(userIdMetadataKey, userId)
}
