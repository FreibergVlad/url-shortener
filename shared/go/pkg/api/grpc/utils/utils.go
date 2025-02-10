package utils

import (
	"context"

	"google.golang.org/grpc/metadata"
)

const userIDMetadataKey = "x-user-id"

func UserIDFromIncomingContext(ctx context.Context) string {
	values := metadata.ValueFromIncomingContext(ctx, userIDMetadataKey)
	if values == nil || len(values) < 1 {
		return ""
	}
	return values[0]
}

func OutgoingContextWithUserID(ctx context.Context, userID string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, userIDMetadataKey, userID)
}

func IncomingContextWithUserID(ctx context.Context, userID string) context.Context {
	return metadata.NewIncomingContext(ctx, UserIDMetadata(userID))
}

func UserIDMetadata(userID string) metadata.MD {
	return metadata.Pairs(userIDMetadataKey, userID)
}
