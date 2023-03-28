// CODE GENERATED BY github.com/volvo-cars/lingon. DO NOT EDIT.

package local

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/volvo-cars/lingon/pkg/terra"
)

func NewFile(name string, args FileArgs) *File {
	return &File{
		Args: args,
		Name: name,
	}
}

var _ terra.Resource = (*File)(nil)

type File struct {
	Name  string
	Args  FileArgs
	state *fileState
}

func (f *File) Type() string {
	return "local_file"
}

func (f *File) LocalName() string {
	return f.Name
}

func (f *File) Configuration() interface{} {
	return f.Args
}

func (f *File) Attributes() fileAttributes {
	return fileAttributes{name: f.Name}
}

func (f *File) ImportState(av io.Reader) error {
	f.state = &fileState{}
	if err := json.NewDecoder(av).Decode(f.state); err != nil {
		return fmt.Errorf("decoding state into resource %s.%s: %w", f.Type(), f.LocalName(), err)
	}
	return nil
}

func (f *File) State() (*fileState, bool) {
	return f.state, f.state != nil
}

func (f *File) StateMust() *fileState {
	if f.state == nil {
		panic(fmt.Sprintf("state is nil for resource %s.%s", f.Type(), f.LocalName()))
	}
	return f.state
}

func (f *File) DependOn() terra.Value[terra.Reference] {
	return terra.ReferenceAttribute("local_file", f.Name)
}

type FileArgs struct {
	// Content: string, optional
	Content terra.StringValue `hcl:"content,attr"`
	// ContentBase64: string, optional
	ContentBase64 terra.StringValue `hcl:"content_base64,attr"`
	// DirectoryPermission: string, optional
	DirectoryPermission terra.StringValue `hcl:"directory_permission,attr"`
	// FilePermission: string, optional
	FilePermission terra.StringValue `hcl:"file_permission,attr"`
	// Filename: string, required
	Filename terra.StringValue `hcl:"filename,attr" validate:"required"`
	// SensitiveContent: string, optional
	SensitiveContent terra.StringValue `hcl:"sensitive_content,attr"`
	// Source: string, optional
	Source terra.StringValue `hcl:"source,attr"`
	// DependsOn contains resources that File depends on
	DependsOn terra.Dependencies `hcl:"depends_on,attr"`
}
type fileAttributes struct {
	name string
}

func (f fileAttributes) Content() terra.StringValue {
	return terra.ReferenceString(terra.ReferenceAttribute("local_file", f.name, "content"))
}

func (f fileAttributes) ContentBase64() terra.StringValue {
	return terra.ReferenceString(terra.ReferenceAttribute("local_file", f.name, "content_base64"))
}

func (f fileAttributes) ContentBase64Sha256() terra.StringValue {
	return terra.ReferenceString(terra.ReferenceAttribute("local_file", f.name, "content_base64sha256"))
}

func (f fileAttributes) ContentBase64Sha512() terra.StringValue {
	return terra.ReferenceString(terra.ReferenceAttribute("local_file", f.name, "content_base64sha512"))
}

func (f fileAttributes) ContentMd5() terra.StringValue {
	return terra.ReferenceString(terra.ReferenceAttribute("local_file", f.name, "content_md5"))
}

func (f fileAttributes) ContentSha1() terra.StringValue {
	return terra.ReferenceString(terra.ReferenceAttribute("local_file", f.name, "content_sha1"))
}

func (f fileAttributes) ContentSha256() terra.StringValue {
	return terra.ReferenceString(terra.ReferenceAttribute("local_file", f.name, "content_sha256"))
}

func (f fileAttributes) ContentSha512() terra.StringValue {
	return terra.ReferenceString(terra.ReferenceAttribute("local_file", f.name, "content_sha512"))
}

func (f fileAttributes) DirectoryPermission() terra.StringValue {
	return terra.ReferenceString(terra.ReferenceAttribute("local_file", f.name, "directory_permission"))
}

func (f fileAttributes) FilePermission() terra.StringValue {
	return terra.ReferenceString(terra.ReferenceAttribute("local_file", f.name, "file_permission"))
}

func (f fileAttributes) Filename() terra.StringValue {
	return terra.ReferenceString(terra.ReferenceAttribute("local_file", f.name, "filename"))
}

func (f fileAttributes) Id() terra.StringValue {
	return terra.ReferenceString(terra.ReferenceAttribute("local_file", f.name, "id"))
}

func (f fileAttributes) SensitiveContent() terra.StringValue {
	return terra.ReferenceString(terra.ReferenceAttribute("local_file", f.name, "sensitive_content"))
}

func (f fileAttributes) Source() terra.StringValue {
	return terra.ReferenceString(terra.ReferenceAttribute("local_file", f.name, "source"))
}

type fileState struct {
	Content             string `json:"content"`
	ContentBase64       string `json:"content_base64"`
	ContentBase64Sha256 string `json:"content_base64sha256"`
	ContentBase64Sha512 string `json:"content_base64sha512"`
	ContentMd5          string `json:"content_md5"`
	ContentSha1         string `json:"content_sha1"`
	ContentSha256       string `json:"content_sha256"`
	ContentSha512       string `json:"content_sha512"`
	DirectoryPermission string `json:"directory_permission"`
	FilePermission      string `json:"file_permission"`
	Filename            string `json:"filename"`
	Id                  string `json:"id"`
	SensitiveContent    string `json:"sensitive_content"`
	Source              string `json:"source"`
}
