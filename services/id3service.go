package services

import (
	"github.com/bogem/id3v2"
)

type ID3Service struct{}
type ID3Tags struct {
	Title, Artist, Album string
}

func (i *ID3Service) Read(file string) (ID3Tags, error) {
	tag, err := id3v2.Open(file, id3v2.Options{
		Parse: true,
	})
	if err != nil {
		return ID3Tags{}, err
	}
	defer tag.Close()
	return ID3Tags{
		Title:  tag.Title(),
		Artist: tag.Artist(),
		Album:  tag.Album(),
	}, nil
}

func (i *ID3Service) Write(file string, t ID3Tags) error {
	tag, err := id3v2.Open(file, id3v2.Options{
		Parse: true,
	})
	if err != nil {
		return err
	}
	defer tag.Close()
	tag.SetTitle(t.Title)
	tag.SetArtist(t.Artist)
	tag.SetAlbum(t.Album)
	return tag.Save()
}
