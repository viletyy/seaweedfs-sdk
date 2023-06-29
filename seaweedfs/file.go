package seaweedfs

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
)

type DirAssign struct {
	Count     int    `json:"count"`
	Fid       string `json:"fid"`
	Url       string `json:"url"`
	PublicUrl string `json:"publicUrl"`
}

// > curl http://localhost:9333/dir/assign
// {"count":1,"fid":"3,01637037d6","url":"127.0.0.1:8080","publicUrl":"localhost:8080"}
func (c *Client) GetDirAssign() (*DirAssign, error) {
	da := new(DirAssign)
	err := c.getMasterParsedResponse("GET", "/dir/assign", nil, nil, da)
	return da, err
}

type File struct {
	Name string `json:"name"`
	Size string `json:"size"`
	ETag string `json:"eTag"`
}

// > curl -F file=@/home/chris/myphoto.jpg http://127.0.0.1:8080/3,01637037d6
// {"name":"myphoto.jpg","size":43234,"eTag":"1cc0118e"}
func (c *Client) PostFile(fid string, file io.Reader, filename string) (*File, error) {
	f := new(File)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}

	if _, err = io.Copy(part, file); err != nil {
		return nil, err
	}
	if err = writer.Close(); err != nil {
		return nil, err
	}

	err = c.getVolumnParsedResponse("POST",
		"/"+fid,
		http.Header{"Content-Type": {writer.FormDataContentType()}},
		body,
		&f)
	return f, err
}

// > curl -X DELETE http://127.0.0.1:8080/3,01637037d6
func (c *Client) DeleteFile(fid string) error {
	_, err := c.getVolumnResponse("DELETE", "/"+fid, nil, nil)

	return err
}
