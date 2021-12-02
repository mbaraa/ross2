package postsgen

import (
	"archive/zip"
	"encoding/base64"
	"io"
	"io/fs"
	"os"
)

const tmpDir = "./tmp"

type PostsZipper struct {
	B64Posts []string
	files    []*os.File
}

func NewPostsZipper(b64PostsImages []string) (z *PostsZipper, err error) {
	_ = os.Mkdir(tmpDir, fs.ModeDir|fs.ModePerm)
	z = &PostsZipper{B64Posts: b64PostsImages}
	err = z.makeFiles()

	return
}

func (z *PostsZipper) MakeZipFile() (*os.File, error) {
	zipFile, err := os.CreateTemp(tmpDir, "posts*.zip")
	if err != nil {
		return nil, err
	}

	zipW := zip.NewWriter(zipFile)
	defer func() {
		_ = zipW.Flush()
		_ = zipW.Close()
	}()

	for _, file := range z.files {
		_, err = file.Seek(io.SeekStart, io.SeekStart)
		if err != nil {
			return nil, err
		}

		err = z.addToZipFile(zipW, file)
		if err != nil {
			return nil, err
		}
	}

	z.finishFiles()

	return zipFile, nil
}

func (z *PostsZipper) makeFiles() error {
	for _, post := range z.B64Posts {
		decoded, err := base64.StdEncoding.DecodeString(post)
		if err != nil {
			return err
		}

		newFile, err := os.CreateTemp(tmpDir, "post*.png")
		if err != nil {
			return err
		}

		_, err = newFile.Write(decoded)
		if err != nil {
			return err
		}

		z.files = append(z.files, newFile)
	}
	return nil
}

func (z *PostsZipper) addToZipFile(zipFile *zip.Writer, file *os.File) error {
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(fileInfo)
	if err != nil {
		return err
	}

	header.Method = zip.Deflate

	fileInZip, err := zipFile.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(fileInZip, file)
	return err
}

func (z *PostsZipper) finishFiles() {
	z.closeFiles()
	for _, file := range z.files {
		_ = os.Remove(file.Name())
	}
}

func (z *PostsZipper) closeFiles() {
	for _, file := range z.files {
		_ = file.Close()
	}
}
