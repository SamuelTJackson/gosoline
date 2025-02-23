package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsCfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/justtrackio/gosoline/pkg/cfg"
)

func GetCredentialsProvider(ctx context.Context, config cfg.Config, settings ClientSettings) (aws.CredentialsProvider, error) {
	if len(settings.AssumeRole) > 0 {
		return GetAssumeRoleCredentialsProvider(ctx, settings.AssumeRole)
	}

	if creds := UnmarshalCredentials(config); creds != nil {
		return credentials.NewStaticCredentialsProvider(creds.AccessKeyID, creds.SecretAccessKey, creds.SessionToken), nil
	}

	return nil, nil
}

func GetAssumeRoleCredentialsProvider(ctx context.Context, roleArn string) (aws.CredentialsProvider, error) {
	var err error
	var cfg aws.Config

	if cfg, err = awsCfg.LoadDefaultConfig(ctx); err != nil {
		return nil, fmt.Errorf("can not load default aws config: %w", err)
	}

	stsClient := sts.NewFromConfig(cfg)

	return stscreds.NewAssumeRoleProvider(stsClient, roleArn), nil
}

func UnmarshalCredentials(config cfg.Config) *Credentials {
	if !config.IsSet("cloud.aws.credentials") {
		return nil
	}

	creds := &Credentials{}
	config.UnmarshalKey("cloud.aws.credentials", creds)

	return creds
}
