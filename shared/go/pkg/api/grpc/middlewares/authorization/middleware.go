package authorization

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	permissionsProtoMessages "github.com/FreibergVlad/url-shortener/proto/pkg/permissions/messages/v1"
	grpcUtils "github.com/FreibergVlad/url-shortener/shared/go/pkg/api/grpc/utils"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type permissionResolver interface {
	HasPermissions(ctx context.Context, scopes []string, userId string, organizationId *string) (bool, error)
}

func New(
	servicesDesc map[string]protoreflect.ServiceDescriptor,
	permissionResolver permissionResolver,
) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		methodDesc, err := getMethodDescriptor(servicesDesc, info.FullMethod)
		if err != nil {
			return nil, err
		}
		isAuthRequired := isAuthenticationRequired(methodDesc)
		if !isAuthRequired {
			return handler(ctx, req)
		}

		userId := grpcUtils.UserIDFromIncomingContext(ctx)
		if userId == "" {
			return nil, errors.NewUnauthenticatedError("unauthenticated")
		}

		scopes := getRequiredScopes(methodDesc)
		if len(scopes) == 0 {
			return handler(ctx, req)
		}

		organizationId := organizationIdFromRequest(req)
		ok, err := permissionResolver.HasPermissions(ctx, scopes, userId, organizationId)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve permissions: %w", err)
		}
		if !ok {
			if organizationId == nil {
				return nil, errors.NewPermissionDeniedError("user %s is not allowed to perform %s", userId, scopes)
			}
			return nil, errors.NewPermissionDeniedError("user %s is not allowed to perform %s within organization %s", userId, scopes, *organizationId)
		}
		return handler(ctx, req)
	}
}

func isAuthenticationRequired(methodDesc protoreflect.MethodDescriptor) bool {
	if proto.HasExtension(methodDesc.Options(), permissionsProtoMessages.E_AuthenticationRequired) {
		ext := proto.GetExtension(methodDesc.Options(), permissionsProtoMessages.E_AuthenticationRequired)
		if extBool, ok := ext.(*wrapperspb.BoolValue); ok && extBool != nil {
			return extBool.Value
		}
	}
	return true
}

func organizationIdFromRequest(req any) *string {
	v := reflect.ValueOf(req).Elem()
	field := v.FieldByName("OrganizationId")
	if !field.IsValid() {
		return nil
	}
	orgId := field.String()
	return &orgId
}

func getRequiredScopes(methodDesc protoreflect.MethodDescriptor) []string {
	if proto.HasExtension(methodDesc.Options(), permissionsProtoMessages.E_RequiredPermissions) {
		scopes, ok := proto.GetExtension(methodDesc.Options(), permissionsProtoMessages.E_RequiredPermissions).([]string)
		if ok {
			return scopes
		}
	}
	return []string{}
}

func getMethodDescriptor(servicesDesc map[string]protoreflect.ServiceDescriptor, fullMethodName string) (protoreflect.MethodDescriptor, error) {
	fullMethodNameParts := strings.Split(fullMethodName, "/")
	if len(fullMethodNameParts) < 3 {
		return nil, errors.NewInternalError("authorization: method %s not found", fullMethodName)
	}
	serviceNameParts := strings.Split(fullMethodNameParts[1], ".")
	shortServiceName, methodName := serviceNameParts[len(serviceNameParts)-1], fullMethodNameParts[2]
	serviceDesc, ok := servicesDesc[shortServiceName]
	if !ok {
		return nil, errors.NewInternalError("authorization: method %s not found", fullMethodName)
	}
	md := serviceDesc.Methods().ByName(protoreflect.Name(methodName))
	if md == nil {
		return nil, errors.NewInternalError("authorization: method %s not found", fullMethodName)
	}
	return md, nil
}
