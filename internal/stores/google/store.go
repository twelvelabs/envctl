package google

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	secretmanagerpb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/twelvelabs/envctl/internal/models"
)

var (
	// https://regex101.com/r/SXjWpH/2
	secretPathRe = regexp.MustCompile(`(?P<version>(?P<secret>(?P<parent>projects\/[\w-]+(?P<parent_suffix>\/locations\/[\w-]+)?)\/secrets\/[\w-]+)(?P<version_suffix>\/versions\/[\w-]+)?)`) //nolint:lll

	ErrNotFound   = errors.New("secret not found")
	ErrInvalidURL = errors.New("invalid URL")
)

// GSMStore provides methods for accessing secrets stored
// in Google Secret Manager.
type GSMStore struct {
	client GSMClient
	cache  map[string]string
}

// NewGSMStore returns a new [GSMStore].
func NewGSMStore(client GSMClient) *GSMStore {
	return &GSMStore{
		client: client,
		cache:  map[string]string{},
	}
}

func (s *GSMStore) Close() error {
	return s.client.Close()
}

// Get returns the plaintext value for the secret at the given URL.
func (s *GSMStore) Get(ctx context.Context, value models.Value) (string, error) {
	names, err := s.names(value)
	if err != nil {
		return "", err
	}

	_, hasDefault := value.URL().Default()
	plaintext, err := s.accessSecretVersion(ctx, names.Version, hasDefault)
	if err != nil {
		return "", err
	}

	return models.ReadPlaintext(plaintext, value.URL())
}

// Set updates the plaintext value for the secret at the given URL.
func (s *GSMStore) Set(ctx context.Context, value models.Value, updated string) error {
	names, err := s.names(value)
	if err != nil {
		return err
	}

	// Ensure the secret exists.
	_, err = s.getOrCreateSecret(ctx, names.Secret, names.Parent)
	if err != nil {
		return err
	}

	// Get the current secret version.
	plaintext, err := s.accessSecretVersion(ctx, names.Version, true)
	if err != nil {
		return err
	}

	// Update the plaintext value.
	plaintext, err = models.WritePlaintext(plaintext, value.URL(), updated)
	if err != nil {
		return err
	}

	// Publish a new secret version w/ the updated plaintext.
	err = s.updateSecretVersion(ctx, names.Secret, names.Version, plaintext)
	if err != nil {
		return err
	}

	return nil
}

// Delete destroys the secret at the given URL.
func (s *GSMStore) Delete(ctx context.Context, value models.Value) error {
	return nil
}

type parsedNames struct {
	Parent  string
	Secret  string
	Version string
}

// Parses the given input into:
//
//   - ParentName: `projects/*`
//   - SecretName: `projects/*/secrets/*`
//   - SecretVersionName: `projects/*/secrets/*/versions/*`
//
// Also accepts long-form URLs that contain `locations/*`.
func (s *GSMStore) names(value models.Value) (*parsedNames, error) {
	url := value.URL()
	path := url.Path
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

	match := secretPathRe.FindStringSubmatch(path)
	if match == nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidURL, url.String())
	}

	results := map[string]string{}
	for i, name := range secretPathRe.SubexpNames() {
		results[name] = match[i]
	}

	if results["version_suffix"] == "" {
		// Ensure a valid version name.
		results["version"] += "/versions/latest"
	}

	return &parsedNames{
		Parent:  results["parent"],
		Secret:  results["secret"],
		Version: results["version"],
	}, nil
}

// Ensures a secret exists for the given name.
func (s *GSMStore) getOrCreateSecret(
	ctx context.Context, name string, parent string,
) (*secretmanagerpb.Secret, error) {
	secret, err := s.client.GetSecret(ctx, &secretmanagerpb.GetSecretRequest{
		Name: name,
	})

	// Handle happy path and unknown error.
	// Letting ErrNotFound fall through to the next section.
	err = rpcErr(err, name)
	if err == nil {
		return secret, nil
	}
	if !errors.Is(err, ErrNotFound) {
		return nil, err
	}

	// SecretID is the last segment of the secret name.
	// i.e. `foo` for `projects/my-proj/secrets/foo`.
	segments := strings.Split(name, "/")
	secretID := segments[len(segments)-1]

	secret, err = s.client.CreateSecret(ctx, &secretmanagerpb.CreateSecretRequest{
		Parent:   parent,
		SecretId: secretID,
		Secret: &secretmanagerpb.Secret{
			Replication: &secretmanagerpb.Replication{
				Replication: &secretmanagerpb.Replication_Automatic_{
					Automatic: &secretmanagerpb.Replication_Automatic{},
				},
			},
		},
	})

	return secret, rpcErr(err, parent)
}

func (s *GSMStore) accessSecretVersion(ctx context.Context, versionName string, hasDefault bool) (string, error) {
	plaintext, found := s.cache[versionName]
	if !found {
		resp, err := s.client.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
			Name: versionName,
		})
		err = rpcErr(err, versionName)
		if err != nil {
			if errors.Is(err, ErrNotFound) && hasDefault {
				// Not really an error if we have a default value.
				err = nil
			}
			return "", err
		}
		plaintext = string(resp.GetPayload().GetData())
		s.cache[versionName] = plaintext
	}
	return plaintext, nil
}

func (s *GSMStore) updateSecretVersion(ctx context.Context, secretName, versionName, plaintext string) error {
	_, err := s.client.AddSecretVersion(ctx, &secretmanagerpb.AddSecretVersionRequest{
		Parent: secretName,
		Payload: &secretmanagerpb.SecretPayload{
			Data: []byte(plaintext),
		},
	})
	err = rpcErr(err, secretName)
	if err != nil {
		return err
	}

	delete(s.cache, versionName)
	return nil
}

// Helper to pretty up the default gRPC 404 errors
// so we can return them directly to users.
func rpcErr(err error, name string) error {
	status, _ := status.FromError(err)
	switch status.Code() {
	case codes.OK:
		return nil
	case codes.NotFound:
		return fmt.Errorf("%w: %s", ErrNotFound, name)
	default:
		return status.Err()
	}
}
