package bundle

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/containerd/containerd/remotes"
	"github.com/deislabs/oras/pkg/content"
	"github.com/deislabs/oras/pkg/oras"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

const (
	MediaTypeKubeYaml MediaType = "application/vnd.k8s.manifest.v1+yaml"
	MediaTypeKustomizeYaml MediaType = "application/vnd.k8s.kustomization.manifest.v1+yaml"

	BlobReaderContextKey = "BlobReader"
)

type MediaType string

type Bundle map[string]*Blob

type Blob struct {
	MediaType MediaType
	Content   []byte
}

type BlobReader = func(file string) (*Blob, error)

func fileBytes(file string) ([]byte, error) {
	fileReader, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("unable to load file %s: %v", file, err)
	}
	bytes, err := ioutil.ReadAll(fileReader)
	if err != nil {
		return nil, fmt.Errorf("unable to read file %s: %v", file, err)
	}
	return bytes, nil
}

func KustomizeBaseBlobReader(file string) (*Blob, error) {
	bytes, err := fileBytes(file)
	if err != nil {
		return nil, err
	}
	mediaType := MediaTypeKubeYaml
	if strings.HasSuffix(file, "kustomization.yaml") {
		mediaType = MediaTypeKustomizeYaml
	}
	return &Blob{MediaType:mediaType, Content:bytes}, nil
}

func KubeYamlBlobReader(file string) (*Blob, error) {
	bytes, err := fileBytes(file)
	if err != nil {
		return nil, err
	}
	return &Blob{MediaType:MediaTypeKubeYaml, Content:bytes}, nil
}


var _ BlobReader = KustomizeBaseBlobReader
var _ BlobReader = KubeYamlBlobReader


func Build(ctx context.Context, dir string) (Bundle, error) {
	blobReaderValue := ctx.Value(BlobReaderContextKey)
	if blobReaderValue == nil {
		return nil, fmt.Errorf("no blob reader configured")
	}
	blobReader, ok := blobReaderValue.(BlobReader)
	if !ok {
		return nil, fmt.Errorf("invalid blob reader configured")
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	bundle := Bundle{}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		path := filepath.Join(dir, f.Name())
		blob, err := blobReader(path)
		if err != nil {
			return nil, err
		}
		bundle[f.Name()] = blob
	}
	return bundle, nil
}

func Push(ctx context.Context, resolver remotes.Resolver, ref string, b Bundle) error {
	memoryStore := content.NewMemoryStore()
	pushContents := []ocispec.Descriptor{}
	for name, blob := range b {
		pushContents = append(pushContents, memoryStore.Add(name, string(blob.MediaType), blob.Content))
	}
	fmt.Printf("Pushing to %s...\n", ref)
	desc, err := oras.Push(ctx, resolver, ref, memoryStore, pushContents)
	if err != nil {
		return err
	}
	fmt.Printf("Pushed  with digest %s\n", desc.Digest)
	return nil
}

func Pull(ctx context.Context, resolver remotes.Resolver, ref, dir string) error {
	fmt.Printf("Pulling from %s and saving...\n", ref)
	fileStore := content.NewFileStore(dir)
	defer fileStore.Close()
	allowedMediaTypes := []string{string(MediaTypeKubeYaml), string(MediaTypeKustomizeYaml)}
	desc, _, err := oras.Pull(ctx, resolver, ref, fileStore, oras.WithAllowedMediaTypes(allowedMediaTypes))
	if err != nil {
		return err
	}
	fmt.Printf("Pulled from %s with digest %s\n", ref, desc.Digest)
	return nil
}
